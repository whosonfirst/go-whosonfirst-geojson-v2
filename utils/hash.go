package utils

import (
       	"github.com/whosonfirst/go-whosonfirst-geojson-v2"
       	// "github.com/whosonfirst/go-whosonfirst-geojson-v2/properties/geometry"
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

// this causes an import loop so we're just going to leave it
// here as a reference for now... (20170801/thisisaaronland)
// HashGeometryForFeature(f geojson.Feature) (string, error)
// geom, err := geometry.ToString(f)
// return HashGeometry([]byte(geom))

func HashGeometry(geom []byte) (string, error) {

	h, err := hash.NewHash(hash_algo)

	if err != nil {
		return "", err
	}

	return h.HashFromJSON(geom)
}
