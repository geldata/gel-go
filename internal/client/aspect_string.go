// Code generated by "stringer -type Aspect"; DO NOT EDIT.

package gel

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[DataDescription-84]
}

const _Aspect_name = "DataDescription"

var _Aspect_index = [...]uint8{0, 15}

func (i Aspect) String() string {
	i -= 84
	if i >= Aspect(len(_Aspect_index)-1) {
		return "Aspect(" + strconv.FormatInt(int64(i+84), 10) + ")"
	}
	return _Aspect_name[_Aspect_index[i]:_Aspect_index[i+1]]
}
