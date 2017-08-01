package feature

import (
	"encoding/json"
	"errors"
	"github.com/skelterjohn/geom"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/geometry"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/properties/whosonfirst"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/utils"
	"github.com/whosonfirst/go-whosonfirst-placetypes"
	"strconv"
)

type WOFFeature struct {
	geojson.Feature
	body []byte
}

func EnsureWOFFeature(body []byte) error {

	required := []string{
		"properties.wof:id",
		"properties.wof:name",
		"properties.wof:repo",
		"properties.wof:placetype",
	}

	err := utils.EnsureProperties(body, required)

	if err != nil {
		return err
	}

	pt := utils.StringProperty(body, []string{"properties.wof:placetype"}, "")

	if !placetypes.IsValidPlacetype(pt) {
		return errors.New("Invalid wof:placetype")
	}

	// check wof:repo here?

	return nil
}

func NewWOFFeature(body []byte) (geojson.Feature, error) {

	var stub interface{}
	err := json.Unmarshal(body, &stub)

	if err != nil {
		return nil, err
	}

	err = EnsureWOFFeature(body)

	if err != nil {
		return nil, err
	}

	f := WOFFeature{
		body: body,
	}

	return &f, nil
}

func (f *WOFFeature) String() string {

	body, err := json.Marshal(f.body)

	if err != nil {
		return ""
	}

	return string(body)
}

func (f *WOFFeature) Bytes() []byte {
	return f.body
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
	return geometry.BoundingBoxesForFeature(f)
}

func (f *WOFFeature) Polygons() ([]geojson.Polygon, error) {
	return geometry.PolygonsForFeature(f)
}

func (f *WOFFeature) ContainsCoord(c geom.Coord) (bool, error) {
	return geometry.FeatureContainsCoord(f, c)
}
