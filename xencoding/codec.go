package xencoding

import (
	"context"
	"sort"
	"sync"
)

type codecKeyType struct{}

// for ark context
func (*codecKeyType) String() string { return "encoding2-codec—key" }

var keyForContext = codecKeyType(struct{}{})

// WithContext 将 Codec 存放在 context.Context 中
func WithContext(ctx context.Context, c Codec) context.Context {
	return context.WithValue(ctx, keyForContext, c)
}

// FromContext 从 context.Context 中 获取 Codec
func FromContext(ctx context.Context) Codec {
	c := ctx.Value(keyForContext)
	if c == nil {
		c = ctx.Value(keyForContext.String())
	}
	if c == nil {
		return nil
	}
	return c.(Codec)
}

// Codec defines the interface link uses to encode and decode messages.
type Codec interface {
	// Marshal returns the wire format of v.
	Marshal(v interface{}) ([]byte, error)
	// Unmarshal parses the wire format into v.
	Unmarshal(data []byte, v interface{}) error
	// Name String returns the name of the Codec implementation.
	Name() string
}

var (
	mu     sync.RWMutex
	codecs = make(map[string]Codec)
)

// RegisterCodec 注册 Codec
// 可以注册自定义的 Codec
func RegisterCodec(c Codec) {
	mu.Lock()
	defer mu.Unlock()

	if c == nil {
		panic("cannot register a nil Codec")
	}
	name := c.Name()
	if len(name) == 0 {
		panic("cannot register Codec with empty string result for Name()")
	}
	if _, dup := codecs[name]; dup {
		panic("register called twice for codec " + name)
	}
	codecs[name] = c
}

// Codecs 获取所有 Codec 的名称
func Codecs() []string {
	mu.RLock()
	defer mu.RUnlock()

	list := make([]string, 0, len(codecs))
	for name := range codecs {
		list = append(list, name)
	}
	sort.Strings(list)
	return list
}

// GetCodec 通过名称来获取注册过的 Codec
func GetCodec(name string) Codec {
	mu.RLock()
	defer mu.RUnlock()

	return codecs[name]
}
