package geometry

import (
	"encoding/json"
	"errors"
	"github.com/tidwall/gjson"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/geojson"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/utils"
)

func ToString(f geojson.Feature) (string, error) {

	geom := gjson.GetBytes(f.ToBytes(), "geometry")

	if !geom.Exists() {
		return "", errors.New("Missing geometry property")
	}

	body, err := json.Marshal(geom)

	if err != nil {
		return "", errors.New("Failed to serialize geometry property")
	}

	return string(body), nil
}

func Type(f geojson.Feature) string {

	possible := []string{
		"geometry.type",
	}

	return utils.StringProperty(f, possible, "unknown")
}
