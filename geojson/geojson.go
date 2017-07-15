package geojson

import ()

type Feature interface {
	Type() string
	Id() int64
	Name() string
	Placetype() string
	ToString() string
	ToBytes() []byte
	Bounds() (Bounds, error)
	Polygons() ([]Polygon, error)
}

type Bounds interface {
	MBR() BoundingBox
	BoundingBoxes() []BoundingBox
}

type Coordinate interface {
	Latitude() float64
	Longitude() float64
}

type BoundingBox interface {
	MinCoordinate() Coordinate
	MaxCoordinate() Coordinate
}

type Ring interface {
	Coordinates() []Coordinate
}

type Polygon interface {
	ExteriorRing() Ring
	InteriorRings() [][]Ring
}
