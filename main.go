package main

import (
	"os"
	"path/filepath"

	_ "github.com/joho/godotenv/autoload"
	log "github.com/sirupsen/logrus"
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
