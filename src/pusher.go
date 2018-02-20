package main

import (
	"os"
	"fmt"
	"path/filepath"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"

	"net/url"
	"net/http"
	"github.com/urfave/cli"
	"strconv"
	"path"
)

type Config struct {
	Generated    string
	Repositories []Repo
}

type Repo struct {
	Name  string
	Url   string
	Cache string
}

func main() {

	app := cli.NewApp()
	app.Name = "push plugin for helm"
	app.Usage = "helm push [reponame] [chartname]"
	app.Version = "0.0.1"
	app.Action = func(c *cli.Context) error {
		fileToPush := locateFile(c.Args().Get(1))
		pushRepo := c.Args().Get(0)
		helmLoc := os.Getenv("HELM_HOME")
		servicePath := os.Getenv("HELM_CHARTS_UPLOAD_ENDPOINT_PATH")
		if servicePath == "" {
			servicePath = "/api/" + pushRepo + "/upload/"
		}

		config := ParseConfig(helmLoc + "/repository/repositories.yaml")
		if len(config.Repositories) == 0 {
			fmt.Println("There are no repositories configured")
			fmt.Println("There are no repositories configured")
			Fail()
		}

		repo := Filter(config.Repositories, func(v Repo) bool {
			return v.Name == pushRepo
		})

		if len(repo) == 0 {
			fmt.Println("There is no matching repository found")
			Fail()
		}
		matching_repo := repo[0]

		uploadUrl := GetQualifiedUrl(matching_repo.Url)

		UploadChart(fileToPush , uploadUrl , servicePath)
		return nil
	}

	app.Run(os.Args)
}

func ParseConfig(location string) Config {

	reposFile, err := filepath.Abs(location)

	if err != nil {
		panic(err)
	}
	fmt.Printf("Loading repos from: %#v\n", reposFile)

	yamlFile, err := ioutil.ReadFile(reposFile)

	if err != nil {
		panic(err)
	}

	var config Config

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		fmt.Println("Failed to deserialize config")
		Fail()
	}
	return config
}

func GetQualifiedUrl(subject string) string{
	url, err := url.Parse(subject)
	if err != nil {
		Fail()
	}
	return string(url.Scheme + "://" + url.Host)
}

func UploadChart(filePath string, remoteUrl string, servicePath string) {
	chartName := path.Base(filePath)

	fmt.Printf("Uploading chart to : %#v\n", remoteUrl)
	data, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Couldn't open file at %s", filePath)
	}
	defer data.Close()

	req, err := http.NewRequest("PUT", remoteUrl + servicePath + chartName, data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "text/plain")

	client := &http.Client{}
	res, err := client.Do(req)
	status, _ := strconv.Atoi(res.Status)

	if ( err != nil || status > 300)  {
		fmt.Printf("Couldn't upload chart")
		fmt.Printf("Status code %d", status)
		Fail()
	}
	fmt.Printf("Successfully pushed chart %s", chartName)

	defer res.Body.Close()
}

func Filter(vs []Repo, f func(Repo) bool) []Repo {
	vsf := make([]Repo, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

func Fail() {
	os.Exit(1)
}

func locateFile(fileName string) string{

	if fileName == "" {
		fmt.Printf("No fileName provided")
		Fail()
	}

	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory")
		os.Exit(1)
	}
	pwd = pwd + "/" + fileName
	if _, err := os.Stat(pwd); err != nil {
		fmt.Printf("No file at %q", pwd)
		Fail()
	}
	return pwd
}
