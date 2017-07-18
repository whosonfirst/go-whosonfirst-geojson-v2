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
	BoundingBoxes() (BoundingBoxes, error)
	Polygons() ([]Polygon, error)
	// Centroids() ([]Centroid, eror)
	// SPR() (spr.StandardPlaceResponse, error)
	ContainsCoord(geom.Coord) (bool, error)
}

// type Centroid interface {
//     Source sources.WOFSource
//     Coord geom.Coord
// }

type BoundingBoxes interface {
	Bounds() []*geom.Rect
	MBR() geom.Rect
}

type Polygon interface {
	ExteriorRing() geom.Polygon
	InteriorRings() []geom.Polygon
	ContainsCoord(geom.Coord) bool
}
