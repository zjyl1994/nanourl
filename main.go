package main

import (
	"net/url"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/joho/godotenv/autoload"
	log "github.com/sirupsen/logrus"
	"github.com/speps/go-hashids/v2"
	"github.com/zjyl1994/nanourl/model"
	"github.com/zjyl1994/nanourl/server"
	"github.com/zjyl1994/nanourl/vars"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	var err error
	vars.Listen = os.Getenv("NANOURL_LISTEN")
	vars.DataDir = os.Getenv("NANOURL_PATH")
	vars.BaseUrl = os.Getenv("NANOURL_BASEURL")

	if vars.Listen == "" {
		vars.Listen = vars.DEFAULT_LISTEN
	}
	if vars.BaseUrl == "" {
		vars.BaseUrl = vars.DEFAULT_BASEURL
	} else {
		if !strings.HasSuffix(vars.BaseUrl, "/") {
			vars.BaseUrl += "/"
		}
		if _, err = url.Parse(vars.BaseUrl); err != nil {
			log.Fatalln("NANOURL_BASEURL not valid url", err.Error())
		}
	}

	vars.HashId, err = hashids.New()
	if err != nil {
		log.Fatalln(err.Error())
	}

	vars.DB, err = gorm.Open(sqlite.Open(filepath.Join(vars.DataDir, "nanourl.sqlite")))
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = vars.DB.AutoMigrate(&model.URLObject{})
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = server.Run(vars.Listen)
	if err != nil {
		log.Fatalln(err.Error())
	}
}
