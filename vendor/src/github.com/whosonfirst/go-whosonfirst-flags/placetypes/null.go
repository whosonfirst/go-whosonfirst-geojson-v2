package placetypes

import (
	"github.com/whosonfirst/go-whosonfirst-flags"
)

type NullFlag struct {
	flags.PlacetypesFlag
}

func NewNullFlag() (*NullFlag, error) {

	f := NullFlag{}
	return &f, nil
}

func (f *NullFlag) Matches(other flags.PlacetypesFlag) bool {
     	return true
}

func (f *NullFlag) Placetypes() []string {
	return []string{}
}

func (f *NullFlag) String() string {
	return "NULL"
}
