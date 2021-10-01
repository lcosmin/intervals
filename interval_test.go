package intervals

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var re = regexp.MustCompile(`(?P<first>empty|\[(?P<l1>-i|-?\d+),(?P<r1>i|-?\d+)\])&(?P<second>empty|\[(?P<l2>-i|-?\d+),(?P<r2>i|-?\d+)\])=(?P<result>empty|\[(?P<l3>-i|-?\d+),(?P<r3>i|-?\d+)\])`)

func TestOverlapping(t *testing.T) {
	i := NewInterval(Infinity, Infinity) // (-inf, +inf)

	assert.Equal(t,
		i.Overlaps(NewInterval(Infinity, Infinity)),
		true,
		"(-inf, +inf) overlaps empty")

	assert.Equal(t,
		i.Overlaps(Empty),
		false,
		"(-inf, +inf) overlaps empty",
	)

	assert.Equal(t,
		i.Overlaps(Empty),
		false,
		"empty overlaps (-inf, +inf)",
	)

	assert.Equal(t,
		i.Overlaps(NewInterval(Bound(1), Bound(2))),
		true,
		"(-inf, +inf) overlaps (1, 2)")

	// (0, +inf)
	i = NewInterval(Bound(0), Infinity)

	assert.Equal(t, i.Overlaps(NewInterval(Infinity, Infinity)), true, "(0, +inf) overlaps (-inf, +inf)")
	assert.Equal(t, i.Overlaps(NewInterval(Bound(0), Infinity)), true, "(0, +inf) overlaps (0, inf)")
	assert.Equal(t, i.Overlaps(NewInterval(Bound(1), Infinity)), true, "(0, +inf) overlaps (1, inf)")
	assert.Equal(t, i.Overlaps(NewInterval(Infinity, Bound(0))), true, "(0, +inf) overlaps (-inf, 0)")
	assert.Equal(t, i.Overlaps(NewInterval(Bound(0), Bound(1))), true, "(0, +inf) overlaps (0, 1)")
	assert.Equal(t, i.Overlaps(NewInterval(Infinity, Bound(-1))), false, "(0, +inf) doesn't overlap (-inf, -1)")

	// (0, 0)
	i = NewInterval(Bound(0), Bound(0))

	assert.Equal(t, i.Overlaps(NewInterval(Infinity, Infinity)), true, "(0, 0) overlaps (-inf, +inf)")
	assert.Equal(t, i.Overlaps(NewInterval(Bound(0), Infinity)), true, "(0, 0) overlaps (0, inf)")
	assert.Equal(t, i.Overlaps(NewInterval(Infinity, Bound(0))), true, "(0, 0) overlaps (-inf, 0)")
	assert.Equal(t, i.Overlaps(NewInterval(Bound(0), Bound(0))), true, "(0, 0) overlaps (0, 0)")
	assert.Equal(t, i.Overlaps(NewInterval(Bound(0), Bound(1))), true, "(0, 0) overlaps (0, 1)")

	assert.Equal(t, i.Overlaps(NewInterval(Infinity, Bound(-1))), false, "(0, 0) doesn't overlap (-inf, -1)")
	assert.Equal(t, i.Overlaps(NewInterval(Bound(1), Infinity)), false, "(0, 0) doesn't overlap (1, inf)")

	// (0, 10)
	i = NewInterval(Bound(0), Bound(10))

	assert.Equal(t, i.Overlaps(NewInterval(Infinity, Infinity)), true, "(0, 10) overlaps (-inf, +inf)")

	assert.Equal(t, i.Overlaps(NewInterval(Infinity, Bound(-1))), false, "(0, 10) doesn't overlap (-inf, -1)")
	assert.Equal(t, i.Overlaps(NewInterval(Bound(11), Infinity)), false, "(0, 10) doesn't overlap (11, inf)")

	assert.Equal(t, i.Overlaps(NewInterval(Bound(0), Bound(10))), true, "(0, 10) overlaps (0, 10)")
	assert.Equal(t, i.Overlaps(NewInterval(Bound(0), Bound(0))), true, "(0, 10) overlaps (0, 0)")
	assert.Equal(t, i.Overlaps(NewInterval(Bound(10), Bound(10))), true, "(0, 10) overlaps (10, 10)")
	assert.Equal(t, i.Overlaps(NewInterval(Bound(10), Bound(10))), true, "(0, 10) overlaps (10, 10)")

	assert.Equal(t, i.Overlaps(NewInterval(Bound(-1), Bound(11))), true, "(0, 10) overlaps (-1, 11)")
	assert.Equal(t, i.Overlaps(NewInterval(Bound(1), Bound(5))), true, "(0, 10) overlaps (1, 5)")
	assert.Equal(t, i.Overlaps(NewInterval(Bound(-1), Bound(1))), true, "(0, 10) overlaps (-1, 1)")
	assert.Equal(t, i.Overlaps(NewInterval(Bound(5), Bound(11))), true, "(0, 10) overlaps (5, 11)")

}

