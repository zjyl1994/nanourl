package main

import (
	"net/url"
	"os"
	"path/filepath"
	"strings"

	lru "github.com/hashicorp/golang-lru/v2"
	_ "github.com/joho/godotenv/autoload"
	log "github.com/sirupsen/logrus"
	"github.com/zjyl1994/nanourl/model/db_model"
	"github.com/zjyl1994/nanourl/model/val_obj"
	"github.com/zjyl1994/nanourl/server"
	"github.com/zjyl1994/nanourl/util"
	"github.com/zjyl1994/nanourl/vars"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	var err error
	// init config
	vars.Listen = os.Getenv("NANOURL_LISTEN")
	vars.DataDir = os.Getenv("NANOURL_PATH")
	vars.BaseUrl = os.Getenv("NANOURL_BASEURL")
	vars.RealIpHeader = os.Getenv("NANOURL_REAL_IP_HEADER")

	if vars.Listen == "" {
		vars.Listen = vars.DEFAULT_LISTEN
	}
	if vars.DataDir != "" {
		err = os.MkdirAll(vars.DataDir, 0644)
		if err != nil {
			log.Fatalln(err.Error())
		}
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

	// init vars
	vars.CodeCache, err = lru.New2Q[string, val_obj.URLObject](vars.CODE_CACHE_SIZE)
	if err != nil {
		log.Fatalln(err.Error())
	}

	geoipDBFile := filepath.Join(vars.DataDir, vars.GEOIP_DATABASE_FILENAME)
	if !util.FileExist(geoipDBFile) {
		if err := util.HttpDownload(vars.GEOIP_DOWNLOAD_URL, geoipDBFile, vars.DEFAULT_DOWNLOAD_TIMEOUT); err != nil {
			log.Fatalln(err.Error())
		}
	}
	// init database
	vars.DB, err = gorm.Open(sqlite.Open(filepath.Join(vars.DataDir, "nanourl.sqlite")))
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = vars.DB.AutoMigrate(&db_model.URLObject{}, &db_model.AccessLog{})
	if err != nil {
		log.Fatalln(err.Error())
	}

	// start web server
	err = server.Run(vars.Listen)
	if err != nil {
		log.Fatalln(err.Error())
	}
}
