package cloud

import (
	"fmt"
	"net/url"
	"strings"
)

const StorageTypeGCS StorageType = "gcs" // google 云

func init() {
	register(StorageTypeGCS, newGCStorage)
}

type gcsStorage struct {
	*baseStorage
}

func newGCStorage(accessKeyID string, secretAccessKey string, bucket string, opts ...StorageOption) (Storage, error) {
	bb, err := newBaseBucket(accessKeyID, secretAccessKey, bucket, func(options *StorageOptions) (ep string, err error) {
		return "storage.googleapis.com", nil
	}, opts...)
	if err != nil {
		return nil, err
	}
	s := &gcsStorage{baseStorage: bb}
	s.setObjectNameResolver(func(u *url.URL) (string, error) {
		return strings.TrimPrefix(u.Path, "/"+s.bucket+"/"), nil
	})
	return s, nil
}

func (c gcsStorage) GetRootUrl() string {
	return fmt.Sprintf("https://storage.googleapis.com/%s", c.bucket)
}
