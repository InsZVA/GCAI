package res

import (
	"net/http"
	"path/filepath"
	"log"
	"strings"
	"os"
)

func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}


func init() {
	staticDir := getCurrentDirectory()
	staticDir += "/static"
	log.Println("Static path:" + staticDir)
	http.Handle("/", http.FileServer(http.Dir(staticDir)))
}
