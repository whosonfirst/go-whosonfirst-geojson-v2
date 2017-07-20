package geometry

import (
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/geojson"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/utils"
)

func Type(f geojson.Feature) string {

	possible := []string{
		"geometry.type",
	}

	return utils.StringProperty(f, possible, "unknown")
}
