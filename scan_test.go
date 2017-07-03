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

// TestParseTRiD tests the ParseTRiDOutput function.
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

// TestParseTRiDSsdeep tests the ParseSsdeepOutput function.
func TestParseTRiDSsdeep(t *testing.T) {
	b, err := ioutil.ReadFile("test/ssdeep.out") // just pass the file name
	if err != nil {
		fmt.Print(err)
	}

	ssdeep := ParseSsdeepOutput(string(b), nil)

	if true {
		t.Log("ssdeep: ", ssdeep)
	}
}

// TestGenerateMarkDownTable tests the ParseSsdeepOutput function.
func TestGenerateMarkDownTable(t *testing.T) {
	exifOut, err := ioutil.ReadFile("test/exiftool.out")
	if err != nil {
		fmt.Print(err)
	}

	tridOut, err := ioutil.ReadFile("test/trid.out")
	if err != nil {
		fmt.Print(err)
	}

	ssdeepOut, err := ioutil.ReadFile("test/ssdeep.out")
	if err != nil {
		fmt.Print(err)
	}

	fileInfo := FileInfo{
		// Magic:    fi.Magic,
		SSDeep:   ParseSsdeepOutput(string(ssdeepOut), nil),
		TRiD:     ParseTRiDOutput(string(tridOut), nil),
		Exiftool: ParseExiftoolOutput(string(exifOut), nil),
	}
	fileInfo.MarkDown = generateMarkDownTable(fileInfo)

	markDown := generateMarkDownTable(fileInfo)

	if true {
		t.Log("markDown: ", markDown)
	}
}
