package segmenttree

import "strconv"

type SegErrno uintptr

// Errors
const (
	EBADINDEX = SegErrno(0x1)
)

// Error table
var errors = [...]string{
	1: "bad index",
}

func (e SegErrno) Error() string {
	if 0 <= int(e) && int(e) < len(errors) {
		s := errors[e]
		if s != "" {
			return s
		}
	}
	return "errno " + strconv.FormatInt(int64(e), 10)
}