func TestIntersections(t *testing.T) {

	mktest := func(s string) (Interval, Interval, Interval) {

		s1 := strings.ReplaceAll(s, " ", "")

		// []string{"", "first", "l1", "r1", "second", "l2", "r2", "result", "l3", "r3"}
		res := re.FindStringSubmatch(s1)
		if res == nil {
			panic("failed")
		}

		var (
			r1, r2, r3 Interval
			l, r       BoundType
		)

		if res[1] == "empty" {
			r1 = Empty
		} else {
			if res[2] == "-i" {
				l = Infinity
			} else {
				i, _ := strconv.ParseInt(res[2], 10, 64)
				l = Bound(i)
			}

			if res[3] == "i" {
				r = Infinity
			} else {
				i, _ := strconv.ParseInt(res[3], 10, 64)
				r = Bound(i)
			}

			r1 = NewInterval(l, r)
		}

		if res[4] == "empty" {
			r2 = Empty
		} else {
			if res[5] == "-i" {
				l = Infinity
			} else {
				i, _ := strconv.ParseInt(res[5], 10, 64)
				l = Bound(i)
			}

			if res[6] == "i" {
				r = Infinity
			} else {
				i, _ := strconv.ParseInt(res[6], 10, 64)
				r = Bound(i)
			}

			r2 = NewInterval(l, r)
		}

		if res[7] == "empty" {
			r3 = Empty
		} else {
			if res[8] == "-i" {
				l = Infinity
			} else {
				i, _ := strconv.ParseInt(res[8], 10, 64)
				l = Bound(i)
			}

			if res[9] == "i" {
				r = Infinity
			} else {
				i, _ := strconv.ParseInt(res[9], 10, 64)
				r = Bound(i)
			}

			r3 = NewInterval(l, r)
		}

		return r1, r2, r3
	}

	tests := []string{
		"empty & [-i, i] = empty",
		"[-i, i] & empty = empty",
		"[-i, i] & [-i, i] = [-i, i]",
		"[-i, i] & [-i, 0] = [-i, 0]",
		"[-i, i] & [0,  0] = [0,  0]",
		"[-i, i] & [0, 10] = [0, 10]",
		"[-i, 0] & [-i, i] = [-i, 0]",
		"[-i, 0] & [-i, -1] = [-i, -1]",
		"[-i, 0] & [-i, 0] = [-i, 0]",
		"[-i, 0] & [-i, 1] = [-i, 0]",
		"[-i, 0] & [-1, i] = [-1, 0]",
		"[-i, 0] & [0, i] = [0, 0]",
		"[-i, 0] & [1, i] = empty",
		"[-i, 0] & [1, 1] = empty",
		"[-i, 0] & [1, 2] = empty",
		"[-i, 0] & [0, 0] = [0, 0]",
		"[-i, 0] & [-1, -1] = [-1, -1]",
		"[-i, 0] & [-1, 0] = [-1, 0]",
		"[-i, 0] & [-1, 1] = [-1, 0]",
		"[0, i] & [-i, i] = [0, i]",
		"[0, i] & [-i, -1] = empty",
		"[0, i] & [-i, 0] = [0, 0]",
		"[0, i] & [-i, 1] = [0, 1]",
		"[0, i] & [-1, i] = [0, i]",
		"[0, i] & [0, i] = [0, i]",
		"[0, i] & [1, i] = [1, i]",
		"[0, i] & [1, 1] = [1, 1]",
		"[0, i] & [1, 2] = [1, 2]",
		"[0, i] & [0, 0] = [0, 0]",
		"[0, i] & [-1, -1] = empty",
		"[0, i] & [-1, 0] = [0, 0]",
		"[0, i] & [-1, 1] = [0, 1]",

		"[0, 0] & [-i, i] = [0, 0]",
		"[0, 0] & [-i, -1] = empty",
		"[0, 0] & [-i, 0] = [0, 0]",
		"[0, 0] & [-i, 1] = [0, 0]",
		"[0, 0] & [-1, i] = [0, 0]",
		"[0, 0] & [0, i] = [0, 0]",
		"[0, 0] & [1, i] = empty",

		"[0, 0] & [0, 0] = [0, 0]",
		"[0, 0] & [-1, 1] = [0, 0]",
		"[0, 0] & [-2, -1] = empty",
		"[0, 0] & [0, 1] = [0, 0]",
		"[0, 0] & [1, 2] = empty",

		"[0, 2] & [-i, i] = [0, 2]",
		"[0, 2] & [-i, -1] = empty",
		"[0, 2] & [-i, 0] = [0, 0]",
		"[0, 2] & [-i, 1] = [0, 1]",
		"[0, 2] & [-i, 3] = [0, 2]",
		"[0, 2] & [-1, i] = [0, 2]",
		"[0, 2] & [0, i] = [0, 2]",
		"[0, 2] & [1, i] = [1, 2]",
		"[0, 2] & [3, i] = empty",

		"[0, 2] & [0, 0] = [0, 0]",
		"[0, 2] & [1, 1] = [1, 1]",
		"[0, 2] & [0, 1] = [0, 1]",
		"[0, 2] & [1, 3] = [1, 2]",
		"[0, 2] & [-2, -1] = empty",
		"[0, 2] & [-1, 2] = [0, 2]",
		"[0, 2] & [-1, 1] = [0, 1]",

		"[0, 2] & [3, 4] = empty",
	}

	for idx, testCase := range tests {

		i, o, expected := mktest(testCase)

		res := i.Intersect(o)

		assert.Equal(t, res, expected, fmt.Sprintf("failed test %d, %s (%s, %s): got %s", idx, testCase, i, o, res))

	}

}

func TestInvalidIntervals(t *testing.T) {
	i := NewInterval(Bound(1), Bound(0))

	assert.Equal(t, i, Empty, "return empty interval if invalid")
}

func TestParticular(t *testing.T) {
	i := NewInterval(Bound(0), Bound(1))
	o := NewInterval(Infinity, Infinity)

	res := i.Intersect(o)

	assert.Equal(t, res, i, "fail")
}
