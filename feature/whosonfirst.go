package feature

import (
	"encoding/json"
	"errors"
	"github.com/skelterjohn/geom"
	"github.com/whosonfirst/go-whosonfirst-flags"
	"github.com/whosonfirst/go-whosonfirst-flags/existential"
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
	MZIsCurrent              int64   `json:"mz:is_current"`
	MZIsCeased               int64   `json:"mz:is_ceased"`
	MZIsDeprecated           int64   `json:"mz:is_deprecated"`
	MZIsSuperseded           int64   `json:"mz:is_superseded"`
	MZIsSuperseding          int64   `json:"mz:is_superseding"`
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

	is_current, err := whosonfirst.IsCurrent(f)

	if err != nil {
		return nil, err
	}

	is_ceased, err := whosonfirst.IsCeased(f)

	if err != nil {
		return nil, err
	}

	is_deprecated, err := whosonfirst.IsDeprecated(f)

	if err != nil {
		return nil, err
	}

	is_superseded, err := whosonfirst.IsSuperseded(f)

	if err != nil {
		return nil, err
	}

	is_superseding, err := whosonfirst.IsSuperseding(f)

	if err != nil {
		return nil, err
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
		MZIsCurrent:     is_current.Flag(),
		MZIsCeased:      is_ceased.Flag(),
		MZIsDeprecated:  is_deprecated.Flag(),
		MZIsSuperseded:  is_superseded.Flag(),
		MZIsSuperseding: is_superseding.Flag(),
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

func (spr *WOFStandardPlacesResult) IsCurrent() flags.ExistentialFlag {
	return existentialFlag(spr.MZIsCurrent)
}

func (spr *WOFStandardPlacesResult) IsCeased() flags.ExistentialFlag {
	return existentialFlag(spr.MZIsCeased)
}

func (spr *WOFStandardPlacesResult) IsDeprecated() flags.ExistentialFlag {
	return existentialFlag(spr.MZIsDeprecated)
}

func (spr *WOFStandardPlacesResult) IsSuperseded() flags.ExistentialFlag {
	return existentialFlag(spr.MZIsSuperseded)
}

func (spr *WOFStandardPlacesResult) IsSuperseding() flags.ExistentialFlag {
	return existentialFlag(spr.MZIsSuperseding)
}

func (spr *WOFStandardPlacesResult) SupersededBy() []int64 {
	return spr.WOFSupersededBy
}

func (spr *WOFStandardPlacesResult) Supersedes() []int64 {
	return spr.WOFSupersedes
}

// we're going to assume that this won't fail since we already go through
// the process of instantiating `flags.ExistentialFlag` thingies in SPR()
// if we need to we'll just cache those instances in the `spr *WOFStandardPlacesResult`
// thingy (and omit them from the JSON output) but today that is unnecessary
// (20170816/thisisaaronland)

func existentialFlag(i int64) flags.ExistentialFlag {
	fl, _ := existential.NewKnownUnknownFlag(i)
	return fl
}
