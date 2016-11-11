package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/crackcomm/go-clitable"
	"github.com/fatih/structs"
	"github.com/maliceio/go-plugin-utils/database/elasticsearch"
	"github.com/maliceio/go-plugin-utils/utils"
	"github.com/parnurzeal/gorequest"
	"github.com/rakyll/magicmime"
	"github.com/urfave/cli"
)

var (
	// Version stores the plugin's version
	Version string
	// BuildTime stores the plugin's build time
	BuildTime string

	fi FileInfo

	wg sync.WaitGroup
)

const (
	name     = "fileinfo"
	category = "metadata"
)

type pluginResults struct {
	ID       string   `structs:"id"`
	FileInfo FileInfo `structs:"fileinfo"`
}

// FileMagic is file magic
type FileMagic struct {
	Mime        string `json:"mime" structs:"mime"`
	Description string `json:"description" structs:"description"`
}

// FileInfo json object
type FileInfo struct {
	Magic FileMagic `json:"magic" structs:"magic"`
	// Ssdeep string `json:"ssdeep"`
	SSDeep   string            `json:"ssdeep" structs:"ssdeep"`
	TRiD     []string          `json:"trid" structs:"trid"`
	Exiftool map[string]string `json:"exiftool" structs:"exiftool"`
}

// GetFileMimeType returns the mime-type of a file path
func GetFileMimeType(ctx context.Context, path string) error {
	defer wg.Done()

	c := make(chan struct {
		mimetype string
		err      error
	}, 1)

	go func() {
		utils.Assert(magicmime.Open(magicmime.MAGIC_MIME_TYPE | magicmime.MAGIC_SYMLINK | magicmime.MAGIC_ERROR))
		defer magicmime.Close()

		mt, err := magicmime.TypeByFile(path)
		pack := struct {
			mimetype string
			err      error
		}{mt, err}
		c <- pack
	}()

	select {
	case <-ctx.Done():
		<-c // Wait for mime
		fmt.Println("Cancel the context")
		return ctx.Err()
	case ok := <-c:
		utils.Assert(ok.err)
		fi.Magic.Mime = ok.mimetype
		return ok.err
	}
}

// GetFileDescription returns the textual libmagic type of a file path
func GetFileDescription(ctx context.Context, path string) error {
	defer wg.Done()

	c := make(chan struct {
		magicdesc string
		err       error
	}, 1)

	go func() {
		utils.Assert(magicmime.Open(magicmime.MAGIC_SYMLINK | magicmime.MAGIC_ERROR))
		defer magicmime.Close()

		magicdesc, err := magicmime.TypeByFile(path)
		pack := struct {
			magicdesc string
			err       error
		}{magicdesc, err}
		c <- pack
	}()

	select {
	case <-ctx.Done():
		<-c // Wait for mime
		fmt.Println("Cancel the context")
		return ctx.Err()
	case ok := <-c:
		utils.Assert(ok.err)
		fi.Magic.Description = ok.magicdesc
		return ok.err
	}
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

	log.Debugln("Exiftool lines: ", lines)

	if utils.SliceContainsString("File not found", lines) {
		return nil
	}

	datas := make(map[string]string, len(lines))

	for _, line := range lines {
		keyvalue := strings.Split(line, ":")
		if len(keyvalue) != 2 {
			continue
		}
		if !utils.StringInSlice(strings.TrimSpace(keyvalue[0]), ignoreTags) {
			datas[strings.TrimSpace(utils.CamelCase(keyvalue[0]))] = strings.TrimSpace(keyvalue[1])
		}
	}

	return datas
}

// ParseSsdeepOutput convert ssdeep output into JSON
func ParseSsdeepOutput(ssdout string) string {

	// Break output into lines
	lines := strings.Split(ssdout, "\n")

	log.Debugln("ssdeep lines: ", lines)

	if utils.SliceContainsString("No such file or directory", lines) {
		return ""
	}

	// Break second line into hash and path
	hashAndPath := strings.Split(lines[1], ",")

	return strings.TrimSpace(hashAndPath[0])
}

