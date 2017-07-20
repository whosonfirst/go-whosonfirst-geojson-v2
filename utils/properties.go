package utils

import (
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/geojson"
)

func EnsureProperties(body []byte, properties []string) error {

	for _, path := range properties {

		r := gjson.GetBytes(body, path)

		if ! r.Exists(){
		   msg := fmt.Sprintf("Feature is missing a %s property", path)
		   return errors.New(msg)
		}
	}

	return nil
}

func Int64Property(f geojson.Feature, possible []string, d int64) int64 {

	for _, path := range possible {

		v := gjson.GetBytes(f.ToBytes(), path)

		if v.Exists() {
			return v.Int()
		}
	}

	return d
}

func StringProperty(f geojson.Feature, possible []string, d string) string {

	for _, path := range possible {

		v := gjson.GetBytes(f.ToBytes(), path)

		if v.Exists() {
			return v.String()
		}
	}

	return d
}
