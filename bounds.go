package intervals

type BoundType *int64

var (
	Infinity BoundType = nil
)

func Bound(v int64) BoundType {
	return BoundType(&v)
}

func Copy(b BoundType) BoundType {
	if b == nil {
		return Infinity
	}
	return Bound(*b)
}
