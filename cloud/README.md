# cloud

对象存储

- 支持 `s3`、`qcloud`、`gcs` 对象存储操作


# 例子

```go
sb := MustNew(StorageTypeS3, key, secret, "xxxx", WithRegion("us-east-2"))
err := sb.PutObject(context.Background(), src, strings.NewReader(str), len(str))
xpanic.WhenError(err)
fmt.Println("OK")
```

Output:
```text
OK
```