package main

import (
	"flag"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/spatial"
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

		fmt.Println("Enspatialize bounding box")

		err = spatial.EnSpatialize(f)

		if err != nil {
			log.Fatal(err)
		}

	}

}
