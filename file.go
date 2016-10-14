package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/go-units"
	er "github.com/maliceio/malice/malice/errors"
	"github.com/rakyll/magicmime"
)

// FileData is a file object's metadata
type FileData struct {
	Name string `json:"name" structs:"name"`
	Path string `json:"path" structs:"path"`
	// Valid bool   `json:"valid"`
	Size string `json:"size" structs:"size"`
	// Size   int64
	// CRC32  string
	MD5    string `json:"md5" structs:"md5"`
	SHA1   string `json:"sha1" structs:"sha1"`
	SHA256 string `json:"sha256" structs:"sha256"`
	SHA512 string `json:"sha512" structs:"sha512"`
	// Ssdeep string `json:"ssdeep"`
	Mime  string `json:"mime" structs:"mime"`
	Magic string `json:"magic" structs:"magic"`
	// Arch string `json:"arch"`
}

// Init initializes the File object
func (file *FileData) Init() {

	if file.Path == "" {
		log.Fatalf("error occured during file.Init() because file.Path was not set.")
	}

	file.GetName()
	file.GetSize()
	file.GetFileMimeType()
	file.GetFileMagicType()

	// Read in file data
	dat, err := ioutil.ReadFile(file.Path)
	er.CheckError(err)

	file.GetMD5(dat)
	file.GetSHA1(dat)
	file.GetSHA256(dat)
	file.GetSHA512(dat)
}

// GetName returns file name
func (file *FileData) GetName() (name string, err error) {

	fileHandle, err := os.Open(file.Path)
	if err != nil {
		return
	}
	defer fileHandle.Close()

	stat, err := fileHandle.Stat()
	if err != nil {
		return
	}

	name = stat.Name()

	file.Name = name

	return
}

// GetSize calculates file Size
func (file *FileData) GetSize() (bytes int64, err error) {

	fileHandle, err := os.Open(file.Path)
	if err != nil {
		return
	}
	defer fileHandle.Close()

	stat, err := fileHandle.Stat()
	if err != nil {
		return
	}

	bytes = stat.Size()

	file.Size = units.HumanSize(float64(bytes))

	return
}

// GetMD5 calculates the Files md5sum
func (file *FileData) GetMD5(data []byte) (hMd5Sum string, err error) {

	hmd5 := md5.New()
	_, err = hmd5.Write(data)
	er.CheckError(err)
	hMd5Sum = fmt.Sprintf("%x", hmd5.Sum(nil))

	file.MD5 = hMd5Sum

	return
}

// GetSHA1 calculates the Files sha256sum
func (file *FileData) GetSHA1(data []byte) (h1Sum string, err error) {

	h1 := sha1.New()
	_, err = h1.Write(data)
	er.CheckError(err)
	h1Sum = fmt.Sprintf("%x", h1.Sum(nil))

	file.SHA1 = h1Sum

	return
}

// GetSHA256 calculates the Files sha256sum
func (file *FileData) GetSHA256(data []byte) (h256Sum string, err error) {

	h256 := sha256.New()
	_, err = h256.Write(data)
	er.CheckError(err)
	h256Sum = fmt.Sprintf("%x", h256.Sum(nil))

	file.SHA256 = h256Sum

	return
}

// GetSHA512 calculates the Files sha256sum
func (file *FileData) GetSHA512(data []byte) (h512Sum string, err error) {

	h512 := sha512.New()
	_, err = h512.Write(data)
	er.CheckError(err)
	h512Sum = fmt.Sprintf("%x", h512.Sum(nil))

	file.SHA512 = h512Sum

	return
}

// GetFileMimeType returns the mime-type of a file path
func (file *FileData) GetFileMimeType() (mimetype string, err error) {

	if err = magicmime.Open(magicmime.MAGIC_MIME_TYPE | magicmime.MAGIC_SYMLINK | magicmime.MAGIC_ERROR); err != nil {
		log.Fatal(err)
	}
	defer magicmime.Close()

	mimetype, err = magicmime.TypeByFile(file.Path)
	if err != nil {
		log.Fatalf("error occured during type lookup: %v", err)
	}
	file.Mime = mimetype

	return
}

// GetFileMagicType returns the textual libmagic type of a file path
func (file *FileData) GetFileMagicType() (magictype string, err error) {

	if err = magicmime.Open(magicmime.MAGIC_SYMLINK | magicmime.MAGIC_ERROR); err != nil {
		log.Fatal(err)
	}
	defer magicmime.Close()

	magictype, err = magicmime.TypeByFile(file.Path)
	if err != nil {
		log.Fatalf("error occured during type lookup: %v", err)
	}
	file.Magic = magictype

	return
}

// ToJSON converts File object to []byte JSON
func (file *FileData) ToJSON() []byte {
	fileJSON, err := json.Marshal(file)
	er.CheckError(err)
	return fileJSON
}
