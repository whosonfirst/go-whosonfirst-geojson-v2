package spatial

import (
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/geojson"
	"log"
)

func EnSpatialize(f geojson.Feature) error {

	/*
		id := wof.Id()
		name := wof.Name()
		placetype := wof.Placetype()
		deprecated := wof.Deprecated()
		superseded := wof.Superseded()
	*/

	bounds, err := f.Bounds()

	if err != nil {
		return err
	}

	for i, bbox := range bounds.BoundingBoxes() {

		sw := bbox.MinCoordinate()
		ne := bbox.MaxCoordinate()

		llat := ne.Latitude() - sw.Latitude()
		llon := ne.Longitude() - sw.Longitude()

		log.Println(i, llat, llon)

		// pt := rtreego.Point{swlon, swlat}
		// rect, err := rtreego.NewRect(pt, []float64{llon, llat})

		if err != nil {
			return err
		}
	}

	// return &WOFSpatial{rect, id, name, placetype, -1, deprecated, superseded}, nil

	return nil
}
