package flags

type ExistentialFlag interface {
	Flag() int64
	IsTrue() bool
	IsFalse() bool
	IsKnown() bool
	Matches(ExistentialFlag) bool
	String() string
}

type PlacetypesFlag interface {
	Matches(PlacetypesFlag) bool
	Placetypes() []string
	String() string
}
