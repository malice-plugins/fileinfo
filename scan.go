package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func assert(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// RunCommand runs cmd on file
func RunCommand(cmd string, path string) string {

	cmdOut, err := exec.Command(cmd, path).Output()
	assert(err)

	return string(cmdOut)
}

// ParseExiftoolOutput convert exiftool output into JSON
func ParseExiftoolOutput(exifout string) []byte {

	var ignoreTags = []string{
		"Directory",
		"File Name",
		"File Permissions",
		"File Modification Date/Time",
	}

	exifJSON := make(map[string]map[string]string)
	lines := strings.Split(exifout, "\n")
	datas := make(map[string]string, len(lines))

	for _, line := range lines {
		keyvalue := strings.Split(line, ":")
		if len(keyvalue) != 2 {
			continue
		}
		if !stringInSlice(strings.TrimSpace(keyvalue[0]), ignoreTags) {
			datas[strings.TrimSpace(keyvalue[0])] = strings.TrimSpace(keyvalue[1])
		}
	}

	exifJSON["exiftool"] = datas
	jsonExif, err := json.Marshal(exifJSON)
	assert(err)

	return jsonExif
}

// ParseSsdeepOutput convert ssdeep output into JSON
func ParseSsdeepOutput(ssdout string) []byte {

	datas := make(map[string]string, 1)
	// Break output into lines
	lines := strings.Split(ssdout, "\n")
	// Break second line into hash and path
	hashAndPath := strings.Split(lines[1], ",")
	// Add hash to map
	datas["ssdeep"] = strings.TrimSpace(hashAndPath[0])

	jsonSsdeep, err := json.Marshal(datas)
	assert(err)

	return jsonSsdeep
}

// ParseTRiDOutput convert trid output into JSON
func ParseTRiDOutput(tridout string) []byte {

	keepLines := []string{}
	datas := make(map[string][]string, 1)

	lines := strings.Split(tridout, "\n")
	lines = lines[6:]
	// fmt.Println(lines)

	for _, line := range lines {
		if len(line) != 0 {
			keepLines = append(keepLines, strings.TrimSpace(line))
		}
	}

	datas["trid"] = keepLines
	j, err := json.Marshal(datas)
	assert(err)

	return j
}

func main() {
	path := os.Args[1]

	if _, err := os.Stat(path); os.IsNotExist(err) {
		assert(err)
	}
	ssdeepJSON := ParseSsdeepOutput(RunCommand("ssdeep", path))
	tridJSON := ParseTRiDOutput(RunCommand("trid", path))
	exiftoolJSON := ParseExiftoolOutput(RunCommand("exiftool", path))
	fmt.Println(string(ssdeepJSON))
	fmt.Println(string(tridJSON))
	fmt.Println(string(exiftoolJSON))
}
