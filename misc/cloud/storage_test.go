package cloud

import (
	"context"
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"strings"
	"testing"
)

func TestCloud(t *testing.T) {
	key := os.Getenv("RELEASE_CLOUD_KEY")
	secret := os.Getenv("RELEASE_CLOUD_SECRET")
	if len(key) == 0 || len(secret) == 0 {
		return
	}
	sb := MustNew(StorageTypeS3, key, secret, "zhongtai", WithRegion("us-east-2"))

	testUtil(sb, t)
}

func TestQCloud(t *testing.T) {
	// https://cloud.tencent.com/document/faq/436/102489
	// virtual-hosted-style
	key := os.Getenv("RELEASE_QCLOUD_KEY")
	secret := os.Getenv("RELEASE_QCLOUD_SECRET")
	bucket := os.Getenv("RELEASE_QCLOUD_BUCKET")
	if len(key) == 0 ||
		len(secret) == 0 ||
		len(bucket) == 0 {
		return
	}
	sb := MustNew(StorageTypeQCloud, key, secret, bucket, WithRegion("ap-beijing"))

	testUtil(sb, t)
}

func TestAliCloud(t *testing.T) {
	key := os.Getenv("RELEASE_ALICLOUD_KEY")
	secret := os.Getenv("RELEASE_ALICLOUD_SECRET")
	bucket := os.Getenv("RELEASE_ALICLOUD_BUCKET")
	if len(key) == 0 ||
		len(secret) == 0 ||
		len(bucket) == 0 {
		return
	}
	sb := MustNew(StorageTypeAliCS, key, secret, bucket, WithRegion("us-east-1"))
	testUtil(sb, t)
}

func testUtil(sb Storage, t *testing.T) {
	Convey("put/stat/list/copy object", t, func() {
		str := "test"
		src := "testtesttest"
		err := sb.PutObject(context.Background(), src, strings.NewReader(str), len(str))
		So(err, ShouldBeNil)

		info, err0 := sb.StatObject(context.Background(), src)
		So(err0, ShouldBeNil)
		t.Log(info)

		myChan := sb.ListObjects(context.Background(), "test")
		for v := range myChan {
			t.Log(v)
		}

		dest := "testtesttest_dest"
		err = sb.CopyObject(context.Background(), dest, src)
		So(err, ShouldBeNil)

		myChan = sb.ListObjects(context.Background(), "test")
		for v := range myChan {
			t.Log(v)
		}
		err = sb.DelObject(context.Background(), src)
		So(err, ShouldBeNil)
		err = sb.DelObject(context.Background(), dest)
		So(err, ShouldBeNil)
	})

	Convey("resolve", t, func() {
		s, err := sb.ResolveObjectName("https://fsadfdsa.com/zhongtai/ddd")
		So(err, ShouldBeNil)
		t.Log(s)
	})
}
