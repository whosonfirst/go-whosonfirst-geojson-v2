package feature

import (
	"encoding/json"
	"errors"
	"github.com/skelterjohn/geom"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/geometry"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/properties/whosonfirst"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/utils"
	"github.com/whosonfirst/go-whosonfirst-placetypes"
	"github.com/whosonfirst/go-whosonfirst-spr"
	"github.com/whosonfirst/go-whosonfirst-uri"
	"strconv"
)

type WOFFeature struct {
	geojson.Feature
	body []byte
}

type WOFExistentialFlag struct {
     spr.ExistentialFlag
     WOFStatusInt int64
     WOFStatus bool
     WOFConfidence bool 
     raw interface{}
}

func (f *WOFExistentialFlag) StatusInt() int64 {
     return f.WOFStatusInt
}

func (f *WOFExistentialFlag) Status() bool {
     return f.WOFStatus
}

func (f *WOFExistentialFlag) Confidence() bool {
     return f.WOFConfidence
}

func (f *WOFExistentialFlag) String() string {
     return fmt.Sprintf("%v", f.raw)
}

type WOFStandardPlacesResult struct {
	spr.StandardPlacesResult `json:",omitempty"`
	WOFId                    int64   `json:"wof:id"`
	WOFParentId              int64   `json:"wof:parent_id"`
	WOFName                  string  `json:"wof:name"`
	WOFPlacetype             string  `json:"wof:placetype"`
	WOFCountry               string  `json:"wof:country"`
	WOFRepo                  string  `json:"wof:repo"`
	WOFPath                  string  `json:"wof:path"`
	MZURI                    string  `json:"mz:uri"`
	WOFSupersededBy          []int64 `json:"wof:superseded_by"`
	WOFSupersedes            []int64 `json:"wof:supersedes"`
	MZIsCurrent              int     `json:"mz:is_current"`
	MZIsCeased               int     `json:"mz:is_ceased"`
	MZIsDeprecated           int     `json:"mz:is_deprecated"`
	MZIsSuperseded           int     `json:"mz:is_superseded"`
	MZIsSuperseding          int     `json:"mz:is_superseding"`
}

func EnsureWOFFeature(body []byte) error {

	required := []string{
		"properties.wof:id",
		"properties.wof:name",
		"properties.wof:repo",
		"properties.wof:placetype",
	}

	err := utils.EnsureProperties(body, required)

	if err != nil {
		return err
	}

	pt := utils.StringProperty(body, []string{"properties.wof:placetype"}, "")

	if !placetypes.IsValidPlacetype(pt) {
		return errors.New("Invalid wof:placetype")
	}

	// check wof:repo here?

	return nil
}

func NewWOFFeature(body []byte) (geojson.Feature, error) {

	var stub interface{}
	err := json.Unmarshal(body, &stub)

	if err != nil {
		return nil, err
	}

	err = EnsureWOFFeature(body)

	if err != nil {
		return nil, err
	}

	f := WOFFeature{
		body: body,
	}

	return &f, nil
}

func (f *WOFFeature) String() string {

	body, err := json.Marshal(f.body)

	if err != nil {
		return ""
	}

	return string(body)
}

func (f *WOFFeature) Bytes() []byte {
	return f.body
}

func (f *WOFFeature) Id() string {
	id := whosonfirst.Id(f)
	return strconv.FormatInt(id, 10)
}

func (f *WOFFeature) Name() string {
	return whosonfirst.Name(f)
}

func (f *WOFFeature) Placetype() string {
	return whosonfirst.Placetype(f)
}

func (f *WOFFeature) BoundingBoxes() (geojson.BoundingBoxes, error) {
	return geometry.BoundingBoxesForFeature(f)
}

func (f *WOFFeature) Polygons() ([]geojson.Polygon, error) {
	return geometry.PolygonsForFeature(f)
}

func (f *WOFFeature) ContainsCoord(c geom.Coord) (bool, error) {
	return geometry.FeatureContainsCoord(f, c)
}

