package hide

import (
	"fmt"
	"testing"
)

func TestHide(t *testing.T) {
	fmt.Println(Do("123456789_should_hide_part", WithSuffix("gmt")))
	fmt.Println(Do("1234567", WithSuffix("gmt")))
}
