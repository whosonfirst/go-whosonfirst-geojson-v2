package main

import (
	"flag"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/geojson"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/spatial"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/whosonfirst"
	"log"
)

type Point struct {
	geojson.Coordinate
	lat float64
	lon float64
}

func (p Point) Latitude() float64 {
	return p.lat
}

func (p Point) Longitude() float64 {
	return p.lon
}

func main() {

	flag.Parse()
	args := flag.Args()

	lat := 45.523668
	lon := -73.600159

	pt := Point{
		lat: lat,
		lon: lon,
	}

	for _, path := range args {

		f, err := whosonfirst.LoadFeatureFromFile(path)

		if err != nil {
			log.Fatal(err)
		}

		contains, err := spatial.ContainsPoint(f, pt)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%s contains point: %t\n", path, contains)
	}

}
