package cloud

import (
	"context"
	"errors"
	"io"
	"time"
)

var (
	// ErrRegionShouldNotEmptyForS3 s3 的 region 必须进行设置，否则返回该错误
	ErrRegionShouldNotEmptyForS3 = errors.New("region should not empty for s3")
	// ErrUnknownStorageType 未知的 StorageType
	ErrUnknownStorageType = errors.New("unknown storage type")
)

// ObjectInfo container for object metadata.
type ObjectInfo struct {
	// An ETag is optionally set to md5sum of an object.  In case of multipart objects,
	// ETag is of the form MD5SUM-N where MD5SUM is md5sum of all individual md5sums of
	// each parts concatenated into one string.
	ETag         string    `json:"etag"`
	Key          string    `json:"name"`         // Name of the object
	LastModified time.Time `json:"lastModified"` // Date and time the object was last modified.
	Size         int64     `json:"size"`         // Size in bytes of the object.
	ContentType  string    `json:"contentType"`  // A standard MIME type describing the format of the object data.
}

// Storage 存储器
type Storage interface {
	// GetRootUrl 获取 root url，不同的 StorageType，拥有不同的 root url
	GetRootUrl() string
	// ResolveObjectName 解析 objectName
	ResolveObjectName(rawUrl string) (string, error)

	// DelObject 通过 objectName 删除指定的 object
	DelObject(ctx context.Context, objectName string) error
	// PutObject 将 object 存放至 Storage 中
	// PutObject Uploads objects that are less than 128MiB in a single PUT operation.Wr481iDXfbN1VntumTyHPM7f4IKDvTr4
	// For objects that are greater than 128MiB in size, PutObject seamlessly
	// uploads the object as parts of 128MiB or more depending on the actual file size.
	// The max upload size for an object is 5TB.
	PutObject(ctx context.Context, objectName string, reader io.Reader, objectSize int, options ...PutOption) (err error)
	// StatObject 通过 objectName 获取 object 元信息 ObjectInfo
	StatObject(ctx context.Context, objectName string) (ObjectInfo, error)
	// ListObjects 通过前缀批量获取 object 元信息 ObjectInfo
	ListObjects(ctx context.Context, prefix string) <-chan ObjectInfo
	// GetObject 通过 objectName 获取 object
	GetObject(ctx context.Context, objectName string) (reader io.Reader, err error)
	// CopyObject 拷贝名为 srcObjectName 的 object 至 destObjectName
	CopyObject(ctx context.Context, destObjectName, srcObjectName string) error
}
