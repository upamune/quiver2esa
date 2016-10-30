package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/upamune/go-esa/esa"
)

var (
	client      *esa.Client
	accessToken string
	path        string
	team        string
	category    string
)

const defaultCategory = "quiver"

func main() {

	flag.StringVar(&accessToken, "token", "", "Access Token")
	flag.StringVar(&path, "path", "", "Walking Path")
	flag.StringVar(&team, "team", "", "Team Name")
	flag.StringVar(&category, "category", defaultCategory, "Category")
	flag.Parse()

	if accessToken == "" {
		log.Fatal("set a token")
	}

	if path == "" {
		log.Fatal("set a path")
	}

	if team == "" {
		log.Fatal("set a team name")
	}

	client = esa.NewClient(accessToken)
	err := filepath.Walk(path, uploadFile)
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}

func uploadFile(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if info.IsDir() {
		return nil
	}
	if filepath.Ext(path) != ".md" {
		return nil
	}

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	body := string(b)

	fileName := info.Name()
	post := esa.Post{
		Category: category,
		BodyMd:   body,
		Name:     strings.TrimSuffix(fileName, filepath.Ext(path)),
	}
	_, err = client.Post.Create(team, post)
	if err != nil {
		return err
	}

	log.Println("Posted: ", fileName)
	return nil
}
