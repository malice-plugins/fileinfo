package main

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/fatih/structs"
	"github.com/gorilla/mux"
	"github.com/malice-plugins/pkgs/database"
	"github.com/malice-plugins/pkgs/database/elasticsearch"
	"github.com/malice-plugins/pkgs/utils"
	"github.com/parnurzeal/gorequest"
	"github.com/pkg/errors"
	"github.com/rakyll/magicmime"
	"github.com/urfave/cli"
)

const (
	name     = "fileinfo"
	category = "metadata"
)

var (
	// Version stores the plugin's version
	Version string
	// BuildTime stores the plugin's build time
	BuildTime string

	fi FileInfo

	mtx sync.Mutex

	// es is the elasticsearch database object
	es elasticsearch.Database
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
	Magic    FileMagic         `json:"magic" structs:"magic"`
	SSDeep   string            `json:"ssdeep" structs:"ssdeep"`
	TRiD     []string          `json:"trid" structs:"trid"`
	Exiftool map[string]string `json:"exiftool" structs:"exiftool"`
	MarkDown string            `json:"markdown,omitempty" structs:"markdown,omitempty"`
}

// GetFileMimeType returns the mime-type of a file path
func GetFileMimeType(ctx context.Context, path string) error {

	utils.Assert(magicmime.Open(magicmime.MAGIC_MIME_TYPE | magicmime.MAGIC_SYMLINK | magicmime.MAGIC_ERROR))
	defer magicmime.Close()

	mimetype, err := magicmime.TypeByFile(path)
	if err != nil {
		fi.Magic.Mime = err.Error()
		return err
	}

	fi.Magic.Mime = mimetype
	return nil
}

// GetFileDescription returns the textual libmagic type of a file path
func GetFileDescription(ctx context.Context, path string) error {

	utils.Assert(magicmime.Open(magicmime.MAGIC_SYMLINK | magicmime.MAGIC_ERROR))
	defer magicmime.Close()

	magicdesc, err := magicmime.TypeByFile(path)
	if err != nil {
		fi.Magic.Description = err.Error()
		return err
	}

	fi.Magic.Description = magicdesc
	return nil
}

