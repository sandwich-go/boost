package sset

// go get github.com/sandwich-go/gotemplate/...
//go:generate gotemplate -outfmt gen_%v "../templates/sset" "Int(int)"
//go:generate gotemplate -outfmt gen_%v "../templates/sset" "Int8(int8)"
//go:generate gotemplate -outfmt gen_%v "../templates/sset" "Int16(int16)"
//go:generate gotemplate -outfmt gen_%v "../templates/sset" "Int32(int32)"
//go:generate gotemplate -outfmt gen_%v "../templates/sset" "Int64(int64)"
//go:generate gotemplate -outfmt gen_%v "../templates/sset" "Uint(uint)"
//go:generate gotemplate -outfmt gen_%v "../templates/sset" "Uint8(uint8)"
//go:generate gotemplate -outfmt gen_%v "../templates/sset" "Uint16(uint16)"
//go:generate gotemplate -outfmt gen_%v "../templates/sset" "Uint32(uint32)"
//go:generate gotemplate -outfmt gen_%v "../templates/sset" "Uint64(uint64)"
//go:generate gotemplate -outfmt gen_%v "../templates/sset" "String(string)"
//go:generate gotemplate -outfmt gen_%v "../templates/sset" "Any(interface{})"
