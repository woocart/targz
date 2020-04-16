package main

import (
	"os"
	"time"

	"github.com/woocart/targz/logger"
	"github.com/woocart/targz/version"
	"go.uber.org/zap"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var log = logger.New("targz")
var app = kingpin.New("targz", "Compress data to Tar using gzip compression.")

var from = app.Arg("from", "Directory or file to compress").Required().ExistingFileOrDir()
var to = app.Arg("to", "Where to write to [path or stdout]").Required().String()
var exclude = app.Flag("exclude", "Exclude files matching PATTERN, a glob(3)-style wildcard pattern.").Strings()
var afterString = app.Flag("afterDate", "Exclude files created before date (RFC3339).").String()

func main() {
	app.Author("dz0ny")
	app.Version(version.String())
	kingpin.MustParse(app.Parse(os.Args[1:]))
	after := time.Unix(0, 0) // from start of counting
	if *afterString != "" {
		if *afterString == "now" {
			after = time.Now()
		} else {
			afterParsed, err := time.Parse(time.RFC3339, *afterString)
			if err != nil {
				log.Fatal("Could not parse afterDate", err)
			}
			after = afterParsed
		}

		log.Infow("Only including files create or modified", zap.String("afterDate", after.Format(time.RFC3339)))
	}

	if *to == "stdout" {
		err := Tar(*from, os.Stdout, *exclude, after)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		f, err := os.Create(*to)
		if err != nil {
			log.Fatal(err)
		}
		err = Tar(*from, f, *exclude, after)
		if err != nil {
			log.Fatal(err)
		}
	}

}
