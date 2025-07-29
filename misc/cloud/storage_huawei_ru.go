package cloud

import "fmt"

const StorageTypeHuaweiRU StorageType = "huaweiru" // 华为云俄罗斯区域

func init() {
	register(StorageTypeHuaweiRU, newHuaweiRUStorage)
}

type huaweiRUStorage struct {
	*baseStorage
}

func newHuaweiRUStorage(accessKeyID string, secretAccessKey string, bucket string, opts ...StorageOption) (Storage, error) {
	bb, err := newBaseBucket(accessKeyID, secretAccessKey, bucket, func(options *StorageOptions) (ep string, err error) {
		if len(options.GetRegion()) == 0 {
			ep = "obs.ru-moscow-1.hc.sbercloud.ru"
		} else {
			ep = fmt.Sprintf("obs.%s.hc.sbercloud.ru", options.GetRegion())
		}
		return
	}, opts...)
	if err != nil {
		return nil, err
	}
	return &huaweiRUStorage{baseStorage: bb}, nil
}

func (c huaweiRUStorage) GetRootUrl() string {
	return fmt.Sprintf("https://%s.obs.%s.hc.sbercloud.ru", c.bucket, c.spec.GetRegion())
}
