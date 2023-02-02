package dns

//go:generate stringer -type=Policy
type Policy int

const (
	PolicyFirst  Policy = 0 // 使用第一个ip
	PolicyRandom Policy = 1 // 随机一个可用ip
)
