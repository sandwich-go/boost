package cloud

import "fmt"

const StorageTypeQCloud StorageType = "qcloud" // 腾讯云

func init() {
	register(StorageTypeQCloud, newQCloudStorage)
}

type qCloudStorage struct {
	*baseStorage
}

func newQCloudStorage(accessKeyID string, secretAccessKey string, bucket string, opts ...StorageOption) (Storage, error) {
	bb, err := newBaseBucket(accessKeyID, secretAccessKey, bucket, func(options *StorageOptions) (ep string, err error) {
		if len(options.GetRegion()) == 0 {
			ep = "cos.ap-beijing.myqcloud.com"
		} else {
			ep = fmt.Sprintf("cos.%s.myqcloud.com", options.GetRegion())
		}
		return
	}, opts...)
	if err != nil {
		return nil, err
	}
	return &qCloudStorage{baseStorage: bb}, nil
}

func (c qCloudStorage) GetRootUrl() string {
	return fmt.Sprintf("https://%s.cos.%s.myqcloud.com", c.bucket, c.spec.GetRegion())
}
