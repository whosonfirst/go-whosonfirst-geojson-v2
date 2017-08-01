package utils

import (
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/properties/geometry"
	"github.com/whosonfirst/go-whosonfirst-hash"
)

const hash_algo = "md5"

func HashFeature(f geojson.Feature) (string, error) {

	h, err := hash.NewHash(hash_algo)

	if err != nil {
		return "", err
	}

	return h.HashFromJSON(f.Bytes())
}

func HashGeometryForFeature(f geojson.Feature) (string, error) {

	geom, err := geometry.ToString(f)

	if err != nil {
		return "", err
	}

	return HashGeomtry([]byte(geom))
}

func HashGeometry(geom []byte) (string, error) {

	h, err := hash.NewHash(hash_algo)

	if err != nil {
		return "", err
	}

	return h.HashFromJSON(f.Bytes())
}
