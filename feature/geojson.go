package feature

import (
	"encoding/json"
	"errors"
	"github.com/skelterjohn/geom"
	"github.com/tidwall/gjson"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/utils"
)

type GeoJSONBoundingBoxes struct {
	geojson.BoundingBoxes
	bounds []*geom.Rect
	mbr    geom.Rect
}

func (b GeoJSONBoundingBoxes) Bounds() []*geom.Rect {
	return b.bounds
}

func (b GeoJSONBoundingBoxes) MBR() geom.Rect {
	return b.mbr
}

type GeoJSONPolygon struct {
	geojson.Polygon
	exterior geom.Polygon
	interior []geom.Polygon
}

func (p GeoJSONPolygon) ExteriorRing() geom.Polygon {
	return p.exterior
}

func (p GeoJSONPolygon) InteriorRings() []geom.Polygon {
	return p.interior
}

func (p GeoJSONPolygon) ContainsCoord(c geom.Coord) bool {

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

type GeoJSONFeature struct {
	geojson.Feature
	body []byte
}

func NewGeoJSONFeature(body []byte) (geojson.Feature, error) {

	var stub interface{}
	err := json.Unmarshal(body, &stub)

	if err != nil {
		return nil, err
	}

	f := GeoJSONFeature{
		body: body,
	}

	return &f, nil
}

func (f *GeoJSONFeature) ContainsCoord(c geom.Coord) (bool, error) {

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

func (f *GeoJSONFeature) ToString() string {

	body, err := json.Marshal(f.body)

	if err != nil {
		return ""
	}

	return string(body)
}

func (f *GeoJSONFeature) ToBytes() []byte {

	return f.body
}

func (f *GeoJSONFeature) Id() string {

	possible := []string{
		"id",
		"properties.id",
	}

	return utils.StringProperty(f.ToBytes(), possible, "")
}

func (f *GeoJSONFeature) Name() string {

	possible := []string{
		"properties.name",
	}

	return utils.StringProperty(f.ToBytes(), possible, "")
}

func (f *GeoJSONFeature) Placetype() string {

	possible := []string{
		"properties.placetype",
	}

	return utils.StringProperty(f.ToBytes(), possible, "")
}

func (f *GeoJSONFeature) BoundingBoxes() (geojson.BoundingBoxes, error) {

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

	wb := GeoJSONBoundingBoxes{
		bounds: bounds,
		mbr:    mbr,
	}

	return wb, nil
}

func (f *GeoJSONFeature) Polygons() ([]geojson.Polygon, error) {

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

		polygon, err := f.gjson_coordsToGeoJSONPolygon(c)

		if err != nil {
			return nil, err
		}

		polys = append(polys, polygon)

	case "MultiPolygon":

		for _, rings := range coords {

			polygon, err := f.gjson_coordsToGeoJSONPolygon(rings)

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

func (f *GeoJSONFeature) gjson_coordsToGeoJSONPolygon(r gjson.Result) (geojson.Polygon, error) {

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

	polygon := GeoJSONPolygon{
		exterior: exterior,
		interior: interior,
	}

	return &polygon, nil
}

func (f *GeoJSONFeature) gjson_linearRingToGeomPolygon(r gjson.Result) (geom.Polygon, error) {

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
