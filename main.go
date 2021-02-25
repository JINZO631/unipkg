package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println(`Missing filename ("unipkg [filename]")`)
		return
	}
	filename := os.Args[1]

	file, err := os.Open(filename)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalln(err)
	}

	gzipReader, err := gzip.NewReader(bytes.NewReader(b))
	if err != nil {
		log.Fatalln(err)
	}
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln(err)
		}

		if strings.HasSuffix(header.Name, "/pathname") {
			guid := strings.TrimSuffix(header.Name, "/pathname")
			bytes, err := ioutil.ReadAll(tarReader)
			if err != nil {
				log.Fatalln(err)
			}
			path := string(bytes)
			fmt.Printf("%s %s\n", guid, path)
		}
	}
}
