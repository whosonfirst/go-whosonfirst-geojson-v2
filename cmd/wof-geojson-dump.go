package main

import (
	"flag"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/utils"
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

		fmt.Printf("# %s\n", path)

		fmt.Printf("ID is %d\n", f.Id())
		fmt.Printf("Name is %s\n", f.Name())
		fmt.Printf("Placetype is %s\n", f.Placetype())

		// fmt.Printf("Hierarchy is %s\n", f.Hierarchy())

		coord, _ := utils.NewCoordinateFromLatLons(0.0, 0.0)
		contains, _ := f.ContainsCoord(coord)

		fmt.Printf("Contains %v %t\n", coord, contains)

		bboxes, _ := f.BoundingBoxes()

		fmt.Printf("Count boxes %d\n", len(bboxes.Bounds()))
		fmt.Printf("MBR %s\n", bboxes.MBR())

	}

}