func (f *WOFFeature) SPR() (spr.StandardPlacesResult, error) {

	id := whosonfirst.Id(f)
	parent_id := whosonfirst.ParentId(f)
	name := whosonfirst.Name(f)
	placetype := whosonfirst.Placetype(f)
	country := whosonfirst.Country(f)
	repo := whosonfirst.Repo(f)

	path, err := uri.Id2RelPath(id)

	if err != nil {
		return nil, err
	}

	uri, err := uri.Id2AbsPath("https://whosonfirst.mapzen.com/data", id)

	if err != nil {
		return nil, err
	}

	is_current := 0
	is_ceased := 0
	is_deprecated := 0
	is_superseded := 0
	is_superseding := 0

	current, current_known := whosonfirst.IsCurrent(f)

	if current == true {
		is_current = 1
	}

	if current == false && current_known == false {
		is_current = -1
	}

	if whosonfirst.IsCeased(f) {
		is_ceased = 1
	}

	if whosonfirst.IsDeprecated(f) {
		is_deprecated = 1
	}

	if whosonfirst.IsSuperseded(f) {
		is_superseded = 1
	}

	if whosonfirst.IsSuperseding(f) {
		is_superseding = 1
	}

	superseded_by := whosonfirst.SupersededBy(f)
	supersedes := whosonfirst.Supersedes(f)

	spr := WOFStandardPlacesResult{
		WOFId:           id,
		WOFParentId:     parent_id,
		WOFPlacetype:    placetype,
		WOFName:         name,
		WOFCountry:      country,
		WOFRepo:         repo,
		WOFPath:         path,
		MZURI:           uri,
		MZIsCurrent:     is_current,
		MZIsCeased:      is_ceased,
		MZIsDeprecated:  is_deprecated,
		MZIsSuperseded:  is_superseded,
		MZIsSuperseding: is_superseding,
		WOFSupersedes:   supersedes,
		WOFSupersededBy: superseded_by,
	}

	return &spr, nil
}

func (spr *WOFStandardPlacesResult) Id() int64 {
	return spr.WOFId
}

func (spr *WOFStandardPlacesResult) ParentId() int64 {
	return spr.WOFParentId
}

func (spr *WOFStandardPlacesResult) Name() string {
	return spr.WOFName
}

func (spr *WOFStandardPlacesResult) Placetype() string {
	return spr.WOFPlacetype
}

func (spr *WOFStandardPlacesResult) Country() string {
	return spr.WOFCountry
}

func (spr *WOFStandardPlacesResult) Repo() string {
	return spr.WOFRepo
}

func (spr *WOFStandardPlacesResult) Path() string {
	return spr.WOFPath
}

func (spr *WOFStandardPlacesResult) URI() string {
	return spr.MZURI
}

func (spr *WOFStandardPlacesResult) IsCurrent() spr.ExistentialFlag {

     return NewWOFExistentialFlag(spr.MZIsCurrent)
}

func (spr *WOFStandardPlacesResult) IsCeased() bool {

     return NewWOFExistentialFlag(spr.MZIsCeased)
}

func (spr *WOFStandardPlacesResult) IsDeprecated() bool {

     return NewWOFExistentialFlag(spr.MZIsDeprecated)
}

func (spr *WOFStandardPlacesResult) IsSuperseded() bool {

     return NewWOFExistentialFlag(spr.MZIsSuperseded)
}

func (spr *WOFStandardPlacesResult) IsSuperseding() bool {

     return NewWOFExistentialFlag(spr.MZIsSuperseding)
}

func (spr *WOFStandardPlacesResult) SupersededBy() []int64 {
	return spr.WOFSupersededBy
}

func (spr *WOFStandardPlacesResult) Supersedes() []int64 {
	return spr.WOFSupersedes
}

func NewWOFExistentialFlag(v int64) (*WOFExistentialFlag) {

     	var status_int int64
	var status bool
	var confidence bool

	switch v {
	       case 0:
	       	    status_int = v
		    status = false
		    confidence = true
	       case 1:
	       	    status_int = v
		    status = true
		    confidence = true
	       default:
	       	    status_int = v
		    status = false
		    confidence = false
        }

	flag := WOFExistentialFlag{
	     WOFStatusInt: status_int,
	     WOFStatus: status
	     WOFConfidence: confidence,
	}

	return &flag
}