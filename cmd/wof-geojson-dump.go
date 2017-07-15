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

		fmt.Printf("# %s\n", path)

		fmt.Printf("ID is %d\n", f.Id())
		fmt.Printf("Name is %s\n", f.Name())
		fmt.Printf("Placetype is %s\n", f.Placetype())

		// fmt.Printf("Hierarchy is %s\n", f.Hierarchy())
	}

}
