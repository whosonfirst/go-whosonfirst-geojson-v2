package geometry

// not sure about this package name yet... (20170722/thisisaaronland)

import (
	"errors"
	"github.com/skelterjohn/geom"
	"github.com/tidwall/gjson"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/utils"
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

type Polygon struct {
	geojson.Polygon
	exterior geom.Polygon
	interior []geom.Polygon
}

func (p Polygon) ExteriorRing() geom.Polygon {
	return p.exterior
}

func (p Polygon) InteriorRings() []geom.Polygon {
	return p.interior
}

func ContainsCoord(f geojson.Feature, c geom.Coord) (bool, error) {

	polys, err := Polygons(f)

	if err != nil {
		return false, err
	}

	contains := false

	for _, p := range polys {

		if p.ContainsCoord(c) {
			contains = true
			break
		}
	}

	return contains, nil
}

func BoundingBoxes(f geojson.Feature) (geojson.BoundingBoxes, error) {

	polys, err := Polygons(f)

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

func Polygons(f geojson.Feature) ([]geojson.Polygon, error) {

	t := gjson.GetBytes(f.ToBytes(), "geometry.type")

	if !t.Exists() {
		return nil, errors.New("Failed to determine geometry.type")
	}

	c := gjson.GetBytes(f.ToBytes(), "geometry.coordinates")

	if !c.Exists() {
		return nil, errors.New("Failed to determine geometry.coordinates")
	}

	coords := c.Array()

	if len(coords) == 0 {
		return nil, errors.New("Invalid geometry.coordinates")
	}

	polys := make([]geojson.Polygon, 0)

	switch t.String() {

	case "Polygon":

		// c === rings (below)

		polygon, err := gjson_coordsToPolygon(c)

		if err != nil {
			return nil, err
		}

		polys = append(polys, polygon)

	case "MultiPolygon":

		for _, rings := range coords {

			polygon, err := gjson_coordsToPolygon(rings)

			if err != nil {
				return nil, err
			}

			polys = append(polys, polygon)

		}

	default:

		return nil, errors.New("Invalid geometry type")
	}

	return polys, nil
}

func gjson_coordsToPolygon(r gjson.Result) (geojson.Polygon, error) {

	rings := r.Array()

	count_rings := len(rings)
	count_interior := count_rings - 1

	exterior, err := gjson_linearRingToGeomPolygon(rings[0])

	if err != nil {
		return nil, err
	}

	interior := make([]geom.Polygon, count_interior)

	for i := 1; i < count_interior; i++ {

		poly, err := gjson_linearRingToGeomPolygon(rings[i])

		if err != nil {
			return nil, err
		}

		interior = append(interior, poly)
	}

	polygon := Polygon{
		exterior: exterior,
		interior: interior,
	}

	return &polygon, nil
}

func gjson_linearRingToGeomPolygon(r gjson.Result) (geom.Polygon, error) {

	coords := make([]geom.Coord, 0)

	for _, pt := range r.Array() {

		lonlat := pt.Array()

		lat := lonlat[1].Float()
		lon := lonlat[0].Float()

		coord, _ := utils.NewCoordinateFromLatLons(lat, lon)
		coords = append(coords, coord)
	}

	return utils.NewPolygonFromCoords(coords)
}
