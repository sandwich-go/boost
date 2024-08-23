package cloud

import (
	"fmt"
	"strings"
)

const StorageTypeAliCS StorageType = "alics" // 腾讯云

const (
	XAliCSMetaMD5 = "x-oss-meta-file-md5"
)

func init() {
	register(StorageTypeAliCS, newAliCloudStorage)
}

type aliCloudStorage struct {
	*baseStorage
}

func newAliCloudStorage(accessKeyID string, secretAccessKey string, bucket string, opts ...StorageOption) (Storage, error) {
	bb, err := newBaseBucket(accessKeyID, secretAccessKey, bucket, func(options *StorageOptions) (ep string, err error) {
		if options.Region == "" {
			return "oss.aliyuncs.com", nil
		}
		region := options.Region
		region = strings.TrimPrefix(region, "oss-")
		return fmt.Sprintf("oss-%s.aliyuncs.com", region), nil
	}, opts...)
	if err != nil {
		return nil, err
	}
	return &aliCloudStorage{baseStorage: bb}, nil
}

func (c aliCloudStorage) GetRootUrl() string {
	return fmt.Sprintf("https://%s.oss-%s.aliyuncs.com", c.bucket, c.spec.GetRegion())
}
