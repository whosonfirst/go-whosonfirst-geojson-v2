package feature

import (
	"encoding/json"
	"github.com/tidwall/gjson"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/utils"
	"io"
	"io/ioutil"
	"os"
)

func LoadFeatureFromFile(path string) (geojson.Feature, error) {

	body, err := UnmarshalFeatureFromFile(path)

	if err != nil {
		return nil, err
	}

	return loadFeatureFromBytes(body)
}

func LoadFeatureFromReader(fh io.Reader) (geojson.Feature, error) {

	body, err := UnmarshalFeatureFromReader(fh)

	if err != nil {
		return nil, err
	}

	return loadFeatureFromBytes(body)
}

func loadFeatureFromBytes(body []byte) (geojson.Feature, error) {

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

func LoadWOFFeatureFromReader(fh io.Reader) (geojson.Feature, error) {

	body, err := UnmarshalFeatureFromReader(fh)

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

	return UnmarshalFeatureFromReader(fh)
}

func UnmarshalFeatureFromReader(fh io.Reader) ([]byte, error) {

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
