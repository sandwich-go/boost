package cloud

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/hex"
	"github.com/sandwich-go/boost/xpanic"
	"github.com/sandwich-go/minio-go"
	"github.com/sandwich-go/minio-go/pkg/credentials"
	"io"
	"net/http"
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
	opts = append(opts, WithStorageType(st))
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

func (c baseStorage) Type() StorageType {
	return c.spec.StorageType
}

func (c baseStorage) String() string {
	return c.bucket
}

func (c baseStorage) DelObject(ctx context.Context, objName string) error {
	return c.cli.RemoveObject(ctx, c.bucket, objName, minio.RemoveObjectOptions{})
}

func (c baseStorage) PutObject(ctx context.Context, objName string, reader io.Reader, objSize int, opts ...PutOption) (err error) {
	spec, err := c.toSpec(opts...)
	if err != nil {
		return err
	}
	op := minio.PutObjectOptions{
		ContentType:          spec.ContentType,
		ContentDisposition:   spec.ContentDisposition,
		CacheControl:         spec.CacheControl,
		CustomHeaders:        spec.CustomHeader,
		SendContentMd5:       spec.SendContentMd5,
		UserMetadata:         spec.CustomMeta,
		DisableContentSha256: spec.DisableContentSha256,
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
		Meta:         info.UserMetadata,
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
				Meta:         info.UserMetadata,
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

func (c baseStorage) toSpec(opts ...PutOption) (*PutOptions, error) {
	storageType := c.spec.StorageType
	spec := NewPutOptions(opts...)
	var defaultOpts []PutOption
	if storageType == StorageTypeGCS || storageType == StorageTypeQCloud || storageType == StorageTypeAliCS {
		defaultOpts = append(defaultOpts, WithDisableContentSha256(true))
	}
	if (storageType == StorageTypeQCloud) && len(spec.FileMD5) != 0 {
		const CosHeaderMD5Key = "x-cos-meta-md5"
		defaultOpts = append(defaultOpts, WithCustomHeader(http.Header{CosHeaderMD5Key: []string{spec.FileMD5}}))
	}
	if (storageType == StorageTypeAliCS || storageType == StorageTypeGCS) && len(spec.FileMD5) != 0 {
		hexStyleMD5 := spec.FileMD5
		binaryData := make([]byte, len(hexStyleMD5)/2)
		_, err := hex.Decode(binaryData, []byte(hexStyleMD5))
		if err != nil {
			return nil, err
		}
		md5Base64ed := base64.StdEncoding.EncodeToString(binaryData)
		defaultOpts = append(defaultOpts, WithCustomHeader(http.Header{"Content-MD5": []string{md5Base64ed}}))
		defaultOpts = append(defaultOpts, WithCustomMeta(map[string]string{XAliCSMetaMD5: spec.FileMD5}))
	}

	if storageType == StorageTypeGCS {
		defaultOpts = append(defaultOpts, WithSendContentMd5(true))
	}
	fullOpts := append(defaultOpts, opts...)
	return NewPutOptions(fullOpts...), nil
}
