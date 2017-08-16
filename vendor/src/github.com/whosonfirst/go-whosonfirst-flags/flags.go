package flags

type ExistentialFlag interface {
	Flag() int64
	IsTrue() bool
	IsFalse() bool
	IsKnown() bool
}
