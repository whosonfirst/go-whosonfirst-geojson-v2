package spatial

import (
	"errors"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/geojson"
)

func ContainsPoint(f geojson.Feature, pt geojson.Coordinate) (bool, error) {

	switch f.Type() {

	case "Point":
		return false, errors.New("Please write me")
	case "Polygon":
		return false, errors.New("Please write me")
	case "MultiPolygon":
		return false, errors.New("Please write me")
	default:
		return false, errors.New("Unsupported geometry type")
	}

}
