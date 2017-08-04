package feature

import (
	"encoding/json"
	"errors"
	"github.com/skelterjohn/geom"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/geometry"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/utils"
	"github.com/whosonfirst/go-whosonfirst-spr"
)

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

	return geometry.FeatureContainsCoord(f, c)
}

func (f *GeoJSONFeature) String() string {

	body, err := json.Marshal(f.body)

	if err != nil {
		return ""
	}

	return string(body)
}

func (f *GeoJSONFeature) Bytes() []byte {

	return f.body
}

func (f *GeoJSONFeature) Id() string {

	possible := []string{
		"id",
		"properties.id",
	}

	return utils.StringProperty(f.Bytes(), possible, "")
}

func (f *GeoJSONFeature) Name() string {

	possible := []string{
		"properties.name",
	}

	return utils.StringProperty(f.Bytes(), possible, "")
}

func (f *GeoJSONFeature) Placetype() string {

	possible := []string{
		"properties.placetype",
	}

	return utils.StringProperty(f.Bytes(), possible, "")
}

func (f *GeoJSONFeature) BoundingBoxes() (geojson.BoundingBoxes, error) {
	return geometry.BoundingBoxesForFeature(f)
}

func (f *GeoJSONFeature) Polygons() ([]geojson.Polygon, error) {
	return geometry.PolygonsForFeature(f)
}

func (f *GeoJSONFeature) SPR() (spr.StandardPlacesResult, error) {
	return nil, errors.New("SPR is not implemented yet.")
}
