package vigenere

// Sanitize 剔除无效的字节
func Sanitize(key []byte) []byte {
	if len(key) == 0 {
		return key
	}
	var out = make([]byte, 0, len(key))
	for _, v := range key {
		if 65 <= v && v <= 90 {
			out = append(out, v)
		} else if 97 <= v && v <= 122 {
			out = append(out, v-32)
		}
	}
	return out
}

func encode(a, b byte) byte {
	return (((a - 'A') + (b - 'A')) % 26) + 'A'
}

func decode(a, b byte) byte {
	return ((((a - 'A') - (b - 'A')) + 26) % 26) + 'A'
}

// Encrypt 使用 key 进行加密
func Encrypt(src []byte, key []byte) []byte {
	if len(key) == 0 {
		return src
	}
	out := make([]byte, 0, len(src))
	for i, v := range src {
		out = append(out, encode(v, key[i%len(key)]))
	}
	return out
}

// Decrypt 使用 key 进行解密
func Decrypt(src, key []byte) []byte {
	if len(key) == 0 {
		return src
	}
	out := make([]byte, 0, len(src))
	for i, v := range src {
		out = append(out, decode(v, key[i%len(key)]))
	}
	return out
}

// EncryptAndInplace  使用 key 进行加密，并且替换 src
func EncryptAndInplace(src, key []byte) {
	if len(key) == 0 {
		return
	}
	for i, v := range src {
		src[i] = encode(v, key[i%len(key)])
	}
}

// DecryptAndInplace  使用 key 进行解密，并且替换 src
func DecryptAndInplace(src, key []byte) {
	if len(key) == 0 {
		return
	}
	for i, v := range src {
		src[i] = decode(v, key[i%len(key)])
	}
}
