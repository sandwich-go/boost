package xcompress

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// 本文件目标：补 xcompress 包内未覆盖的导出 API + 错误路径：
//   - 包级 Flat / Inflate（之前 0%，调 Default 全局压缩器）
//   - MustNew panic 路径
//   - gzip / snappy Inflate 在非法数据上返错误
//   - Type.String（之前 0%，stringer 生成）
//
// 已有 *_test.go 覆盖 gzip / snappy / dummy 的 happy path。

func TestPackageLevel_FlatInflate(t *testing.T) {
	Convey("包级 Flat / Inflate 走 Default 全局 GZIP 压缩器", t, func() {
		input := []byte("hello world, this is a test message for compression.")
		compressed, err := Flat(input)
		So(err, ShouldBeNil)
		So(len(compressed), ShouldBeGreaterThan, 0)

		decompressed, err := Inflate(compressed)
		So(err, ShouldBeNil)
		So(decompressed, ShouldResemble, input)
	})

	Convey("包级 Flat 处理空 bytes", t, func() {
		compressed, err := Flat([]byte{})
		So(err, ShouldBeNil)
		// gzip 即使空数据也有 header，所以 compressed 非空
		decompressed, err := Inflate(compressed)
		So(err, ShouldBeNil)
		So(len(decompressed), ShouldEqual, 0)
	})
}

func TestMustNew_PanicOnInvalidLevel(t *testing.T) {
	Convey("MustNew GZIP + 非法 level panic", t, func() {
		// gzip 合法 level: -2 (HuffmanOnly), -1 (Default), 0 (NoCompression),
		// 1 (BestSpeed), ..., 9 (BestCompression)
		// 100 是非法的
		So(func() { MustNew(WithType(GZIP), WithLevel(100)) }, ShouldPanic)
		So(func() { MustNew(WithType(GZIP), WithLevel(-99)) }, ShouldPanic)
	})

	Convey("MustNew 合法参数返回 Compressor", t, func() {
		c := MustNew(WithType(GZIP), WithLevel(BestSpeed))
		So(c, ShouldNotBeNil)
	})
}

func TestGzip_Inflate_InvalidData(t *testing.T) {
	Convey("Gzip Inflate 在非法数据上返错误", t, func() {
		c, err := New(WithType(GZIP), WithLevel(DefaultCompression))
		So(err, ShouldBeNil)

		// 完全不是 gzip 格式
		_, err = c.Inflate([]byte("this is not gzip"))
		So(err, ShouldNotBeNil)

		// 截断的 gzip header
		_, err = c.Inflate([]byte{0x1f, 0x8b}) // 只有 magic bytes
		So(err, ShouldNotBeNil)
	})
}

func TestSnappy_Inflate_InvalidData(t *testing.T) {
	Convey("Snappy Inflate 在非法数据上返错误", t, func() {
		c, err := New(WithType(Snappy))
		So(err, ShouldBeNil)

		// 不是 snappy 格式
		_, err = c.Inflate([]byte("not snappy"))
		So(err, ShouldNotBeNil)
	})
}

func TestType_String(t *testing.T) {
	Convey("Type.String stringer 输出", t, func() {
		// Dummy / GZIP / Snappy 三个 Type 值都应该有 String() 输出
		So(Dummy.String(), ShouldNotBeEmpty)
		So(GZIP.String(), ShouldNotBeEmpty)
		So(Snappy.String(), ShouldNotBeEmpty)

		// 不同 Type 的 String 输出应该不同
		So(Dummy.String(), ShouldNotEqual, GZIP.String())
		So(GZIP.String(), ShouldNotEqual, Snappy.String())
	})
}

// BenchmarkGzip / BenchmarkSnappy hot path bench：压缩是 IO 路径里的常见
// 步骤（缓存 / RPC payload），值得量化两种实现的成本对比。
func BenchmarkFlat_Gzip(b *testing.B) {
	c, _ := New(WithType(GZIP), WithLevel(DefaultCompression))
	payload := make([]byte, 4096)
	for i := range payload {
		payload[i] = byte(i % 256)
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = c.Flat(payload)
	}
}

func BenchmarkFlat_Snappy(b *testing.B) {
	c, _ := New(WithType(Snappy))
	payload := make([]byte, 4096)
	for i := range payload {
		payload[i] = byte(i % 256)
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = c.Flat(payload)
	}
}

func BenchmarkInflate_Gzip(b *testing.B) {
	c, _ := New(WithType(GZIP), WithLevel(DefaultCompression))
	payload := make([]byte, 4096)
	for i := range payload {
		payload[i] = byte(i % 256)
	}
	compressed, _ := c.Flat(payload)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = c.Inflate(compressed)
	}
}

func BenchmarkInflate_Snappy(b *testing.B) {
	c, _ := New(WithType(Snappy))
	payload := make([]byte, 4096)
	for i := range payload {
		payload[i] = byte(i % 256)
	}
	compressed, _ := c.Flat(payload)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = c.Inflate(compressed)
	}
}
