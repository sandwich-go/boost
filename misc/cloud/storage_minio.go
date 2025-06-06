package cloud

const StorageTypeMinio StorageType = "minio" // aws s3

func init() {
	register(StorageTypeMinio, newMinioStorage)
}

type minioStorage struct {
	*baseStorage
}

func newMinioStorage(accessKeyID string, secretAccessKey string, bucket string, opts ...StorageOption) (Storage, error) {
	bb, err := newBaseBucket(accessKeyID, secretAccessKey, bucket, func(options *StorageOptions) (ep string, err error) {
		if len(options.GetRegion()) == 0 {
			err = ErrRegionShouldNotEmptyForS3
		} else {
			ep = options.GetRegion()
		}
		return
	}, opts...)
	if err != nil {
		return nil, err
	}
	return &minioStorage{baseStorage: bb}, nil
}

func (c minioStorage) GetRootUrl() string {
	return c.baseStorage.spec.GetRegion()
}
