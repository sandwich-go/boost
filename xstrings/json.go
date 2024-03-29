package xstrings

import "encoding/json"

// UglyJSON removes insignificant space characters from the input json byte slice
// and returns the compacted result.
func UglyJSON(json []byte) []byte {
	buf := make([]byte, 0, len(json))
	return ugly(buf, json)
}

func ugly(dst, src []byte) []byte {
	dst = dst[:0]
	for i := 0; i < len(src); i++ {
		if src[i] > ' ' {
			dst = append(dst, src[i])
			if src[i] == '"' {
				for i = i + 1; i < len(src); i++ {
					dst = append(dst, src[i])
					if src[i] == '"' {
						j := i - 1
						for ; ; j-- {
							if src[j] != '\\' {
								break
							}
						}
						if (j-i)%2 != 0 {
							break
						}
					}
				}
			}
		}
	}
	return dst
}

// ValidJSON 是否为合法json字符串
func ValidJSON(bb []byte) bool {
	if len(bb) == 0 {
		return true
	}
	var jsonStr interface{}
	err := json.Unmarshal(bb, &jsonStr)
	return err == nil
}
