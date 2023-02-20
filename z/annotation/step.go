package annotation

//go:generate stringer -type=step
type step int

const (
	initialStep        step = iota // 初始化
	annotationNameStep             // 注释名
	attributeNameStep              // 属性名
	attributeValueStep             // 属性值
	doneStep                       // 结束
)
