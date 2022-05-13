# go-whosonfirst-geojson-v2

Go tools for working with Who's On First documents

## Notice

This package is officially. It is still being used in other packages but they are slowly being updated.

Basically this package tries to do too many things, specifically around defining GeoJSON-related structs and interfaces. It's not really worth the effort and better to use [paulmach/orb/geojson](https://github.com/paulmach/orb) for geometry and GeoJSON-related operations and [tidwall/gjson](https://github.com/tidwall/gjson] for query-related operations using plain-vanilla `[]byte` elements. This is the approach taken by the `go-whosonfirst-feature` package.

If you are using _this package_ in your code it would be best to migrate it to use equivalent functionality defined in the `go-whosonfirst-feature` package.

## History

The goal of this package was to replace the existing [go-whosonfirst-geojson](https://github.com/whosonfirst/go-whosonfirst-geojson) package.

### geojson-v2?

Yeah, I don't really like it either but this package is basically 100% backwards incompatible with `github.com/whosonfirst/go-whosonfirst-geojson` and while I don't _really_ think anyone else is using it I don't like the idea of suddenly breaking everyone's code.

## Interfaces

Unlike the first `go-whosonfirst-geojson` package this one at least attempts to define a simplified interface for working with GeoJSON features. These are still in flux.

_Please finish writing me._

### geojson.Feature

```
type Feature interface {
	Type() string
	Id() int64
	Name() string
	Placetype() string
	ToString() string
	ToBytes() []byte
	BoundingBoxes() (BoundingBoxes, error)
	Polygons() ([]Polygon, error)
	ContainsCoord(geom.Coord) (bool, error)
}
```

### geojson.BoundingBoxes

```
type BoundingBoxes interface {
	Bounds() []*geom.Rect
	MBR() geom.Rect
}
```

### geojson.Centroid

```
type Centroid interface {
	Coord() geom.Coord
	Source() string
}
```

### geojson.Polygon

```
type Polygon interface {
	ExteriorRing() geom.Polygon
	InteriorRings() []geom.Polygon
	ContainsCoord(geom.Coord) bool
}
```

## Usage

### Simple

```
import (
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/whosonfirst"
	"log"
)

func main() {

	path := "/usr/local/data/whosonfirst-data/data/101/736/545/101736545.geojson"
	f, err := whosonfirst.LoadFeatureFromFile(path)

	if err != nil {
		log.Fatal(err)
	}

	// prints "Montreal"
	log.Println("Name is ", f.Name())
}
```

## See also

* github.com/skelterjohn/geom
* https://github.com/whosonfirst/go-whosonfirst-geojson

