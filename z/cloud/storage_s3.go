package cloud

import "fmt"

const StorageTypeS3 StorageType = "s3" // aws s3

func init() {
	register(StorageTypeS3, newS3Storage)
}

type s3Storage struct {
	*baseStorage
}

func newS3Storage(accessKeyID string, secretAccessKey string, bucket string, opts ...StorageOption) (Storage, error) {
	bb, err := newBaseBucket(accessKeyID, secretAccessKey, bucket, func(options *StorageOptions) (ep string, err error) {
		if len(options.GetRegion()) == 0 {
			err = ErrRegionShouldNotEmptyForS3
		} else {
			ep = fmt.Sprintf("s3.dualstack.%s.amazonaws.com", options.GetRegion())
		}
		return
	}, opts...)
	if err != nil {
		return nil, err
	}
	return &s3Storage{baseStorage: bb}, nil
}

func (c s3Storage) GetRootUrl() string {
	return fmt.Sprintf("https://%s.s3.amazonaws.com", c.bucket)
}
