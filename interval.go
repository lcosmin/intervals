package intervals

import (
	"fmt"
	"strings"
)

var Empty = Interval{Empty: true}

// Interval is a structure representing an interval of integer numbers, [Lower, Upper]
type Interval struct {
	// Lower is the lower bound of the interval; can be a concrete number or the special constant
	// Infinity (interpreted as -infinity)
	Lower BoundType
	// Upper is the upper bound of the interval; can be a concrete number or the special constant
	// Infinity
	Upper BoundType
	// Empty indicates that this interval is empty
	Empty bool
}

func (i *Interval) setLower(b BoundType) {
	i.Lower = boundCopy(b)
	i.Empty = false
}

func (i *Interval) setUpper(b BoundType) {
	i.Upper = boundCopy(b)
	i.Empty = false
}

func (i Interval) String() string {

	if i.Empty {
		return "(empty)"
	}

	sb := &strings.Builder{}
	sb.WriteRune('(')
	if i.Lower == Infinity {
		sb.WriteString("-inf, ")
	} else {
		sb.WriteString(fmt.Sprintf("%d", int64(*i.Lower)))
		sb.WriteString(", ")
	}

	if i.Upper == Infinity {
		sb.WriteString("inf)")
	} else {
		sb.WriteString(fmt.Sprintf("%d", int64(*i.Upper)))
		sb.WriteRune(')')
	}

	return sb.String()
}

func (i Interval) Overlaps(o Interval) bool {
	// Test for cases when intervals DON'T overlap
	if i == Empty || o == Empty {
		return false
	}

	// 1) interval o comes before interval i
	if o.Upper != nil && i.Lower != nil && *o.Upper < *i.Lower {
		return false
	}

	// 2) interval o comes after interval i
	if o.Lower != nil && i.Upper != nil && *o.Lower > *i.Upper {
		return false
	}

	return true
}

func (i Interval) Intersect(o Interval) Interval {

	if !i.Overlaps(o) {
		return Empty
	}

	var r Interval

	if i.Lower == Infinity && i.Upper == Infinity {
		// (-inf, +inf) & anything -> anything
		return NewInterval(o.Lower, o.Upper)
	}

	if o.Lower == Infinity && o.Upper == Infinity {
		// anything & (-inf, +inf) -> anything
		return NewInterval(i.Lower, i.Upper)
	}

	if i.Lower == Infinity {
		// i is (-i, i.Upper]

		// o can be:

		// (-i, o.Upper] ; in this case, result is (-i, min(i.Upper, o.Upper))
		if o.Lower == Infinity {
			r.setLower(Infinity)
			if *o.Upper < *i.Upper {
				r.setUpper(o.Upper)
			} else {
				r.setUpper(i.Upper)
			}
			return r
		}

		// [o.Lower, i)
		// [o.Lower, o.Upper]

		// Intervals overlap for sure, so the lower bound is o.Lower
		r.setLower(o.Lower)

		if o.Upper == Infinity || *i.Upper < *o.Upper {
			r.setUpper(i.Upper)
		} else {
			r.setUpper(o.Upper)
		}

		return r
	}

	// i is [i.Lower, i)
	if i.Upper == Infinity {
		// o can be:

		// (-i, o.Upper]
		if o.Lower == Infinity {
			return NewInterval(i.Lower, o.Upper)
		}

		// [o.Lower, i)
		// [o.Lower, o.Upper]

		r.setUpper(o.Upper)
		if *i.Lower > *o.Lower {
			r.setLower(i.Lower)
		} else {
			r.setLower(o.Lower)
		}
		return r
	}

	// i is [A, B]

	// o cant be:
	// (-i, o.Upper]
	if o.Lower == Infinity {
		r.setLower(i.Lower)

		if *o.Upper < *i.Upper {
			r.setUpper(o.Upper)
		} else {
			r.setUpper(i.Upper)
		}
		return r
	}

	// [o.Lower, i)
	if o.Upper == Infinity {
		r.setUpper(i.Upper)

		if *i.Lower < *o.Lower {
			r.setLower(o.Lower)
		} else {
			r.setLower(i.Lower)
		}
		return r
	}

	// [o.Lower, o.Upper]
	if *o.Lower < *i.Lower {
		r.setLower(i.Lower)
	} else {
		r.setLower(o.Lower)
	}
	if *o.Upper < *i.Upper {
		r.setUpper(o.Upper)
	} else {
		r.setUpper(i.Upper)
	}
	return r
}

func NewInterval(lower, upper BoundType) Interval {
	if lower != nil && upper != nil && *lower > *upper {
		return Empty
	}

	var r Interval

	r.setLower(lower)
	r.setUpper(upper)
	return r
}
