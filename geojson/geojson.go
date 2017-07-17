package geojson

import (
	"github.com/skelterjohn/geom"
)

type Feature interface {
	Type() string
	Id() int64
	Name() string
	Placetype() string
	ToString() string
	ToBytes() []byte
	// Bounds() (Bounds, error)
	Polygons() ([]Polygon, error)
	// SPR() (StandardPlaceResponse, error)
}

type StandardPlaceResponse interface {
}

type Polygon interface {
	ExteriorRing() geom.Polygon
	InteriorRings() [][]geom.Polygon
}
