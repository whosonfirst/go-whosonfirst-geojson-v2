package main

import (
	"flag"
	"github.com/whosonfirst/go-whosonfirst-placetypes"
	"github.com/whosonfirst/go-whosonfirst-cli/flags"	
	"log"
)

func main() {

	var roles flags.MultiString
	flag.Var(&roles, "role", "...")
	
	flag.Parse()
	
	for _, str_pt := range flag.Args() {
	
		pt, err := placetypes.GetPlacetypeByName(str_pt)

		if err != nil {
			log.Fatal(err)
		}

		children := placetypes.Children(pt)
		
		for i, p := range children {
			log.Println(i, p.Name)
		}
	}
}
