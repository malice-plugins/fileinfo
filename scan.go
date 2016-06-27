package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/crackcomm/go-clitable"
	"github.com/parnurzeal/gorequest"
	"github.com/urfave/cli"
	r "gopkg.in/dancannon/gorethink.v2"
)

// Version stores the plugin's version
var Version string

// BuildTime stores the plugin's build time
var BuildTime string

const (
	name     = "fileinfo"
	category = "metadata"
)

type pluginResults struct {
	ID       string   `gorethink:"id"`
	FileInfo FileInfo `gorethink:"fileinfo"`
}

// FileInfo json object
type FileInfo struct {
	SSDeep   string            `json:"ssdeep" gorethink:"ssdeep"`
	TRiD     []string          `json:"trid" gorethink:"trid"`
	Exiftool map[string]string `json:"exiftool" gorethink:"exiftool"`
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

// getSHA256 calculates a file's sha256sum
func getSHA256(name string) string {

	dat, err := ioutil.ReadFile(name)
	assert(err)

	h256 := sha256.New()
	_, err = h256.Write(dat)
	assert(err)

	return fmt.Sprintf("%x", h256.Sum(nil))
}

// RunCommand runs cmd on file
func RunCommand(cmd string, path string) string {

	cmdOut, err := exec.Command(cmd, path).Output()
	if len(cmdOut) == 0 {
		assert(err)
	}

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
	fmt.Println()
	// print trid
	fmt.Println("#### TRiD")
	table := clitable.New([]string{"TRiD", ""})
	for _, trd := range finfo.TRiD {
		fmt.Println(" - ", trd)
	}
	fmt.Println()
	// print exiftool
	fmt.Println("#### Exiftool")
	table = clitable.New([]string{"Field", "Value"})
	for key, value := range finfo.Exiftool {
		table.AddRow(map[string]interface{}{"Field": key, "Value": value})
	}
	table.Markdown = true
	table.Print()
}

// writeToDatabase upserts plugin results into Database
func writeToDatabase(results pluginResults) {

	// connect to RethinkDB
	session, err := r.Connect(r.ConnectOpts{
		Address:  fmt.Sprintf("%s:28015", getopt("MALICE_RETHINKDB", "rethink")),
		Timeout:  5 * time.Second,
		Database: "malice",
	})
	defer session.Close()

	if err == nil {
		res, err := r.Table("samples").Get(results.ID).Run(session)
		assert(err)
		defer res.Close()

		if res.IsNil() {
			// upsert into RethinkDB
			resp, err := r.Table("samples").Insert(results, r.InsertOpts{Conflict: "replace"}).RunWrite(session)
			assert(err)
			log.Debug(resp)
		} else {
			resp, err := r.Table("samples").Get(results.ID).Update(map[string]interface{}{
				"plugins": map[string]interface{}{
					category: map[string]interface{}{
						name: results.FileInfo,
					},
				},
			}).RunWrite(session)
			assert(err)

			log.Debug(resp)
		}

	} else {
		log.Debug(err)
	}
}

var appHelpTemplate = `Usage: {{.Name}} {{if .Flags}}[OPTIONS] {{end}}COMMAND [arg...]

{{.Usage}}

Version: {{.Version}}{{if or .Author .Email}}

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
	cli.AppHelpTemplate = appHelpTemplate
	app := cli.NewApp()
	app.Name = "fileinfo"
	app.Author = "blacktop"
	app.Email = "https://github.com/blacktop"
	app.Version = Version + ", BuildTime: " + BuildTime
	app.Compiled, _ = time.Parse("20060102", BuildTime)
	app.Usage = "Malice File Info Plugin - ssdeep/exiftool/TRiD"
	var rethinkdb string
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "verbose, V",
			Usage: "verbose output",
		},
		cli.BoolFlag{
			Name:  "table, t",
			Usage: "output as Markdown table",
		},
		cli.BoolFlag{
			Name:   "post, p",
			Usage:  "POST results to Malice webhook",
			EnvVar: "MALICE_ENDPOINT",
		},
		cli.BoolFlag{
			Name:   "proxy, x",
			Usage:  "proxy settings for Malice webhook endpoint",
			EnvVar: "MALICE_PROXY",
		},
		cli.StringFlag{
			Name:        "rethinkdb",
			Value:       "",
			Usage:       "rethinkdb address for Malice to store results",
			EnvVar:      "MALICE_RETHINKDB",
			Destination: &rethinkdb,
		},
	}
	app.Action = func(c *cli.Context) error {
		path := c.Args().First()

		if _, err := os.Stat(path); os.IsNotExist(err) {
			assert(err)
		}

		if c.Bool("verbose") {
			log.SetLevel(log.DebugLevel)
		} else {
			r.Log.Out = ioutil.Discard
		}

		fileInfo := FileInfo{
			SSDeep:   ParseSsdeepOutput(RunCommand("ssdeep", path)),
			TRiD:     ParseTRiDOutput(RunCommand("trid", path)),
			Exiftool: ParseExiftoolOutput(RunCommand("exiftool", path)),
		}

		// upsert into Database
		writeToDatabase(pluginResults{
			ID:       getopt("MALICE_SCANID", getSHA256(path)),
			FileInfo: fileInfo,
		})

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
			// write to stdout
			fmt.Println(string(fileInfoJSON))
		}
		return nil
	}

	err := app.Run(os.Args)
	assert(err)
}
