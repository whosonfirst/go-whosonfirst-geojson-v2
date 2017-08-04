package spr

type StandardPlacesResult interface {
	Id() int64
	ParentId() int64
	Name() string
	Placetype() string
	Country() string
	Repo() string
	Path() string
	URI() string
	IsCurrent() bool
	IsCeased() bool
	IsDeprecated() bool
	IsSuperseded() bool
	IsSuperseding() bool
	SupersededBy() []int64
	Supersedes() []int64
}

type StandardPlacesResults interface {
	Results() []StandardPlacesResult
	// Pagination ?
}
