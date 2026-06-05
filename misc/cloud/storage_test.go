package cloud

import (
	"context"
	"os"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// 这组 cloud 端到端测试需要真实云服务凭证，仅在本地集成场景跑。
// CI 不设 RELEASE_*_KEY/SECRET env，会显式 t.Skipf 跳过（不是 silent
// PASS），日志里看得清是 SKIP 不是 RUN。
//
// 历史：原代码用 'if len(key)==0 { return }' silent 早返，CI 上看着像
// 真跑了。改 t.Skipf 让 SKIP 在测试日志里有显式输出便于排查。

func TestCloud(t *testing.T) {
	key := os.Getenv("RELEASE_CLOUD_KEY")
	secret := os.Getenv("RELEASE_CLOUD_SECRET")
	if len(key) == 0 || len(secret) == 0 {
		t.Skipf("need RELEASE_CLOUD_KEY / RELEASE_CLOUD_SECRET env to run AWS S3 integration test")
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
		t.Skipf("need RELEASE_QCLOUD_{KEY,SECRET,BUCKET} env to run Tencent COS integration test")
	}
	sb := MustNew(StorageTypeQCloud, key, secret, bucket, WithRegion("ap-beijing"))

	testUtil(sb, t)
}
func TestMinio(t *testing.T) {
	key := os.Getenv("RELEASE_MINIO_KEY")
	secret := os.Getenv("RELEASE_MINIO_SECRET")
	bucket := os.Getenv("RELEASE_MINIO_BUCKET")
	region := os.Getenv("RELEASE_MINIO_REGION")
	if len(key) == 0 ||
		len(secret) == 0 ||
		len(bucket) == 0 ||
		len(region) == 0 {
		t.Skipf("need RELEASE_MINIO_{KEY,SECRET,BUCKET,REGION} env to run Minio integration test")
	}
	sb := MustNew(StorageTypeMinio, key, secret, bucket, WithRegion(region))

	testUtil(sb, t)
}

func TestAliCloud(t *testing.T) {
	key := os.Getenv("RELEASE_ALICLOUD_KEY")
	secret := os.Getenv("RELEASE_ALICLOUD_SECRET")
	bucket := os.Getenv("RELEASE_ALICLOUD_BUCKET")
	if len(key) == 0 ||
		len(secret) == 0 ||
		len(bucket) == 0 {
		t.Skipf("need RELEASE_ALICLOUD_{KEY,SECRET,BUCKET} env to run AliCloud OSS integration test")
	}
	sb := MustNew(StorageTypeAliCS, key, secret, bucket, WithRegion("us-east-1"))
	testUtil(sb, t)
}

func TestHuaweiRUCloud(t *testing.T) {
	key := os.Getenv("RELEASE_HUAWEIRU_KEY")
	secret := os.Getenv("RELEASE_HUAWEIRU_SECRET")
	bucket := os.Getenv("RELEASE_HUAWEIRU_BUCKET")
	if len(key) == 0 ||
		len(secret) == 0 ||
		len(bucket) == 0 {
		t.Skipf("need RELEASE_HUAWEIRU_{KEY,SECRET,BUCKET} env to run Huawei OBS RU integration test")
	}
	sb := MustNew(StorageTypeHuaweiRU, key, secret, bucket, WithRegion("ru-moscow-1"))
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
