# encryption

加密

- 支持 `aes` 加解密
- 支持生成 `SecretKey`、`PublicKey`、`SharedKey`


# 例子

```go
// 生成随机私钥
secretKey0, err0 := curve25519.GenerateSecretKey()
if err0 != nil {
    panic(err0)
}
secretKey1, err1 := curve25519.GenerateSecretKey()
if err1 != nil {
    panic(err1)
}

// 通过私钥生成公钥
publicKey0, err2 := curve25519.GeneratePublicKey(secretKey0)
if err2 != nil {
    panic(err2)
}
publicKey1, err3 := curve25519.GeneratePublicKey(secretKey1)
if err3 != nil {
    panic(err3)
}

// 交换共享密钥
sharedKey0, err4 := curve25519.GenerateSharedKey(secretKey0, publicKey1)
if err4 != nil {
    panic(err4)
}
sharedKey1, err5 := curve25519.GenerateSharedKey(secretKey1, publicKey0)
if err5 != nil {
    panic(err5)
}

// aes 加解密
frame := []byte("time.Duration,[]time.Duration,map[string]*Redis此类的非基础类型的slice或者map都需要辅助指明类型")
encryptFrame, err6 := aes.Encrypt(frame, secretKey0)
if err6 != nil {
    panic(err6)
}
decryptFrame, err7 := aes.Decrypt(encryptFrame, secretKey0)
if err7 != nil {
    panic(err7)
}
if string(frame) != string(decryptFrame) {
    panic("encrypt/decrypt fail")
}
fmt.Println("OK")
```

Output:
```text
OK
```