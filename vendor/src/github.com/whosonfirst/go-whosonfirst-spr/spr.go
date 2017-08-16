package spr

import (
	"github.com/whosonfirst/go-whosonfirst-flags"
)

type StandardPlacesResult interface {
	Id() int64
	ParentId() int64
	Name() string
	Placetype() string
	Country() string
	Repo() string
	Path() string
	URI() string
	IsCurrent() (flags.ExistentialFlag, error)
	IsCeased() (flags.ExistentialFlag, error)
	IsDeprecated() (flags.ExistentialFlag, error)
	IsSuperseded() (flags.ExistentialFlag, error)
	IsSuperseding() (flags.ExistentialFlag, error)
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
