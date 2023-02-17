package sarray

// go get github.com/sandwich-go/gotemplate/...
//go:generate gotemplate -outfmt gen_%v "../templates/sarray" "Int(int)"
//go:generate gotemplate -outfmt gen_%v "../templates/sarray" "Int8(int8)"
//go:generate gotemplate -outfmt gen_%v "../templates/sarray" "Int16(int16)"
//go:generate gotemplate -outfmt gen_%v "../templates/sarray" "Int32(int32)"
//go:generate gotemplate -outfmt gen_%v "../templates/sarray" "Int64(int64)"
//go:generate gotemplate -outfmt gen_%v "../templates/sarray" "Uint(uint)"
//go:generate gotemplate -outfmt gen_%v "../templates/sarray" "Uint8(uint8)"
//go:generate gotemplate -outfmt gen_%v "../templates/sarray" "Uint16(uint16)"
//go:generate gotemplate -outfmt gen_%v "../templates/sarray" "Uint32(uint32)"
//go:generate gotemplate -outfmt gen_%v "../templates/sarray" "Uint64(uint64)"
//go:generate gotemplate -outfmt gen_%v "../templates/sarray" "String(string)"
//go:generate gotemplate -outfmt gen_%v "../templates/sarray" "Any(interface{})"
