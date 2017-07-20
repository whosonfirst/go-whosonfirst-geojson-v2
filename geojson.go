package geojson

import (
	"github.com/skelterjohn/geom"
)

type Feature interface {
	Id() string
	Name() string
	Placetype() string
	ToString() string
	ToBytes() []byte
	BoundingBoxes() (BoundingBoxes, error)
	Polygons() ([]Polygon, error)
	// SPR() (spr.StandardPlaceResponse, error)
	ContainsCoord(geom.Coord) (bool, error)
}

type BoundingBoxes interface {
	Bounds() []*geom.Rect
	MBR() geom.Rect
}

type Polygon interface {
	ExteriorRing() geom.Polygon
	InteriorRings() []geom.Polygon
	ContainsCoord(geom.Coord) bool
}
