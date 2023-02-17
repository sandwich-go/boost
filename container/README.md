# container

##容器
包含以下容器：
- `Ringbuf`环形缓冲区。
- `SyncArray`提供了多种数据类型的数组实现，包含协程安全和非协程安全的版本。
- `SyncList`提供了多种数据类型的链表实现，包含协程安全和非协程安全的版本。
- `SyncSet`提供了多种数据类型的集合的实现，可提供协程安全或者非协程安全的版本。
- `SyncMap`提供了一个同步的映射实现。

### Ringbuf
`Ringbuf`是一个非协程安全的环形缓冲区

####例子
```go
import "github.com/sandwich-go/boost/container/ringbuf"

var buf = ringbuf.New(10)
fmt.Println(buf.Capacity())
err := buf.Write([]byte("helloworld"))

tmp := make([]byte, 5)
buf.Read(tmp, 5)
fmt.Println(string(tmp))

buf.Read(tmp, 4)
fmt.Println(string(tmp))
```

Output:
```text
10
hello
worlo
```

### SyncArray
`SyncArray`可以创建出非线程安全的数组，
```go
import "github.com/sandwich-go/boost/container/sarray"

sarray.NewInt8()
```
也可以创建出带读写锁的`数组`，从而在多个协程中安全地并发使用。
```go
import "github.com/sandwich-go/boost/container/sarray"

sarray.NewSyncInt8()
```
####例子

```go
import "github.com/sandwich-go/boost/container/sarray"

tr := sarray.NewInt8()
tr.PushLeft(3)
fmt.Println(tr.Get(0)) // 3

tr.InsertBefore(0, 11)
fmt.Println(tr.Get(0)) // 11

fmt.Println(tr.Contains(11)) // true

tr.DeleteValue(11)
fmt.Println(tr.Contains(11)) // false

fmt.Println(tr.Search(3)) // 0
```

Output:
```txt
3
11
true
false
0
```

### SyncList
`SyncList`提供了存储多种数据类型的双向链表的实现，包含提供协程安全的版本或非协程安全的版本。

```go
import "github.com/sandwich-go/boost/container/slist"

slist.NewInt8() // 非协程安全的版本

slist.NewSyncInt8() // 协程安全的版本

```

#### 例子
```go
import "github.com/sandwich-go/boost/container/slist"

tr := slist.NewSyncInt8()
tr.PushBack(8)
tr.PushBack(9)
fmt.Println(tr.Len()) // 2
tr.PushFront(7)
fmt.Println(tr.PopFrontAll()) // 7 8 9
```

Output:
```txt
2
7 8 9
```

### SyncMap
`SyncMap`提供了一个同步的映射实现。

#### 例子
````go
import "github.com/sanndwich-go/boost/container/syncmap"

tr := syncmap.NewInt8Int()
tr.Store(1, 2)
v, h := tr.Load(1)
fmt.Println(h)
fmt.Println(v)
````

Output:
```txt
true
2
```