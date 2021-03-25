package tests

import (
	_ "fmt"
	"github.com/sfomuseum/go-edtf"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/feature"
	"testing"
)

func TestWhosOnFirstSPR(t *testing.T) {

	path := "../fixtures/101851199.geojson"

	f, err := feature.LoadFeatureFromFile(path)

	if err != nil {
		t.Fatalf("Failed to load '%s', %v", path, err)
	}

	spr, err := f.SPR()

	if err != nil {
		t.Fatalf("SAD Failed to create SPR, %v", err)
	}

	if spr.Id() != "101851199" {
		t.Fatalf("Invalid SPR Id() result")
	}

	if spr.Country() != "FR" {
		t.Fatalf("Invalid SPR Country() result")
	}

	inception := spr.Inception()

	if inception == nil {
		t.Fatal("Invalid SPR Inception() result")
	}

	if inception.String() != edtf.UNKNOWN {
		t.Fatal("Invalid SPR Inception() result")
	}

	cessation := spr.Cessation()

	if cessation == nil {
		t.Fatal("Invalid SPR Cessation() result")
	}

	if cessation.String() != "2018-10-24" {
		t.Fatal("Invalid SPR Cessation() result")
	}

	belongsto := spr.BelongsTo()

	if len(belongsto) != 8 {
		t.Fatalf("Invalid SPR BelongsTo() result")
	}

}
