package cloud

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestCRUD(t *testing.T) {
	key := os.Getenv("RELEASE_CLOUD_KEY")
	secret := os.Getenv("RELEASE_CLOUD_SECRET")
	if len(key) == 0 || len(secret) == 0 {
		return
	}
	sb := MustNew(StorageTypeS3, key, secret, "zhongtai", WithRegion("us-east-2"))
	str := "test"
	src := "testtesttest"
	dest := "testtesttest_dest"
	err := sb.PutObject(context.Background(), src, strings.NewReader(str), len(str))
	if err != nil {
		t.Error(err)
	}
	info, err := sb.StatObject(context.Background(), src)
	if err != nil {
		fmt.Println(info)
		t.Error(err)
	}
	myChan := sb.ListObjects(context.Background(), "test")
	for v := range myChan {
		fmt.Println(v)
	}
	err = sb.CopyObject(context.Background(), dest, src)
	myChan = sb.ListObjects(context.Background(), "test")
	for v := range myChan {
		fmt.Println(v)
	}
	err = sb.DelObject(context.Background(), src)
	if err != nil {
		t.Error(err)
	}
	err = sb.DelObject(context.Background(), dest)
	if err != nil {
		t.Error(err)
	}
}
