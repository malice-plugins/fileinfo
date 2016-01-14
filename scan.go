package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/codegangsta/cli"
	"github.com/crackcomm/go-clitable"
	"github.com/parnurzeal/gorequest"
)

// Version stores the plugin's version
var Version string

// BuildTime stores the plugin's build time
var BuildTime string

// FileInfo json object
type FileInfo struct {
	SSDeep   string            `json:"ssdeep"`
	TRiD     []string          `json:"trid"`
	Exiftool map[string]string `json:"exiftool"`
}

func getopt(name, dfault string) string {
	value := os.Getenv(name)
	if value == "" {
		value = dfault
	}
	return value
}

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
func ParseExiftoolOutput(exifout string) map[string]string {

	var ignoreTags = []string{
		"Directory",
		"File Name",
		"File Permissions",
		"File Modification Date/Time",
	}

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

	return datas
}

// ParseSsdeepOutput convert ssdeep output into JSON
func ParseSsdeepOutput(ssdout string) string {

	// Break output into lines
	lines := strings.Split(ssdout, "\n")
	// Break second line into hash and path
	hashAndPath := strings.Split(lines[1], ",")

	return strings.TrimSpace(hashAndPath[0])
}

// ParseTRiDOutput convert trid output into JSON
func ParseTRiDOutput(tridout string) []string {

	keepLines := []string{}

	lines := strings.Split(tridout, "\n")
	lines = lines[6:]
	// fmt.Println(lines)

	for _, line := range lines {
		if len(strings.TrimSpace(line)) != 0 {
			keepLines = append(keepLines, strings.TrimSpace(line))
		}
	}

	return keepLines
}

func printStatus(resp gorequest.Response, body string, errs []error) {
	fmt.Println(resp.Status)
}

func printMarkDownTable(finfo FileInfo) {

	// print ssdeep
	fmt.Println("#### SSDeep")
	fmt.Println(finfo.SSDeep)
	// print exiftool
	fmt.Println("#### Exiftool")
	table := clitable.New([]string{"Field", "Value"})
	for key, value := range finfo.Exiftool {
		table.AddRow(map[string]interface{}{"Field": key, "Value": value})
	}
	table.Markdown = true
	table.Print()
	// print trid
	fmt.Println("#### TRiD")
	table = clitable.New([]string{"TRiD", ""})
	for _, trd := range finfo.TRiD {
		fmt.Println(" - ", trd)
	}
}

// AppHelpTemplate custom app help template
var AppHelpTemplate = `Usage: {{.Name}} {{if .Flags}}[OPTIONS] {{end}}COMMAND [arg...]

{{.Usage}}

Version: {{.Version}}{{if or .Author .Email}}
BuildTime: {{.BuildTime}}

Author:{{if .Author}}
  {{.Author}}{{if .Email}} - <{{.Email}}>{{end}}{{else}}
  {{.Email}}{{end}}{{end}}
{{if .Flags}}
Options:
  {{range .Flags}}{{.}}
  {{end}}{{end}}
Commands:
  {{range .Commands}}{{.Name}}{{with .ShortName}}, {{.}}{{end}}{{ "\t" }}{{.Usage}}
  {{end}}
Run '{{.Name}} COMMAND --help' for more information on a command.
`

func main() {
	cli.AppHelpTemplate = AppHelpTemplate
	app := cli.NewApp()
	app.Name = "fileinfo"
	app.Author = "blacktop"
	app.Email = "https://github.com/blacktop"
	app.Version = Version
	app.Compiled = time.Now()
	app.Usage = "Malice File Info Plugin - ssdeep/exiftool/TRiD"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "table, t",
			Usage: "output as Markdown table",
		},
	}
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:   "post, p",
			Usage:  "POST results to Malice webhook",
			EnvVar: "MALICE_ENDPOINT",
		},
	}
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:   "proxy, x",
			Usage:  "proxy settings for Malice webhook endpoint",
			EnvVar: "MALICE_PROXY",
		},
	}
	app.Action = func(c *cli.Context) {
		path := c.Args().First()

		if _, err := os.Stat(path); os.IsNotExist(err) {
			assert(err)
		}

		fileInfo := FileInfo{
			SSDeep:   ParseSsdeepOutput(RunCommand("ssdeep", path)),
			TRiD:     ParseTRiDOutput(RunCommand("trid", path)),
			Exiftool: ParseExiftoolOutput(RunCommand("exiftool", path)),
		}

		if c.Bool("table") {
			printMarkDownTable(fileInfo)
		} else {
			fileInfoJSON, err := json.Marshal(fileInfo)
			assert(err)
			if c.Bool("post") {
				request := gorequest.New()
				if c.Bool("proxy") {
					request = gorequest.New().Proxy(os.Getenv("MALICE_PROXY"))
				}
				request.Post(os.Getenv("MALICE_ENDPOINT")).
					Set("Task", path).
					Send(fileInfoJSON).
					End(printStatus)
			}
			fmt.Println(string(fileInfoJSON))
		}
	}

	app.Run(os.Args)

	// if len(os.Args) < 2 {
	// 	fmt.Println("[ERROR] Missing input file.")
	// 	os.Exit(2)
	// }
	// if len(os.Args) == 2 && os.Args[1] == "--version" {
	// 	fmt.Println("Version: ", Version)
	// 	fmt.Println("BuildTime: ", BuildTime)
	// 	os.Exit(0)
	// }

	// path := os.Args[1]
	//
	// if _, err := os.Stat(path); os.IsNotExist(err) {
	// 	assert(err)
	// }
	//
	// fileInfo := FileInfo{
	// 	SSDeep:   ParseSsdeepOutput(RunCommand("ssdeep", path)),
	// 	TRiD:     ParseTRiDOutput(RunCommand("trid", path)),
	// 	Exiftool: ParseExiftoolOutput(RunCommand("exiftool", path)),
	// }
	//
	// fileInfoJSON, err := json.Marshal(fileInfo)
	// assert(err)
	//
	// fmt.Println(string(fileInfoJSON))
	// printMarkDownTable(fileInfo)
}