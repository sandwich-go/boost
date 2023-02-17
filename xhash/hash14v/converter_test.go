package hash14v

import (
	"fmt"
	"github.com/sandwich-go/boost/xcrypto/algorithm/vigenere"
	"math"
	"testing"
)

func Test_HashID(t *testing.T) {
	test(t, nil)
}

func Test_HashID_WithHashKey(t *testing.T) {
	test(t, []byte("test"))
}

func test(t *testing.T, hk []byte) {
	oo := New(WithHashKey(hk), WithHashOffset([]byte("FAAAAAA")), WithUsingReservedBuff(true))
	conv := oo.(*converter)
	hashOffsetUsing := conv.hashOffset
	hashOffsetUsing = vigenere.Encrypt(hashOffsetUsing, conv.hashKey)
	uid := oo.ToId(hashOffsetUsing)
	if uid != 0 {
		t.Fatalf("uid with offset got:%d , except:%d", uid, 0)
	}
	hashID := oo.ToV(uid)
	if string(hashID) != string(hashOffsetUsing) {
		t.Fatalf("hashID got:%s , except:%s", string(hashID), string(hashOffsetUsing))
	}

	{
		uid = uint64(0)
		hashID = oo.ToV(uid)
		fmt.Println(fmt.Sprintf("with hash offset:%s uid:%d got:%s ", hashOffsetUsing, uid, hashID))
		if oo.ToId(hashID) != uid {
			t.Fatalf("got:%d , except:%d", oo.ToId(hashID), uid)
		}
	}
	{
		uid = math.MaxUint64 - oo.Offset()
		hashID = oo.ToV(uid)
		fmt.Println(fmt.Sprintf("with hash offset:%s uid:%d got:%s ", hashOffsetUsing, uid, hashID))
		if oo.ToId(hashID) != uid {
			t.Fatalf("got:%d , except:%d", oo.ToId(hashID), uid)
		}

		if oo.ToV(uint64(math.MaxUint64)) != nil {
			t.Fatal("should overflow")
		}
	}
	{
		var hashEqualOffsetLengthMax []byte
		var hashEqualOffsetLengthMaxSamePrefix []byte
		for range hashOffsetUsing {
			hashEqualOffsetLengthMax = append(hashEqualOffsetLengthMax, 'Z')
			hashEqualOffsetLengthMaxSamePrefix = append(hashEqualOffsetLengthMaxSamePrefix, 'Z')
		}
		hashEqualOffsetLengthMaxSamePrefix[0] = conv.hashOffset[0]
		vigenere.EncryptAndInplace(hashEqualOffsetLengthMax, conv.hashKey)
		vigenere.EncryptAndInplace(hashEqualOffsetLengthMaxSamePrefix, conv.hashKey)
		uid1 := oo.ToId(hashEqualOffsetLengthMax)
		uid2 := oo.ToId(hashEqualOffsetLengthMaxSamePrefix)
		fmt.Println(fmt.Sprintf("with hash offset:%s, %s got :%d  %s got :%d", hashOffsetUsing, string(hashEqualOffsetLengthMax), uid1, hashEqualOffsetLengthMaxSamePrefix, uid2))
	}
}

func BenchmarkFromUInt64(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ToV(123456789)
	}
}

func BenchmarkToUInt64(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ToId([]byte("YESTTES"))
	}
}
