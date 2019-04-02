package main

import (
	"os"

	"github.com/woocart/targz/logger"
	"github.com/woocart/targz/version"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var log = logger.New("targz")
var app = kingpin.New("targz", "Compress data to Tar using gzip compression.")

var from = app.Arg("from", "Directory or file to compress").Required().ExistingFileOrDir()
var to = app.Arg("to", "Where to write to [path or stdout]").Required().String()
var exclude = app.Flag("exclude", "Exclude files matching PATTERN, a glob(3)-style wildcard pattern.").Strings()

func main() {
	app.Author("dz0ny")
	app.Version(version.String())
	kingpin.MustParse(app.Parse(os.Args[1:]))

	if *to == "stdout" {
		err := Tar(*from, os.Stdout, *exclude)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		f, err := os.Create(*to)
		if err != nil {
			log.Fatal(err)
		}
		err = Tar(*from, f, *exclude)
		if err != nil {
			log.Fatal(err)
		}
	}

}
