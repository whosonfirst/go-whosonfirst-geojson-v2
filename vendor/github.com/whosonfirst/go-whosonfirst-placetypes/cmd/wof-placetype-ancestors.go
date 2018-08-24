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

		var ancestors []*placetypes.WOFPlacetype
		
		if len(roles) > 0 {
			ancestors = placetypes.AncestorsForRoles(pt, roles)
		} else {
			ancestors = placetypes.Ancestors(pt)
		}
		
		for i, p := range ancestors {
			log.Println(i, p.Name)
		}
	}
}
