package feature

import (
	"encoding/json"
	"github.com/tidwall/gjson"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/utils"
	"io/ioutil"
	"os"
)

func LoadFeatureFromFile(path string) (geojson.Feature, error) {

	body, err := UnmarshalFeatureFromFile(path)

	if err != nil {
		return nil, err
	}

	wofid := gjson.GetBytes(body, "properties.wof:id")

	if wofid.Exists() {
		return NewWOFFeature(body)
	}

	return NewGeoJSONFeature(body)
}

func LoadWOFFeatureFromFile(path string) (geojson.Feature, error) {

	body, err := UnmarshalFeatureFromFile(path)

	if err != nil {
		return nil, err
	}

	return NewWOFFeature(body)
}

func UnmarshalFeatureFromFile(path string) ([]byte, error) {

	fh, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer fh.Close()

	body, err := ioutil.ReadAll(fh)

	if err != nil {
		return nil, err
	}

	var stub interface{}
	err = json.Unmarshal(body, &stub)

	if err != nil {
		return nil, err
	}

	properties := []string{
		"geometry",
		"geometry.type",
		"geometry.coordinates",
		"properties",
		"id",
	}

	err = utils.EnsureProperties(body, properties)

	if err != nil {
		return nil, err
	}

	return body, nil
}
