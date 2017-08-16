package spr

type ExistentialFlag interface {
     StatusInt() int64
     Status() bool     		     
     Confidence() bool
}

type StandardPlacesResult interface {
	Id() int64
	ParentId() int64
	Name() string
	Placetype() string
	Country() string
	Repo() string
	Path() string
	URI() string
	IsCurrent() ExistentialFlag
	IsCeased() ExistentialFlag
	IsDeprecated() ExistentialFlag
	IsSuperseded() ExistentialFlag
	IsSuperseding() ExsitentialFlag
	SupersededBy() []int64
	Supersedes() []int64
}

type Pagination interface {
	Pages() int
	Page() int
	PerPage() int
	Total() int
	Cursor() string
	NextQuery() string
}

type StandardPlacesResults interface {
	Results() []StandardPlacesResult
}
