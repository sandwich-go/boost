# xchan

Unbounded chan with ring buffer.

Refer to the below articles and issues:
1. https://github.com/golang/go/issues/20352
2. https://stackoverflow.com/questions/41906146/why-go-channels-limit-the-buffer-size
3. https://medium.com/capital-one-tech/building-an-unbounded-channel-in-go-789e175cd2cd
4. https://erikwinter.nl/articles/2020/channel-with-infinite-buffer-in-golang/


## Usage

```go
ch := NewUnboundedChan(1000)
// or ch := NewUnboundedChanSize(10,200,1000)

go func() {
    for ...... {
        ...
        ch.In <- ... // send values
        ...
    }

    close(ch.In) // close In channel
}()


for v := range ch.Out { // read values
    fmt.Println(v)
}

```