// ParseTRiDOutput convert trid output into JSON
func ParseTRiDOutput(tridout string) []string {

	keepLines := []string{}

	lines := strings.Split(tridout, "\n")

	log.Debugln("TRiD lines: ", lines)

	if utils.SliceContainsString("Error: found no file(s) to analyze!", lines) {
		return nil
	}

	lines = lines[6:]

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

	fmt.Println("#### Magic")
	table := clitable.New([]string{"Field", "Value"})
	table.AddRow(map[string]interface{}{"Field": "Mime", "Value": finfo.Magic.Mime})
	table.AddRow(map[string]interface{}{"Field": "Description", "Value": finfo.Magic.Description})
	table.Markdown = true
	table.Print()
	fmt.Println()

	if len(finfo.SSDeep) > 0 {
		// print ssdeep
		fmt.Println("#### SSDeep")
		fmt.Println(finfo.SSDeep)
		fmt.Println()
	}

	if finfo.TRiD != nil {
		// print trid
		fmt.Println("#### TRiD")
		for _, trd := range finfo.TRiD {
			fmt.Println(" - ", trd)
		}
		fmt.Println()
	}

	if finfo.Exiftool != nil {
		// print exiftool
		fmt.Println("#### Exiftool")
		table := clitable.New([]string{"Field", "Value"})
		for key, value := range finfo.Exiftool {
			table.AddRow(map[string]interface{}{"Field": key, "Value": value})
		}
		table.Markdown = true
		table.Print()
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
			Name:  "mime, m",
			Usage: "output only mimetype",
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
		cli.IntFlag{
			Name:   "timeout",
			Value:  3,
			Usage:  "malice plugin timeout (in seconds)",
			EnvVar: "MALICE_TIMEOUT",
		},
		cli.StringFlag{
			Name:   "elasitcsearch",
			Value:  "",
			Usage:  "elasitcsearch address for Malice to store results",
			EnvVar: "MALICE_ELASTICSEARCH",
		},
	}
	app.Action = func(c *cli.Context) error {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.Int("timeout"))*time.Second)
		defer cancel()

		path := c.Args().First()

		if _, err := os.Stat(path); os.IsNotExist(err) {
			utils.Assert(err)
		}

		if c.Bool("verbose") {
			log.SetLevel(log.DebugLevel)
		}

		wg.Add(2)
		go GetFileMimeType(ctx, path)
		go GetFileDescription(ctx, path)
		wg.Wait()

		if c.Bool("mime") {
			fmt.Println(fi.Magic.Mime)
			return nil
		}

		var output string
		var err error

		// run ssdeep
		output, err = utils.RunCommand("ssdeep", c.Int("timeout"), path)
		utils.Assert(err)
		ParseSsdeepOutput(output)
		// run trid
		output, err = utils.RunCommand("trid", c.Int("timeout"), path)
		utils.Assert(err)
		ParseTRiDOutput(output)
		// run exiftool
		output, err = utils.RunCommand("exiftool", c.Int("timeout"), path)
		utils.Assert(err)
		ParseExiftoolOutput(output)

		fileInfo := FileInfo{
			Magic:    fi.Magic,
			SSDeep:   fi.SSDeep,
			TRiD:     fi.TRiD,
			Exiftool: fi.Exiftool,
		}

		// upsert into Database
		elasticsearch.InitElasticSearch(c.String("elasitcsearch"))
		elasticsearch.WritePluginResultsToDatabase(elasticsearch.PluginResults{
			ID:       utils.Getopt("MALICE_SCANID", utils.GetSHA256(path)),
			Name:     name,
			Category: category,
			Data:     structs.Map(fileInfo),
		})

		if c.Bool("table") {
			printMarkDownTable(fileInfo)
		} else {
			fileInfoJSON, err := json.Marshal(fileInfo)
			utils.Assert(err)
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
	utils.Assert(err)
}
