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

	r := Interval{}

	if i.Lower == Infinity && i.Upper == Infinity {
		// (-inf, +inf) & anything -> anything
		return NewInterval(Copy(o.Lower), Copy(o.Upper))
	}

	if o.Lower == Infinity && o.Upper == Infinity {
		// (-inf, +inf) & anything -> anything
		return NewInterval(Copy(i.Lower), Copy(i.Upper))
	}

	if i.Lower == Infinity {
		// i is (-i, A]

		// o is
		// (-i, B]
		// [B, i)
		// [B, C]

		if o.Lower == Infinity {
			r.Lower = Infinity
			if *o.Upper < *i.Upper {
				r.Upper = Copy(o.Upper)
			} else {
				r.Upper = Copy(i.Upper)
			}
			return r
		}

		// we know for sure the intervals overlap
		r.Lower = Copy(o.Lower)

		if o.Upper == Infinity {
			r.Upper = Copy(i.Upper)
		} else {
			if *i.Upper < *o.Upper {
				r.Upper = Copy(i.Upper)
			} else {
				r.Upper = Copy(o.Upper)
			}
		}

		return r

	}

	if i.Upper == Infinity {

		// i is [A, i)

		// o is:
		// (-i, B]
		if o.Lower == Infinity {
			r.Lower = Copy(i.Lower)
			r.Upper = Copy(o.Upper)
			return r
		}

		// [B, i)
		// [B, C]

		r.Upper = Copy(o.Upper)
		if *i.Lower > *o.Lower {
			r.Lower = Copy(i.Lower)
		} else {
			r.Lower = Copy(o.Lower)
		}
		return r
	}

	// i is [A, B]

	// o is:
	// (-i, C]
	if o.Lower == Infinity {
		r.Lower = Copy(i.Lower)

		if *o.Upper < *i.Upper {
			r.Upper = Copy(o.Upper)
		} else {
			r.Upper = Copy(i.Upper)
		}
		return r
	}

	// [C, i)
	if o.Upper == Infinity {
		r.Upper = Copy(i.Upper)

		if *i.Lower < *o.Lower {
			r.Lower = Copy(o.Lower)
		} else {
			r.Lower = Copy(i.Lower)
		}
		return r
	}

	// [C, D]
	if *o.Lower < *i.Lower {
		r.Lower = Copy(i.Lower)
		if *o.Upper < *i.Upper {
			r.Upper = Copy(o.Upper)
		} else {
			r.Upper = Copy(i.Upper)
		}
		return r
	} else {
		r.Lower = Copy(o.Lower)
		if *o.Upper < *i.Upper {
			r.Upper = Copy(o.Upper)
		} else {
			r.Upper = Copy(i.Upper)
		}
		return r
	}
}

func NewInterval(lower, upper BoundType) Interval {
	if lower != nil && upper != nil && *lower > *upper {
		return Empty
	}

	return Interval{
		Lower: lower,
		Upper: upper,
		Empty: false,
	}
}
