package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"go.uber.org/zap"
)

// Tar takes a source and variable writers and walks 'source' writing each file
// found to the tar writer; the purpose for accepting multiple writers is to allow
// for multiple outputs (for example a file, or md5 hash)
func Tar(src string, writer io.Writer, excludes []string) error {

	source, err := filepath.Abs(src)
	if err != nil {
		return fmt.Errorf("Unable to expand src - %v", err.Error())
	}

	// ensure the src actually exists before trying to tar it
	if _, err := os.Stat(source); err != nil {
		return fmt.Errorf("Unable to tar files - %v", err.Error())
	}

	gzw, _ := gzip.NewWriterLevel(writer, gzip.BestSpeed)
	defer gzw.Close()

	tw := tar.NewWriter(gzw)
	defer tw.Close()

	// walk path
	return filepath.Walk(source, func(file string, fi os.FileInfo, err error) error {

		// return on any error
		if err != nil {
			log.Warn(err)
			return err
		}

		// Skip folders, we only care about files
		if !fi.Mode().IsRegular() {
			return nil
		}

		filePath := normalize(source, file, strings.HasSuffix(src, string(filepath.Separator)))

		for _, exclude := range excludes {
			if matched, _ := regexp.MatchString(exclude, filePath); matched {
				log.Infow("Skipped", zap.String("file", filePath))
				return nil
			}
		}

		// create a new dir/file header
		fileSize := fi.Size()

		header, err := tar.FileInfoHeader(fi, filePath)
		if err != nil {
			log.Warn(err)
			return err
		}

		// update the name to correctly reflect the desired destination when untaring
		header.Name = filePath
		header.Gid = 0
		header.Uid = 0

		// write the header
		if err := tw.WriteHeader(header); err != nil {
			log.Warn(err)
			return err
		}

		// open files for taring
		f, err := os.Open(file)
		if err != nil {
			log.Warn(err)
			return err
		}

		// copy file data into tar writer
		if _, err := io.Copy(tw, f); err != nil {
			log.Warn(err)
			return err
		}

		if fileSize != fi.Size() {
			log.Fatal("File size changed after write", zap.String("file", filePath))
		}

		log.Infow("Added", zap.String("file", filePath))

		// manually close here after each file operation; defering would cause each file close
		// to wait until all operations have completed.
		f.Close()

		return nil
	})
}

func normalize(base, file string, stripBase bool) string {
	out := strings.Replace(file, base, "", -1)

	// Preserver tar behavior where bases without "filepath.Separator" include the innermost directory
	if !stripBase {
		out = strings.Replace(file, filepath.Dir(base), "", -1)
	}

	return strings.TrimPrefix(out, string(filepath.Separator))
}
