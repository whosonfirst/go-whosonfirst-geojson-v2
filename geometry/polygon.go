package geometry

import (
	"encoding/json"
	"errors"
	"fmt"
	pm_geojson "github.com/paulmach/go.geojson"
	"github.com/skelterjohn/geom"
	"github.com/tidwall/gjson"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2"
	_ "log"
	_ "time"
)

type Polygon struct {
	geojson.Polygon `json:",omitempty"`
	Exterior        geom.Polygon   `json:"exterior"`
	Interior        []geom.Polygon `json:"interior"`
}

func (p Polygon) ExteriorRing() geom.Polygon {
	return p.Exterior
}

func (p Polygon) InteriorRings() []geom.Polygon {
	return p.Interior
}

func (p Polygon) ContainsCoord(c geom.Coord) bool {

	ext := p.ExteriorRing()

	contains := false

	if ext.ContainsCoord(c) {

		contains = true

		for _, int := range p.InteriorRings() {

			if int.ContainsCoord(c) {
				contains = false
				break
			}
		}
	}

	return contains
}

func GeometryForFeature(f geojson.Feature) (*geojson.Geometry, error) {

	// see notes below in PolygonsForFeature

	t := gjson.GetBytes(f.Bytes(), "geometry.type")

	if !t.Exists() {
		return nil, errors.New("Failed to determine geometry.type")
	}

	c := gjson.GetBytes(f.Bytes(), "geometry.coordinates")

	if !c.Exists() {
		return nil, errors.New("Failed to determine geometry.coordinates")
	}

	g := gjson.GetBytes(f.Bytes(), "geometry")

	var geom geojson.Geometry
	err := json.Unmarshal([]byte(g.Raw), &geom)

	if err != nil {
		return nil, err
	}

	return &geom, nil
}

func PolygonsForFeature(f geojson.Feature) ([]geojson.Polygon, error) {

	geom_rsp := gjson.GetBytes(f.Bytes(), "geometry")

	g, err := pm_geojson.UnmarshalGeometry([]byte(geom_rsp.String()))

	if err != nil {
		return nil, err
	}

	polys := make([]geojson.Polygon, 0)

	switch g.Type {

	case "LineString":

		exterior_ring := newRing(g.LineString)

		polygon := Polygon{
			Exterior: exterior_ring,
		}

		polys = []geojson.Polygon{polygon}

	case "Polygon":

		polygon := newPolygon(g.Polygon)
		polys = []geojson.Polygon{polygon}

	case "MultiPolygon":

		for _, poly := range g.MultiPolygon {
			polygon := newPolygon(poly)
			polys = append(polys, polygon)
		}

	case "Point":

		lat := g.Point[1]
		lon := g.Point[0]

		pt := []float64{
			lon,
			lat,
		}

		coords := [][]float64{
			pt, pt,
			pt, pt,
			pt,
		}

		exterior_ring := newRing(coords)

		if err != nil {
			return nil, err
		}

		interior_rings := make([]geom.Polygon, 0)

		polygon := Polygon{
			Exterior: exterior_ring,
			Interior: interior_rings,
		}

		polys = []geojson.Polygon{polygon}
		return polys, nil

	case "MultiPoint":

		exterior_ring := newRing(g.MultiPoint)

		polygon := Polygon{
			Exterior: exterior_ring,
		}

		polys = []geojson.Polygon{polygon}

	default:

		msg := fmt.Sprintf("Invalid geometry type '%s'", g.Type)
		return nil, errors.New(msg)
	}

	return polys, nil
}

func newRing(coords [][]float64) geom.Polygon {

	poly := geom.Polygon{}

	for _, pt := range coords {
		poly.AddVertex(geom.Coord{X: pt[0], Y: pt[1]})
	}

	return poly
}

func newPolygon(rings [][][]float64) Polygon {

	exterior := newRing(rings[0])
	interior := make([]geom.Polygon, 0)

	if len(rings) > 1 {

		for _, coords := range rings {
			interior = append(interior, newRing(coords))
		}
	}

	polygon := Polygon{
		Exterior: exterior,
		Interior: interior,
	}

	return polygon
}