// ParseExiftoolOutput convert exiftool output into JSON
func ParseExiftoolOutput(exifout string, err error) map[string]string {

	if err != nil {
		m := make(map[string]string)
		m["error"] = err.Error()
		return m
	}

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
func ParseSsdeepOutput(ssdout string, err error) string {

	if err != nil {
		return err.Error()
	}

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
func ParseTRiDOutput(tridout string, err error) []string {

	if err != nil {
		return []string{err.Error()}
	}

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

func generateMarkDownTable(fi FileInfo) string {
	var tplOut bytes.Buffer

	t := template.Must(template.New("fileinfo").Parse(tpl))

	err := t.Execute(&tplOut, fi)
	if err != nil {
		log.Println("executing template:", err)
	}

	return tplOut.String()
}

func printStatus(resp gorequest.Response, body string, errs []error) {
	fmt.Println(body)
}

func webService() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/scan", webAvScan).Methods("POST")
	log.Info("web service listening on port :3993")
	log.Fatal(http.ListenAndServe(":3993", router))
}

func getMD5(text string) string {
	hMD5 := md5.New()
	_, _ = hMD5.Write([]byte(text))
	return fmt.Sprintf("%x", hMD5.Sum(nil))
}

func webAvScan(w http.ResponseWriter, r *http.Request) {

	r.ParseMultipartForm(32 << 20)
	file, header, err := r.FormFile("malware")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Please supply a valid file to scan.")
		log.Error(err)
	}
	defer file.Close()

	log.Debug("Uploaded fileName: ", header.Filename)

	tmpfile, err := ioutil.TempFile("/malware", "web_")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	if _, err = tmpfile.Write(data); err != nil {
		log.Fatal(err)
	}
	if err = tmpfile.Close(); err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(60)*time.Second)
	defer cancel()

	// Do FileInfo scan
	mtx.Lock()
	defer mtx.Unlock()
	path := tmpfile.Name()
	GetFileMimeType(ctx, path)
	GetFileDescription(ctx, path)

	fileInfo := FileInfo{
		Magic:    fi.Magic,
		SSDeep:   ParseSsdeepOutput(utils.RunCommand(ctx, "ssdeep", path)),
		TRiD:     ParseTRiDOutput(utils.RunCommand(ctx, "trid", path)),
		Exiftool: ParseExiftoolOutput(utils.RunCommand(ctx, "exiftool", path)),
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(fileInfo); err != nil {
		log.Fatal(err)
	}
}

func main() {

	cli.AppHelpTemplate = utils.AppHelpTemplate
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
			Name:   "callback, c",
			Usage:  "POST results to Malice webhook",
			EnvVar: "MALICE_ENDPOINT",
		},
		cli.BoolFlag{
			Name:   "proxy, x",
			Usage:  "proxy settings for Malice webhook endpoint",
			EnvVar: "MALICE_PROXY",
		},
		cli.StringFlag{
			Name:        "elasticsearch",
			Value:       "",
			Usage:       "elasticsearch url for Malice to store results",
			EnvVar:      "MALICE_ELASTICSEARCH_URL",
			Destination: &es.URL,
		},
		cli.IntFlag{
			Name:   "timeout",
			Value:  10,
			Usage:  "malice plugin timeout (in seconds)",
			EnvVar: "MALICE_TIMEOUT",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:  "web",
			Usage: "Create a File Info web service",
			Action: func(c *cli.Context) error {
				webService()
				return nil
			},
		},
	}
	app.Action = func(c *cli.Context) error {
		var err error

		if c.Bool("verbose") {
			log.SetLevel(log.DebugLevel)
		}

		if c.Args().Present() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.Int("timeout"))*time.Second)
			defer cancel()

			path := c.Args().First()

			if _, err = os.Stat(path); os.IsNotExist(err) {
				utils.Assert(err)
			}

			if c.Bool("verbose") {
				log.SetLevel(log.DebugLevel)
			}

			if c.Bool("mime") {
				GetFileMimeType(ctx, path)
				fmt.Println(fi.Magic.Mime)
				return nil
			}

			// run libmagic
			err = GetFileMimeType(ctx, path)
			if err != nil && ctx.Err() == nil {
				// try again
				GetFileMimeType(ctx, path)
			}
			err = GetFileDescription(ctx, path)
			if err != nil && ctx.Err() == nil {
				// try again
				GetFileDescription(ctx, path)
			}

			fileInfo := FileInfo{
				Magic:    fi.Magic,
				SSDeep:   ParseSsdeepOutput(utils.RunCommand(ctx, "ssdeep", path)),
				TRiD:     ParseTRiDOutput(utils.RunCommand(ctx, "trid", path)),
				Exiftool: ParseExiftoolOutput(utils.RunCommand(ctx, "exiftool", path)),
			}
			fileInfo.MarkDown = generateMarkDownTable(fileInfo)

			// upsert into Database
			if len(c.String("elasticsearch")) > 0 {
				err := es.Init()
				if err != nil {
					return errors.Wrap(err, "failed to initalize elasticsearch")
				}
				err = es.StorePluginResults(database.PluginResults{
					ID:       utils.Getopt("MALICE_SCANID", utils.GetSHA256(path)),
					Name:     name,
					Category: category,
					Data:     structs.Map(fileInfo),
				})
				if err != nil {
					return errors.Wrapf(err, "failed to index malice/%s results", name)
				}
			}

			if c.Bool("table") {
				fmt.Println(fileInfo.MarkDown)
			} else {
				fileInfo.MarkDown = ""
				fileInfoJSON, err := json.Marshal(fileInfo)
				utils.Assert(err)
				if c.Bool("post") {
					request := gorequest.New()
					if c.Bool("proxy") {
						request = gorequest.New().Proxy(os.Getenv("MALICE_PROXY"))
					}
					request.Post(os.Getenv("MALICE_ENDPOINT")).
						Set("X-Malice-ID", utils.Getopt("MALICE_SCANID", utils.GetSHA256(path))).
						Send(string(fileInfoJSON)).
						End(printStatus)

					return nil
				}
				// write to stdout
				fmt.Println(string(fileInfoJSON))
			}
		} else {
			log.Fatal(fmt.Errorf("Please supply a file to scan with malice/%s", name))
		}
		return nil
	}

	err := app.Run(os.Args)
	utils.Assert(err)
}
