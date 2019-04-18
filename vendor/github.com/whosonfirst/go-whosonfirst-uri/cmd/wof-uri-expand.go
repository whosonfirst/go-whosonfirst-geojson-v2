package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-uri"
	"strconv"
	"log"
	"path/filepath"
	"os"
)

func ExpandString(str_id string) (string, error) {

	id, err := strconv.ParseInt(str_id, 10, 64)

	if err != nil {
		return "", err
	}

	return uri.Id2RelPath(id)	
}

func main() {

	root := flag.String("root", "", "An optional (filesystem) root to prepend URIs with")
	stdin := flag.Bool("stdin", false, "Read IDs from STDINT")
	
	flag.Parse()

	if *root != "" {

		abs_root, err := filepath.Abs(*root)

		if err != nil {
			log.Fatal(err)
		}

		*root = abs_root
	}

	expand := func(str_id string) {
		
		path, err := ExpandString(str_id)

		if err != nil {
			log.Fatal(err)
		}
		
		if *root != "" {
			path = filepath.Join(*root, path)
		}

		fmt.Println(path)
	}
	
	if *stdin {

		scanner := bufio.NewScanner(os.Stdin)

		for scanner.Scan() {
			expand(scanner.Text())
		}

		err := scanner.Err()

		if err != nil {
			log.Fatal(err)
		}
		
	} else {
	
		for _, str_id := range flag.Args(){
			expand(str_id)
		}
	}
	
	os.Exit(0)
}
	
