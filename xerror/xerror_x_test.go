package xerror_test

import (
	"errors"
	"fmt"
	"github.com/sandwich-go/boost/xerror"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestErrors(t *testing.T) {
	Convey("xerror errors", t, func() {
		{
			err := xerror.Wrap(errors.New("1"), "wrap with xerror")
			So(errors.Unwrap(err), ShouldNotBeNil)
		}
		{
			err := xerror.NewText("1")
			So(errors.Unwrap(err), ShouldBeNil)
		}
		{
			errChild := xerror.NewText("1")
			errParent := xerror.Wrap(errChild, "wrap with xerror")
			So(errors.Is(errParent, errChild), ShouldBeTrue)
		}
	})
}

func TestCause(t *testing.T) {
	Convey("xerror case from system errors", t, func() {
		{
			err := errors.New("1")
			So(xerror.Cause(err), ShouldEqual, err)
		}
		{
			err1 := errors.New("1")
			err := xerror.Wrap(err1, "2")
			err = xerror.Wrap(err, "3")
			So(xerror.Cause(err), ShouldEqual, err1)
		}
	})

	Convey("xerror case from xerror new", t, func() {
		{
			err := xerror.NewText("1")
			So(xerror.Cause(err).Error(), ShouldEqual, err.Error())
		}
		{
			err1 := xerror.NewText("1")
			err := xerror.Wrap(err1, "2")
			err = xerror.Wrap(err, "3")
			So(xerror.Cause(err).Error(), ShouldEqual, err1.Error())
		}
	})
}

func TestStack(t *testing.T) {
	Convey("xerror stack", t, func() {
		{
			err := errors.New("1")
			So(xerror.Stack(err), ShouldEqual, err.Error())
		}
		{
			old := xerror.IsErrorWithStack
			xerror.IsErrorWithStack = true
			defer func() {
				xerror.IsErrorWithStack = old
			}()
			err := xerror.New(xerror.WithText("io error"), xerror.WithStack())
			errW := xerror.Wrap(err, "link error")
			errW = xerror.Wrap(errW, "session error")
			fmt.Printf("no stack===> %v\n", errW)
			fmt.Printf("with stack===> %+s\n", errW)

			fmt.Println(xerror.Caller(err, 0))         // xerror_x_test.go github.com/sandwich-go/boost/xerror/xerror_test.TestStack.func1 51
			fmt.Println(xerror.Caller(err.Cause(), 0)) // xerror_x_test.go github.com/sandwich-go/boost/xerror/xerror_test.TestStack.func1 49
		}
	})
}

func TestCode(t *testing.T) {
	Convey("xerror code", t, func() {
		{
			err := errors.New("1")
			So(xerror.Code(err), ShouldEqual, xerror.ErrorCodeUnsetAsDefault)
			err = xerror.Wrap(err, "wrap 1")
			So(xerror.Code(err), ShouldEqual, xerror.ErrorCodeUnsetAsDefault)
			err = xerror.WrapCode(2, err, "wrap 2")
			So(xerror.Code(err), ShouldEqual, 2)
			So(xerror.Wrap(nil, "wrap nil"), ShouldEqual, nil)

			So(xerror.Code(xerror.Wrap(nil, "wrap nil")), ShouldEqual, xerror.ErrorCodeOk)
		}
	})
}

func TestUserCase(t *testing.T) {
	{
		err := xerror.NewText("error info is %s,error code will be %d", "crash", xerror.ErrorCodeUnsetAsDefault)
		fmt.Println(err.Error())
		fmt.Println(xerror.Code(err))
	}
	{
		err := xerror.NewCode(10000, "error info")
		fmt.Println(err.Error())
		fmt.Println(xerror.Code(err))
	}
	{
		err1 := errors.New("from some lib")
		err2 := xerror.WrapCode(10000, err1, "wrap with error code")
		fmt.Println(err2.Error())
		fmt.Println(xerror.Code(err1)) // 使用默认的 error code：xerror.ErrorCodeUnsetAsDefault
		fmt.Println(xerror.Code(err2))
	}
}

func TestDyncOpenStackInfo(t *testing.T) {
	old := xerror.IsErrorWithStack
	xerror.IsErrorWithStack = false
	defer func() {
		xerror.IsErrorWithStack = old
	}()
	fmt.Printf("with stack===> %+s\n", xerror.NewText("open stack by function").WithStack())
	fmt.Printf("with stack===> %+s\n", xerror.New(xerror.WithText("open stack by option"), xerror.WithStack()))
}

func TestWrapNilError(t *testing.T) {
	Convey("wrap nil error", t, func() {
		{
			nilErr := func() error {
				return nil
			}
			err := xerror.Wrap(nilErr(), "should nil")
			So(err, ShouldBeNil)
			{
				var errInterface error
				errInterface = err
				if errInterface == nil {
					fmt.Printf("errInterface is nil, and error should be nil:%v\n", err)
				}
			}
			{
				errInterface := err
				if errInterface == nil {
					fmt.Printf("errInterface is nil, and error should be nil:%v\n", err)
				}
			}
			testFiler := func(msg interface{}) {
				switch msg.(type) {
				case error:
					fmt.Println("msg is a error, but is nil", msg)
				default:
					fmt.Println("msg is not error")
				}
			}
			testFiler(err)
		}
	})
}
