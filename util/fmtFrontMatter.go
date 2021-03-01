package util

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/kjk/notionapi"
)

const descLen = 120

type frontMatter struct {
	PageID      string
	Title       string
	Time        int64
}

func (fm *frontMatter) Check(block *notionapi.Block) {
	switch block.Type {
	case notionapi.BlockPage:
		// frontMatter.Date
		fm.Time = block.CreatedTime
		// frontMatter.Title
		fm.Title = block.Title
	}
}

func (fm frontMatter) Print() {
	log.Println(fm.PageID)
	log.Println(fm.Title)
	log.Println(fm.Date())
}

func (fm frontMatter) Date() string {
	return time.Unix(fm.Time/1000, 0).Format("2006-01-02")
}

func (fm frontMatter) GetFile() string {
	var ret string
    log.Println(fm.PageID)
	err := filepath.Walk(contentPath, func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, fm.PageID) {
			ret = path
		}
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return ret
}

func (fm frontMatter) AddFrontMatter(in string) string {
	md := in

	reg := regexp.MustCompile(`{{\ title\ }}`)
	md = reg.ReplaceAllString(md, "\""+fm.Title+"\"")

	reg = regexp.MustCompile(`{{\ date\ }}`)
	md = reg.ReplaceAllString(md, "\""+fm.Date()+"\"")

	return md
}

func fmtFrontMatter(file string, fm frontMatter) {
    log.Println(file)
	mdBytes, err := ioutil.ReadFile(file)
    log.Println(mdBytes)
	if err != nil {
		panic(err)
	}
    tempMD := bytes.NewBuffer(mdBytes)
	// add frontmatter
	doneMD := fm.AddFrontMatter(tempMD.String())
	ioutil.WriteFile(file, []byte(doneMD), os.ModePerm)

	// mv post to folder
	createPostDir()
	extractMD(file)
}

func extractMD(pathwithfile string) {
	_, file := filepath.Split(pathwithfile)
	err := os.Rename(
		filepath.Join(contentPath, file),
		filepath.Join(contentPath, postPath, file),
	)
	if err != nil {
		panic(err)
	}
	log.Println(filepath.Join(contentPath, postPath, file))
}
