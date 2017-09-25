package res

import (
	"net/http"
	"path/filepath"
	"log"
	"strings"
	"os"
	"flag"
)

var staticPath = flag.String("sp", "", "Static File Path.")

func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}


func init() {
	if !flag.Parsed() {
		flag.Parse()
	}
	if *staticPath == "" {
		*staticPath = getCurrentDirectory()
		*staticPath += "/static"
	}

	log.Println("Static path:" + *staticPath)
	http.Handle("/", http.FileServer(http.Dir(*staticPath)))
}
