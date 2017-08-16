package existential

import (
	"github.com/whosonfirst/go-whosonfirst-flags"
)

type KnownUnknownFlag struct {
	flags.ExistentialFlag
	flag       int64
	status     bool
	confidence bool
}

func NewKnownUnknownFlag(i int64) flags.ExistentialFlag {

	var status bool
	var confidence bool

	switch i {
	case 0:
		status = false
		confidence = true
	case 1:
		status = true
		confidence = true
	default:
		status = false
		confidence = false
	}

	f := KnownUnknownFlag{
		flag:       i,
		status:     status,
		confidence: confidence,
	}

	return &f
}

func (f *KnownUnknownFlag) Flag() int64 {
	return f.flag
}

func (f *KnownUnknownFlag) IsTrue() bool {
	return f.status == true
}

func (f *KnownUnknownFlag) IsFalse() bool {
	return f.status == false
}

func (f *KnownUnknownFlag) IsKnown() bool {
	return f.confidence
}
