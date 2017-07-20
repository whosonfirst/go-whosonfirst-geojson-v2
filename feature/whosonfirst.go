package feature

import (
	"encoding/json"
	"errors"
	"github.com/skelterjohn/geom"
	"github.com/tidwall/gjson"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/properties/geometry"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/properties/whosonfirst"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/utils"
	"strconv"
)

type WOFBoundingBoxes struct {
	geojson.BoundingBoxes
	bounds []*geom.Rect
	mbr    geom.Rect
}

func (b WOFBoundingBoxes) Bounds() []*geom.Rect {
	return b.bounds
}

func (b WOFBoundingBoxes) MBR() geom.Rect {
	return b.mbr
}

type WOFPolygon struct {
	geojson.Polygon
	exterior geom.Polygon
	interior []geom.Polygon
}

func (p WOFPolygon) ExteriorRing() geom.Polygon {
	return p.exterior
}

func (p WOFPolygon) InteriorRings() []geom.Polygon {
	return p.interior
}

func (p WOFPolygon) ContainsCoord(c geom.Coord) bool {

	ext := p.ExteriorRing()

	contains := false

	if ext.ContainsCoord(c) {

		contains = true

		for _, int := range p.InteriorRings() {

			if int.ContainsCoord(c) {
				contains = false
				break
			}
		}
	}

	return contains
}

type WOFFeature struct {
	geojson.Feature
	body []byte
}

func NewWOFFeature(body []byte) (geojson.Feature, error) {

	var stub interface{}
	err := json.Unmarshal(body, &stub)

	if err != nil {
		return nil, err
	}

	required := []string{
		"properties.wof:id",
		"properties.wof:name",
		"properties.wof:placetype",
	}

	err = utils.EnsureProperties(body, required)

	if err != nil {
		return nil, err
	}

	f := WOFFeature{
		body: body,
	}

	return &f, nil
}

func (f *WOFFeature) ContainsCoord(c geom.Coord) (bool, error) {

	polys, err := f.Polygons()

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

func (f *WOFFeature) ToString() string {

	body, err := json.Marshal(f.body)

	if err != nil {
		return ""
	}

	return string(body)
}

func (f *WOFFeature) ToBytes() []byte {

	return f.body
}

func (f *WOFFeature) Type() string {

	return geometry.Type(f)
}

func (f *WOFFeature) Id() string {

	id := whosonfirst.Id(f)
	return strconv.FormatInt(id, 10)
}

func (f *WOFFeature) Name() string {

	return whosonfirst.Name(f)
}

func (f *WOFFeature) Placetype() string {

	return whosonfirst.Placetype(f)
}

func (f *WOFFeature) BoundingBoxes() (geojson.BoundingBoxes, error) {

	polys, err := f.Polygons()

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

	wb := WOFBoundingBoxes{
		bounds: bounds,
		mbr:    mbr,
	}

	return wb, nil
}

func (f *WOFFeature) Polygons() ([]geojson.Polygon, error) {

	t := gjson.GetBytes(f.body, "geometry.type")

	if !t.Exists() {
		return nil, errors.New("Failed to determine geometry.type")
	}

	c := gjson.GetBytes(f.body, "geometry.coordinates")

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

		polygon, err := f.gjson_coordsToWOFPolygon(c)

		if err != nil {
			return nil, err
		}

		polys = append(polys, polygon)

	case "MultiPolygon":

		for _, rings := range coords {

			polygon, err := f.gjson_coordsToWOFPolygon(rings)

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

func (f *WOFFeature) gjson_coordsToWOFPolygon(r gjson.Result) (geojson.Polygon, error) {

	rings := r.Array()

	count_rings := len(rings)
	count_interior := count_rings - 1

	exterior, err := f.gjson_linearRingToGeomPolygon(rings[0])

	if err != nil {
		return nil, err
	}

	interior := make([]geom.Polygon, count_interior)

	for i := 1; i < count_interior; i++ {

		poly, err := f.gjson_linearRingToGeomPolygon(rings[i])

		if err != nil {
			return nil, err
		}

		interior = append(interior, poly)
	}

	polygon := WOFPolygon{
		exterior: exterior,
		interior: interior,
	}

	return &polygon, nil
}

func (f *WOFFeature) gjson_linearRingToGeomPolygon(r gjson.Result) (geom.Polygon, error) {

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
