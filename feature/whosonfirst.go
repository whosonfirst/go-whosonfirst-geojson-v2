package feature

import (
	"encoding/json"
	"errors"
	"github.com/tidwall/gjson"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/geojson"
)

type WOFCoordinate struct {
	geojson.Coordinate
	lat float64
	lon float64
}

type WOFBoundingBox struct {
	geojson.BoundingBox
	sw geojson.Coordinate
	ne geojson.Coordinate
}

type WOFBounds struct {
	geojson.Bounds
	bboxes []geojson.BoundingBox
}

type WOFRing struct {
	geojson.Ring
	coords []geojson.Coordinate
}

type WOFPolygon struct {
	exterior geojson.Ring
	interior []geojson.Ring
}

type WOFFeature struct {
	geojson.Feature
	body []byte
}

func NewWOFCoordinateFromLatLons(lat float64, lon float64) (geojson.Coordinate, error) {

	coord := WOFCoordinate{
		lat: lat,
		lon: lon,
	}

	return &coord, nil
}

func NewWOFBoundingBoxFromLatLons(minlat float64, minlon float64, maxlat float64, maxlon float64) (geojson.BoundingBox, error) {

	min_coord, err := NewWOFCoordinateFromLatLons(minlat, minlon)

	if err != nil {
		return nil, err
	}

	max_coord, err := NewWOFCoordinateFromLatLons(maxlat, maxlon)

	if err != nil {
		return nil, err
	}

	bbox := WOFBoundingBox{
		sw: min_coord,
		ne: max_coord,
	}

	return &bbox, nil
}

func NewWOFBoundsFromBoundingBoxes(bboxes []geojson.BoundingBox) (*WOFBounds, error) {

	bounds := WOFBounds{
		bboxes: bboxes,
	}

	return &bounds, nil
}

func NewWOFRingFromCoords(coords []geojson.Coordinate) (geojson.Ring, error) {

	r := WOFRing{
		coords: coords,
	}

	return &r, nil
}

func (c WOFCoordinate) Latitude() float64 {
	return c.lat
}

func (c WOFCoordinate) Longitude() float64 {
	return c.lon
}

func (b WOFBoundingBox) MinCoordinate() geojson.Coordinate {
	return b.sw
}

func (b WOFBoundingBox) MaxCoordinate() geojson.Coordinate {
	return b.ne
}

func (b WOFBounds) MBR() geojson.BoundingBox {

	minlat := 0.0
	minlon := 0.0
	maxlat := 0.0
	maxlon := 0.0

	for _, bbox := range b.BoundingBoxes() {

		min := bbox.MinCoordinate()
		max := bbox.MaxCoordinate()

		if min.Latitude() < minlat {
			minlat = min.Latitude()
		}

		if min.Longitude() < minlon {
			minlon = min.Longitude()
		}

		if max.Latitude() > maxlat {
			maxlat = max.Latitude()
		}

		if max.Longitude() > maxlon {
			maxlon = max.Longitude()
		}
	}

	min_coord := WOFCoordinate{
		lat: minlat,
		lon: minlon,
	}

	max_coord := WOFCoordinate{
		lat: maxlat,
		lon: maxlon,
	}

	bbox := WOFBoundingBox{
		sw: &min_coord,
		ne: &max_coord,
	}

	return &bbox
}

func (b WOFBounds) BoundingBoxes() []geojson.BoundingBox {
	return b.bboxes
}

func (r WOFRing) Coordinates() []geojson.Coordinate {
	return r.coords
}

func (p WOFPolygon) ExteriorRing() geojson.Ring {
	return p.exterior
}

func (p WOFPolygon) InteriorRings() []geojson.Ring {
	return p.interior
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

	return nil, errors.New("Please write me")
}

func (f *WOFFeature) Bounds() (geojson.Bounds, error) {

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

	boxes := make([]geojson.BoundingBox, 0)

	if t.String() == "Point" {

		lat := coords[1].Float()
		lon := coords[0].Float()

		bbox, err := NewWOFBoundingBoxFromLatLons(lat, lon, lat, lon)

		if err != nil {
			return nil, err
		}

		boxes = append(boxes, bbox)
	}

	if t.String() == "Polygon" {

		outer := coords[0]
		bbox, err := f.gjson_linearRingToBoundingBox(outer)

		if err != nil {
			return nil, err
		}

		boxes = append(boxes, bbox)
	}

	if t.String() == "MultiPolygon" {

		for _, rings := range coords {

			poly := rings.Array()

			outer := poly[0]
			bbox, err := f.gjson_linearRingToBoundingBox(outer)

			if err != nil {
				return nil, err
			}

			boxes = append(boxes, bbox)
		}
	}

	return NewWOFBoundsFromBoundingBoxes(boxes)
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

func (f *WOFFeature) gjson_linearRingToBoundingBox(ring gjson.Result) (geojson.BoundingBox, error) {

	minlat := 0.0
	minlon := 0.0
	maxlat := 0.0
	maxlon := 0.0

	for _, pt := range ring.Array() {

		lonlat := pt.Array()

		lat := lonlat[1].Float()
		lon := lonlat[0].Float()

		if lat < minlat {
			minlat = lat
		}

		if lon < minlon {
			minlon = lon
		}

		if lat > maxlat {
			maxlat = lat
		}

		if lon > maxlon {
			maxlon = lon
		}
	}

	return NewWOFBoundingBoxFromLatLons(minlat, minlon, maxlat, maxlon)
}

func (f *WOFFeature) gjson_linearRingToWOFRing(ring gjson.Result) (geojson.Ring, error) {

	coords := make([]geojson.Coordinate, 0)

	for _, pt := range ring.Array() {

		lonlat := pt.Array()

		lat := lonlat[1].Float()
		lon := lonlat[0].Float()

		coord, err := NewWOFCoordinateFromLatLons(lat, lon)

		if err != nil {
			return nil, err
		}

		coords = append(coords, coord)
	}

	return NewWOFRingFromCoords(coords)
}
