// Code generated by "stringer -type=step"; DO NOT EDIT.

package annotation

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[initialStep-0]
	_ = x[annotationNameStep-1]
	_ = x[attributeNameStep-2]
	_ = x[attributeValueStep-3]
	_ = x[doneStep-4]
}

const _step_name = "initialStepannotationNameStepattributeNameStepattributeValueStepdoneStep"

var _step_index = [...]uint8{0, 11, 29, 46, 64, 72}

func (i step) String() string {
	if i < 0 || i >= step(len(_step_index)-1) {
		return "step(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _step_name[_step_index[i]:_step_index[i+1]]
}
