package geometry

import (
	"github.com/skelterjohn/geom"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2"
)

type Bboxes struct {
	geojson.BoundingBoxes
	bounds []*geom.Rect
	mbr    geom.Rect
}

func (b Bboxes) Bounds() []*geom.Rect {
	return b.bounds
}

func (b Bboxes) MBR() geom.Rect {
	return b.mbr
}

func BoundingBoxesForFeature(f geojson.Feature) (geojson.BoundingBoxes, error) {

	polys, err := PolygonsForFeature(f)

	if err != nil {
		return nil, err
	}

	mbr := geom.NilRect()
	bounds := make([]*geom.Rect, 0)

	for _, poly := range polys {

		ext := poly.ExteriorRing()
		b := ext.Path.Bounds()

		mbr.ExpandToContainRect(*b)
		bounds = append(bounds, b)
	}

	wb := Bboxes{
		bounds: bounds,
		mbr:    mbr,
	}

	return wb, nil
}
