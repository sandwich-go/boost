package xbit

import (
	"fmt"
	"math"
)

var (
	ErrInvalidField = fmt.Errorf("invalid field, bit should be <= %d", FieldMax)
)

type (
	// Field 位
	Field byte //取值范围是[1,64]
	// FieldSet 位集
	FieldSet uint64
)

var (
	// FieldMax 最大位数
	FieldMax Field = 64
	// FieldSetFull 所有位数均为1
	FieldSetFull FieldSet = math.MaxUint64
)

// mustValidField 合法的位数
func mustValidField(i Field) {
	if i == 0 || i > FieldMax {
		panic(ErrInvalidField)
	}
}

// Set 位设置
// 返回值为 true 表示 FieldSet 还未设置该 Field，有效操作
// 返回值为 false, 表示 FieldSet 已经设置过该 Field，无效操作
func (x *FieldSet) Set(i Field) bool {
	mustValidField(i)
	orig := *x
	curr := orig | 1<<(i-1)
	if orig == curr {
		return false
	}

	*x = curr
	return true
}

// Clear 清除位
// 返回值为 true 表示 FieldSet 已经设置该 Field，有效操作
// 返回值为 false, 表示 FieldSet 未设置过该 Field，无效操作
func (x *FieldSet) Clear(i Field) bool {
	mustValidField(i)
	orig := *x
	curr := orig & ^(1 << (i - 1))
	if orig == curr {
		return false
	}

	*x = curr
	return true
}

// ClearAll 清理所有的位设置
func (x *FieldSet) ClearAll() { *x = 0 }

// Union 并集，不会修改原有值
func (x FieldSet) Union(fs FieldSet) FieldSet { return x | fs }

// Intersect 交集，不会修改原有值
func (x FieldSet) Intersect(fs FieldSet) FieldSet { return x & fs }

// IsSet 某位是否已经被设置
func (x FieldSet) IsSet(i Field) bool { return (x>>(i-1))&1 == 1 }
