package whosonfirst

import (
	"github.com/tidwall/gjson"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/utils"
)

func Centroid(f geojson.Feature) (float64, float64) {

	var lat gjson.Result
	var lon gjson.Result

	lat = gjson.GetBytes(f.Bytes(), "properties.lbl:latitude")
	lon = gjson.GetBytes(f.Bytes(), "properties.lbl:longitude")

	if lat.Exists() && lon.Exists() {
		return lat.Float(), lon.Float()
	}

	lat = gjson.GetBytes(f.Bytes(), "properties.geom:latitude")
	lon = gjson.GetBytes(f.Bytes(), "properties.geom:longitude")

	if lat.Exists() && lon.Exists() {
		return lat.Float(), lon.Float()
	}

	return 0.0, 0.0
}

func Country(f geojson.Feature) string {

	possible := []string{
		"properties.wof:country",
	}

	return utils.StringProperty(f.Bytes(), possible, "XX")
}

func Id(f geojson.Feature) int64 {

	possible := []string{
		"properties.wof:id",
		"id",
	}

	return utils.Int64Property(f.Bytes(), possible, -1)
}

func Name(f geojson.Feature) string {

	possible := []string{
		"properties.wof:name",
		"properties.name",
	}

	return utils.StringProperty(f.Bytes(), possible, "a place with no name")
}

func ParentId(f geojson.Feature) int64 {

	possible := []string{
		"properties.wof:parent_id",
	}

	return utils.Int64Property(f.Bytes(), possible, -1)
}

func Placetype(f geojson.Feature) string {

	possible := []string{
		"properties.wof:placetype",
		"properties.placetype",
	}

	return utils.StringProperty(f.Bytes(), possible, "here be dragons")
}

func Repo(f geojson.Feature) string {

	possible := []string{
		"properties.wof:repo",
	}

	return utils.StringProperty(f.Bytes(), possible, "whosonfirst-data-xx")
}

func IsCurrent(f geojson.Feature) (bool, bool) {

	possible := []string{
		"properties.mz_iscurrent",
	}

	v := utils.Int64Property(f.Bytes(), possible, -1)

	if v == 1 {
		return true, true
	}

	if v == 0 {
		return false, true
	}

	if IsDeprecated(f) {
		return false, true
	}

	if IsCeased(f) {
		return false, true
	}

	if IsSuperseded(f) {
		return false, true
	}

	// as in -1

	return false, false
}

func IsDeprecated(f geojson.Feature) bool {

	possible := []string{
		"properties.edtf:deprecated",
	}

	v := utils.StringProperty(f.Bytes(), possible, "uuuu")

	if v != "" && v != "u" && v != "uuuu" {
		return true
	}

	return false
}

func IsCeased(f geojson.Feature) bool {

	possible := []string{
		"properties.edtf:cessation",
	}

	v := utils.StringProperty(f.Bytes(), possible, "uuuu")

	if v != "" && v != "u" && v != "uuuu" {
		return true
	}

	return false
}

func IsSuperseded(f geojson.Feature) bool {

	possible := []string{
		"properties.edtf:superseded",
	}

	v := utils.StringProperty(f.Bytes(), possible, "uuuu")

	if v != "" && v != "u" && v != "uuuu" {
		return true
	}

	by := gjson.GetBytes(f.Bytes(), "properties.wof:superseded_by")

	if by.Exists() && len(by.Array()) > 0 {
		return true
	}

	return false
}

func IsSuperseding(f geojson.Feature) bool {

	sc := gjson.GetBytes(f.Bytes(), "properties.wof:supersedes")

	if sc.Exists() && len(sc.Array()) > 0 {
		return true
	}

	return false
}

func SupersededBy(f geojson.Feature) []int64 {

	superseded_by := make([]int64, 0)

	possible := gjson.GetBytes(f.Bytes(), "properties.wof:superseded_by")

	if possible.Exists() {

		for _, id := range possible.Array() {
			superseded_by = append(superseded_by, id.Int())
		}
	}

	return superseded_by
}

func Supersedes(f geojson.Feature) []int64 {

	supersedes := make([]int64, 0)

	possible := gjson.GetBytes(f.Bytes(), "properties.wof:supersedes")

	if possible.Exists() {

		for _, id := range possible.Array() {
			supersedes = append(supersedes, id.Int())
		}
	}

	return supersedes
}

func Hierarchy(f geojson.Feature) []map[string]int64 {

	hierarchies := make([]map[string]int64, 0)

	possible := gjson.GetBytes(f.Bytes(), "properties.wof:hierarchy")

	if possible.Exists() {

		for _, h := range possible.Array() {

			foo := make(map[string]int64)

			for k, v := range h.Map() {

				foo[k] = v.Int()
			}

			hierarchies = append(hierarchies, foo)
		}
	}

	return hierarchies
}
