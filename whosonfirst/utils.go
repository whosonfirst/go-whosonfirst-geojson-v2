package whosonfirst

import (
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/feature"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/geojson"
	"io/ioutil"
	"os"
)

func LoadFeatureFromFile(path string) (geojson.Feature, error) {

	fh, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer fh.Close()

	body, err := ioutil.ReadAll(fh)

	if err != nil {
		return nil, err
	}

	return feature.NewWOFFeature(body)
}
