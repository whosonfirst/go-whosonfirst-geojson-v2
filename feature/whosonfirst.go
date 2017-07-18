package feature

import (
	"encoding/json"
	"errors"
	"github.com/skelterjohn/geom"
	"github.com/tidwall/gjson"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/geojson"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/utils"
)

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

	possible := []string{
		"geometry.type",
	}

	return f.possibleString(possible, "unknown")
}

func (f *WOFFeature) Id() int64 {

	possible := []string{
		"properties.f:id",
		"id",
	}

	return f.possibleInt64(possible, -1)
}

func (f *WOFFeature) Name() string {

	possible := []string{
		"properties.wof:name",
		"properties.name",
	}

	return f.possibleString(possible, "a place with no name")
}

func (f *WOFFeature) Placetype() string {

	possible := []string{
		"properties.wof:placetype",
		"properties.placetype",
	}

	return f.possibleString(possible, "here be dragons")
}

func (f *WOFFeature) Hierarchy() []map[string]int64 {

	hierarchies := make([]map[string]int64, 0)

	possible := gjson.GetBytes(f.body, "properties.wof:hierarchy")

	if possible.Exists() {

		for _, h := range possible.Array() {

			foo := make(map[string]int64)

			for k, v := range h.Map() {

				foo[k] = v.Int()
			}

			hierarchies = append(hierarchies, foo)
		}
	}

	return hierarchies
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

func (f *WOFFeature) IsCurrent() (bool, bool) {

	possible := []string{
		"properties.mz_iscurrent",
	}

	v := f.possibleInt64(possible, -1)

	if v == 1 {
		return true, true
	}

	if v == 0 {
		return true, false
	}

	if f.IsDeprecated() {
		return true, false
	}

	if f.IsSuperseded() {
		return true, false
	}

	return false, false
}

func (f *WOFFeature) IsDeprecated() bool {

	possible := []string{
		"properties.edtf:deprecated",
	}

	v := f.possibleString(possible, "uuuu")

	if v != "" && v != "u" && v != "uuuu" {
		return true
	}

	return false
}

func (f *WOFFeature) IsSuperseded() bool {

	possible := []string{
		"properties.edtf:superseded",
	}

	v := f.possibleString(possible, "uuuu")

	if v != "" && v != "u" && v != "uuuu" {
		return true
	}

	by := gjson.GetBytes(f.body, "properties.wof:superseded_by")

	if by.Exists() && len(by.Array()) > 0 {
		return true
	}

	return false
}

func (f *WOFFeature) possibleInt64(possible []string, d int64) int64 {

	for _, path := range possible {

		v := gjson.GetBytes(f.body, path)

		if v.Exists() {
			return v.Int()
		}
	}

	return d
}

func (f *WOFFeature) possibleString(possible []string, d string) string {

	for _, path := range possible {

		v := gjson.GetBytes(f.body, path)

		if v.Exists() {
			return v.String()
		}
	}

	return d
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
