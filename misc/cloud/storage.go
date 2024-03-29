package cloud

import (
	"bytes"
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sandwich-go/boost/xpanic"
	"io"
	"net/url"
	"strings"
)

var (
	creators        = map[StorageType]storageCreator{}
	emptyObjectInfo = ObjectInfo{}
)

func register(st StorageType, creator storageCreator) { creators[st] = creator }

type (
	// StorageType 存储器类型
	StorageType string
	// storageCreator 存储器创建
	storageCreator = func(accessKeyID string, secretAccessKey string, bucket string, opts ...StorageOption) (Storage, error)
	// endpointGetter endpoint 获取器
	endpointGetter = func(*StorageOptions) (ep string, err error)
	// objectNameResolver 解析 object name
	objectNameResolver = func(*url.URL) (string, error)
)

// New 新建 Storage
func New(st StorageType, accessKeyID string, secretAccessKey string, bucket string, opts ...StorageOption) (Storage, error) {
	creator, ok := creators[st]
	if !ok {
		return nil, ErrUnknownStorageType
	}
	return creator(accessKeyID, secretAccessKey, bucket, opts...)
}

// MustNew 新建 Storage，失败会 panic
func MustNew(st StorageType, accessKeyID string, secretAccessKey string, bucket string, opts ...StorageOption) Storage {
	s, err := New(st, accessKeyID, secretAccessKey, bucket, opts...)
	xpanic.WhenError(err)
	return s
}

type baseStorage struct {
	cli      *minio.Client
	bucket   string
	resolver objectNameResolver
	spec     *StorageOptions
}

func newBaseBucket(accessKeyID string, secretAccessKey string, bucket string, getter endpointGetter, opts ...StorageOption) (*baseStorage, error) {
	spec := NewStorageOptions(opts...)
	ep, err := getter(spec)
	if err != nil {
		return nil, err
	}
	var cli *minio.Client
	cli, err = minio.New(ep, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}
	return &baseStorage{cli: cli, bucket: bucket, spec: spec}, nil
}

func (c *baseStorage) setObjectNameResolver(resolver objectNameResolver) {
	c.resolver = resolver
}

func (c baseStorage) ResolveObjectName(rawUrl string) (string, error) {
	u, err := url.Parse(rawUrl)
	if err != nil {
		return "", err
	}
	if c.resolver != nil {
		return c.resolver(u)
	}
	return strings.Trim(u.Path, "/"), nil
}

func (c baseStorage) String() string { return c.bucket }
func (c baseStorage) DelObject(ctx context.Context, objName string) error {
	return c.cli.RemoveObject(ctx, c.bucket, objName, minio.RemoveObjectOptions{})
}

func (c baseStorage) PutObject(ctx context.Context, objName string, reader io.Reader, objSize int, opts ...PutOption) (err error) {
	spec := NewPutOptions(opts...)
	op := minio.PutObjectOptions{
		ContentType:        spec.ContentType,
		ContentDisposition: spec.ContentDisposition,
		CacheControl:       spec.CacheControl,
	}
	if objSize == 0 {
		op.DisableMultipart = true
	}
	_, err = c.cli.PutObject(ctx, c.bucket, objName, reader, int64(objSize), op)
	return
}

func (c baseStorage) StatObject(ctx context.Context, objName string) (ObjectInfo, error) {
	info, err := c.cli.StatObject(ctx, c.bucket, objName, minio.StatObjectOptions{})
	if err != nil {
		return emptyObjectInfo, err
	}
	return ObjectInfo{
		ETag:         info.ETag,
		Key:          info.Key,
		LastModified: info.LastModified,
		Size:         info.Size,
		ContentType:  info.ContentType,
	}, nil
}

func (c baseStorage) ListObjects(ctx context.Context, prefix string) <-chan ObjectInfo {
	ch0 := c.cli.ListObjects(ctx, c.bucket, minio.ListObjectsOptions{
		Prefix: prefix,
	})
	ch1 := make(chan ObjectInfo)
	go func() {
		defer close(ch1)
		for info := range ch0 {
			ch1 <- ObjectInfo{
				ETag:         info.ETag,
				Key:          info.Key,
				LastModified: info.LastModified,
				Size:         info.Size,
				ContentType:  info.ContentType,
			}
		}
	}()
	return ch1
}

func (c baseStorage) GetObject(ctx context.Context, objName string) (io.Reader, error) {
	obj, err := c.cli.GetObject(ctx, c.bucket, objName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer func(o io.ReadCloser) {
		_ = o.Close()
	}(obj)
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, obj)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func (c baseStorage) CopyObject(ctx context.Context, destObjName, srcObjName string) error {
	_, err := c.cli.CopyObject(ctx, minio.CopyDestOptions{
		Bucket: c.bucket,
		Object: destObjName,
	}, minio.CopySrcOptions{
		Bucket: c.bucket,
		Object: srcObjName,
	})
	return err
}
