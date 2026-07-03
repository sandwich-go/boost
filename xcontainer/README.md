# xcontainer

# 容器

泛型容器集合。自 1.4 起，各容器统一改用 Go 泛型实现（`any.go`），旧的
`gotemplate` 多类型生成代码（`gen_*.go` / `templates/`）已删除。用泛型类型参数
指定元素类型即可，无需再选 `NewInt8` / `NewInt32String` 这类具名构造。

包含以下容器：

- `ringbuf` —— 非协程安全的字节环形缓冲区。
- `sarray` —— 数组 `Array[T]`，提供协程安全（`NewSync`）与非协程安全（`New`）两种。
- `slist` —— 双向链表 `List[T]`，提供协程安全与非协程安全两种。
- `sset` —— 集合 `Set[T]`，提供协程安全与非协程安全两种。
- `smap` —— 分片的协程安全映射 `Concurrent[K, V]`。
- `syncmap` —— 同步映射：`SyncMap[K, V]`（封装 `sync.Map`）与 `BucketMap[K, V]`（分桶）。
- `redblacktree` —— 红黑树 `Tree[K, V]`，有序键，支持 Floor/Ceiling。
- `sortedmap` —— 插入序有序映射 `Map[K, V]`，基于二分插入。

> 迁移提示（旧具名 API → 泛型）：
> `sarray.NewInt8()` → `sarray.New[int8]()`，
> `slist.NewSyncAny()` → `slist.NewSync[any]()`，
> `sset.NewAny()` → `sset.New[any]()`，
> `smap.NewInt32String()` → `smap.New[int32, string]()`，
> `syncmap.NewInt8Int()` → `syncmap.New[int8, int]()`。

## ringbuf

`ringbuf` 是一个非协程安全的字节环形缓冲区。

```go
import "github.com/sandwich-go/boost/xcontainer/ringbuf"

buf := ringbuf.New(10)
fmt.Println(buf.Capacity()) // 10
_ = buf.Write([]byte("helloworld"))

tmp := make([]byte, 5)
buf.Read(tmp, 5)
fmt.Println(string(tmp)) // hello
```

## sarray

`sarray.New[T]()` 创建非协程安全数组；`sarray.NewSync[T]()` 创建带读写锁、可在多
协程中并发使用的数组。类型参数 `T` 需满足 `comparable`。

```go
import "github.com/sandwich-go/boost/xcontainer/sarray"

tr := sarray.New[int8]()
tr.PushLeft(3)
v, _ := tr.Get(0)
fmt.Println(v) // 3

_ = tr.InsertBefore(0, 11)
v, _ = tr.Get(0)
fmt.Println(v) // 11

fmt.Println(tr.Contains(11)) // true
tr.DeleteValue(11)
fmt.Println(tr.Contains(11)) // false

fmt.Println(tr.Search(3)) // 0
```

## slist

`slist.New[T]()` 创建非协程安全链表；`slist.NewSync[T]()` 创建协程安全链表。
类型参数 `T` 为 `any`。

```go
import "github.com/sandwich-go/boost/xcontainer/slist"

tr := slist.NewSync[int8]()
tr.PushBack(8)
tr.PushBack(9)
fmt.Println(tr.Len()) // 2
tr.PushFront(7)
fmt.Println(tr.PopFrontAll()) // [7 8 9]
```

## sset

`sset.New[T]()` 创建非协程安全集合；`sset.NewSync[T]()` 创建协程安全集合。
类型参数 `T` 需满足 `comparable`。

```go
import "github.com/sandwich-go/boost/xcontainer/sset"

s := sset.New[int]()
s.Add(1, 2, 3)
fmt.Println(s.Contains(2)) // true
fmt.Println(s.AddIfNotExist(2)) // false（已存在）
fmt.Println(s.AddIfNotExist(4)) // true
```

## smap

`smap` 提供分片的协程安全映射 `Concurrent[K, V]`（`K` 需 `comparable`）。
`New` 使用默认分片数；`NewWithSharedCount` 指定分片数量。

```go
import "github.com/sandwich-go/boost/xcontainer/smap"

tr := smap.New[int32, string]()
tr.Set(1, "1")
fmt.Println(tr.Len()) // 1
v, ret := tr.Get(1)
fmt.Println(v, ret) // 1 true

tr2 := smap.NewWithSharedCount[int32, string](64) // 指定分片数量为 64
_ = tr2
```

## syncmap

`syncmap` 提供两种同步映射：

- `SyncMap[K, V]`：封装标准库 `sync.Map`，附带 `LoadOrStoreFuncLock` 等一次填充语义。用 `New[K, V]()` 构造。
- `BucketMap[K, V]`：按用户提供的 `hashFunc(K) int64` 分桶的并发映射，`K` 需 `comparable`。用 `NewBucketMap` 构造；`NewBucketMapWithCacheKey` 额外维护键快照。

```go
import "github.com/sandwich-go/boost/xcontainer/syncmap"

// SyncMap（封装 sync.Map）
tr := syncmap.New[int8, int]()
tr.Store(1, 2)
v, ok := tr.Load(1)
fmt.Println(ok, v) // true 2

// BucketMap（分桶，需提供 hashFunc）
bm := syncmap.NewBucketMap[string, int](16, func(k string) int64 {
    return int64(len(k))
})
bm.Store("a", 1)
bv, bok := bm.Load("a")
fmt.Println(bok, bv) // true 1
```

## redblacktree

`redblacktree.Tree[K, V]` 是红黑树，键有序，支持 `Floor`/`Ceiling`/`Walk` 及节点池化回收。
`New` 用于 `cmp.Ordered` 键；`NewWith` 接受自定义 `Comparator[K]`。

```go
import "github.com/sandwich-go/boost/xcontainer/redblacktree"

tr := redblacktree.New[int, string]()
tr.Put(2, "two")
tr.Put(1, "one")
tr.Put(3, "three")
v, found := tr.Get(2)
fmt.Println(found, v) // true two
```

## sortedmap

`sortedmap.Map[K, V]` 是插入序有序映射（基于二分插入），非协程安全。
构造时需提供键比较函数 `sort func(K, K) int`。

```go
import (
    "cmp"
    "github.com/sandwich-go/boost/xcontainer/sortedmap"
)

m := sortedmap.New[int, string](cmp.Compare[int])
m.Set(2, "two")
m.Set(1, "one")
v, ok := m.Get(1)
fmt.Println(ok, v) // true one
```
