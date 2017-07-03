package main

import (
	"fmt"
	"io/ioutil"
	"testing"
)

// TestParseExiftool tests the ParseExiftoolOutput function.
func TestParseExiftool(t *testing.T) {
	b, err := ioutil.ReadFile("test/exiftool.out")
	if err != nil {
		fmt.Print(err)
	}

	results := ParseExiftoolOutput(string(b), nil)

	if err != nil {
		t.Log(err)
	}

	if true {
		t.Log("results: ", results)
	}

}

// TestParseVersion tests the GetFSecureVersion function.
func TestParseTRiD(t *testing.T) {
	b, err := ioutil.ReadFile("test/trid.out") // just pass the file name
	if err != nil {
		fmt.Print(err)
	}

	trid := ParseTRiDOutput(string(b), nil)

	if true {
		t.Log("trid: ", trid)
	}

}
