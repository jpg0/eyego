package eyego

import (
	"testing"
	"bytes"
	"fmt"
)

func TestSingleAP(t *testing.T){

	log, err := ParseLog(bytes.NewBufferString("8,8,POWERON\n" +
	"12,12,NEWAP,0023691de7c2,17,0042a5c6\n" +
	"35,1358360774,NEWPHOTO,P1040355.JPG,4875264"))

	if err != nil {
		t.Error(err)
	}

	aps := log.AccessPoints("P1040355.JPG")

	fmt.Println(log)

	if len(aps) != 1 {
		t.Fatalf("Failed to find access point")
	}

	if aps[0].MacAddress != "00:23:69:1d:e7:c2" {
		t.Errorf("Bad mac address, expected 00:23:69:1d:e7:c2, was %s", aps[0].MacAddress)
	}
}

func TestNoAP(t *testing.T){

	log, err := ParseLog(bytes.NewBufferString("8,8,POWERON\n" +
			"35,1358360774,NEWPHOTO,P1040355.JPG,4875264"))

	if err != nil {
		t.Error(err)
	}

	aps := log.AccessPoints("P1040355.JPG")

	fmt.Println(log)

	if len(aps) != 0 {
		t.Fatalf("Found invalid access point")
	}
}

func TestSurroundingAP(t *testing.T){

	log, err := ParseLog(bytes.NewBufferString("8,8,POWERON\n" +
			"12,12,NEWAP,0023691de7c2,17,0042a5c6\n" +
			"35,1358360774,NEWPHOTO,P1040355.JPG,4875264\n" +
			"47,47,NEWAP,0023691de7c3,17,0042a5c6\n"	))

	if err != nil {
		t.Error(err)
	}

	aps := log.AccessPoints("P1040355.JPG")

	fmt.Println(log)

	if len(aps) != 2 {
		t.Fatalf("Expected to find %d access points, found %d", 2, len(aps))
	}

	var aps0, aps1 AccessPointSightingInfo

	if aps[0].MacAddress == "00:23:69:1d:e7:c2" {
		aps0 = aps[0]
		aps1 = aps[1]
	} else {
		aps0 = aps[1]
		aps1 = aps[0]
	}

	if aps0.MacAddress != "00:23:69:1d:e7:c2" {
		t.Errorf("Bad mac address, expected 00:23:69:1d:e7:c2, was %s", aps0.MacAddress)
	}

	if aps1.MacAddress != "00:23:69:1d:e7:c3" {
		t.Errorf("Bad mac address, expected 00:23:69:1d:e7:c3, was %s", aps1.MacAddress)
	}


}
