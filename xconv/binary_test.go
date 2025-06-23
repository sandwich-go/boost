package xconv

import (
	"fmt"
	"math"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestBinary(t *testing.T) {
	Convey(`test binary`, t, func() {
		var buf = []byte{0x00, 0x00, 0x48, 0x42}
		So(LittleEndianDecodeToFloat32(buf), ShouldEqual, float32(50))
		buf = []byte{24, 45, 68, 84, 251, 33, 9, 64}
		So(LittleEndianDecodeToFloat64(buf), ShouldEqual, 3.141592653589793)
		buf = []byte{0x48}
		So(LittleEndianDecodeToInt64(buf), ShouldEqual, 72)
		So(LittleEndianDecodeToUint64(buf), ShouldEqual, 72)
	})
}

func TestBinaryLittleEndian(t *testing.T) {
	Convey("Test Little-Endian Encoding/Decoding", t, func() {
		Convey("Int64", func() {
			testCases := []int64{
				0,
				72,
				123456789,
				-987654321,
				math.MaxInt64,
				math.MinInt64,
			}

			for idx, tc := range testCases {
				expected := tc
				encoded := LittleEndianEncodeFromInt64(expected)

				Convey(fmt.Sprintf("Decodes correctly %d", idx), func() {
					decoded := LittleEndianDecodeToInt64(encoded)
					So(decoded, ShouldEqual, expected)
				})

				Convey(fmt.Sprintf("Encodes/decodes symmetrically%d", idx), func() {
					decoded := LittleEndianDecodeToInt64(encoded)
					reEncoded := LittleEndianEncodeFromInt64(decoded)
					So(reEncoded, ShouldResemble, encoded)
				})
			}
		})

		Convey("Uint64", func() {
			testCases := []uint64{
				0,
				72,
				123456789,
				math.MaxUint64,
			}

			for idx, tc := range testCases {
				expected := tc
				encoded := LittleEndianEncodeFromUint64(expected)

				Convey(fmt.Sprintf("Decodes correctly %d", idx), func() {
					decoded := LittleEndianDecodeToUint64(encoded)
					So(decoded, ShouldEqual, expected)
				})

				Convey(fmt.Sprintf("Encodes/decodes symmetrically %d", idx), func() {
					decoded := LittleEndianDecodeToUint64(encoded)
					reEncoded := LittleEndianEncodeFromUint64(decoded)
					So(reEncoded, ShouldResemble, encoded)
				})
			}
		})

		Convey("Float32", func() {
			testCases := []float32{
				0.0,
				50.0,
				3.14159,
				-1.5,
				math.MaxFloat32,
				math.SmallestNonzeroFloat32,
			}

			for idx, tc := range testCases {
				expected := tc
				encoded := LittleEndianEncodeFromFloat32(expected)

				Convey(fmt.Sprintf("Decodes correctly %d", idx), func() {
					decoded := LittleEndianDecodeToFloat32(encoded)
					So(decoded, ShouldEqual, expected)
				})

				Convey(fmt.Sprintf("Encodes/decodes symmetrically %d", idx), func() {
					decoded := LittleEndianDecodeToFloat32(encoded)
					reEncoded := LittleEndianEncodeFromFloat32(decoded)
					So(reEncoded, ShouldResemble, encoded)
				})
			}

			// 保持与原有测试的兼容
			Convey("Existing test case passes", func() {
				buf := []byte{0x00, 0x00, 0x48, 0x42}
				So(LittleEndianDecodeToFloat32(buf), ShouldEqual, float32(50))
			})
		})

		Convey("Float64", func() {
			testCases := []float64{
				0.0,
				3.141592653589793,
				-1.5,
				math.MaxFloat64,
				math.SmallestNonzeroFloat64,
			}

			for idx, tc := range testCases {
				expected := tc
				encoded := LittleEndianEncodeFromFloat64(expected)

				Convey(fmt.Sprintf("Decodes correctly %d", idx), func() {
					decoded := LittleEndianDecodeToFloat64(encoded)
					So(decoded, ShouldEqual, expected)
				})

				Convey(fmt.Sprintf("Encodes/decodes symmetrically %d", idx), func() {
					decoded := LittleEndianDecodeToFloat64(encoded)
					reEncoded := LittleEndianEncodeFromFloat64(decoded)
					So(reEncoded, ShouldResemble, encoded)
				})
			}

			// 保持与原有测试的兼容
			Convey("Existing test case passes", func() {
				buf := []byte{24, 45, 68, 84, 251, 33, 9, 64}
				So(LittleEndianDecodeToFloat64(buf), ShouldEqual, 3.141592653589793)
			})
		})

		Convey("LeFillUpSize", func() {
			testCases := []struct {
				input  []byte
				length int
				expect []byte
			}{
				{[]byte{0x01}, 4, []byte{0x01, 0x00, 0x00, 0x00}},
				{[]byte{0x01, 0x02}, 2, []byte{0x01, 0x02}},
				{[]byte{0xFF, 0xEE, 0xDD}, 8, []byte{0xFF, 0xEE, 0xDD, 0x00, 0x00, 0x00, 0x00, 0x00}},
			}

			for idx, tc := range testCases {
				Convey(fmt.Sprintf("Pads bytes correctly %d", idx), func() {
					result := LeFillUpSize(tc.input, tc.length)
					So(result, ShouldResemble, tc.expect)
				})
			}
		})
	})
}
