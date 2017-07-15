package main

import (
	"flag"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/whosonfirst"
	"log"
)

func main() {

	flag.Parse()
	args := flag.Args()

	for _, path := range args {

		f, err := whosonfirst.LoadFeatureFromFile(path)

		if err != nil {
			log.Fatal(err)
		}

		polys, err := f.Polygons()

		if err != nil {
			log.Fatal(err)
		}

		for _, p := range polys {
			ext := p.ExteriorRing()
			fmt.Printf("%d points\n", len(ext.Coordinates()))
		}
	}

}
