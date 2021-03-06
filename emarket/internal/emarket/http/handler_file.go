package http

import (
	"bytes"
	"emarket/internal/pkg/minify"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func (e *EMarketHandler) setupFileHandler() {
	var csss = []string{
		"/static/css/custom_bootstrap.css",
		"/static/css/app.css",
	}

	var jss = []string{
		"/static/js/app.js",
	}

	allCSS, err := concatFiles(e.webRoot, csss)
	if err != nil {
		log.Fatalln(err)
	}

	allJS, err := concatFiles(e.webRoot, jss)
	if err != nil {
		log.Fatalln(err)
	}

	fileCache := make(map[string][]byte)
	fileCache["/static/css/all.css"] = minify.DoMinify(allCSS, "text/css")
	fileCache["/static/js/all.js"] = minify.DoMinify(allJS, "application/javascript")

	const favicon = "/favicon.ico"
	faviconPath := "/static" + favicon

	e.router.HandleFunc(favicon, func(w http.ResponseWriter, r *http.Request) {
		e.fileHandler(w, r, faviconPath, fileCache)
	})

	e.router.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		e.fileHandler(w, r, r.URL.Path, fileCache)
	})

}

func (e *EMarketHandler) fullpath(file string) string {
	full := e.webRoot + file

	if _, err := os.Stat(full); err == nil {
		return full
	}

	return ""
}

func (e *EMarketHandler) fileHandler(w http.ResponseWriter, r *http.Request, filename string, fileCache map[string][]byte) {
	log := func(err error) {
		fmt.Printf("%v %v", r.URL.Path, err)
	}

	requestedFile, err := filepath.Abs(filename)

	if err != nil {
		log(err)
		e.notFound(w, r)
		return
	}

	content := fileCache[requestedFile]
	ctype, err := detectType(requestedFile)

	if err != nil {
		log(err)
		e.notFound(w, r)
		return
	}

	if content == nil {
		fullPath := e.fullpath(requestedFile)

		if fullPath == "" {
			e.notFound(w, r)
			return
		}

		var err error
		content, err = ioutil.ReadFile(fullPath)

		if err != nil {
			log(err)
			e.notFound(w, r)
			return
		}

		fileCache[requestedFile] = minify.DoMinify(content, ctype)
	}

	setCacheControl(w)
	w.Header().Set("Content-Type", ctype)
	writeResponse(w, r.URL.Path, content)
}

func concatFiles(rootDir string, files []string) ([]byte, error) {
	buf := &bytes.Buffer{}
	for _, file := range files {
		data, err := ioutil.ReadFile(rootDir + "/" + file)
		if err != nil {
			return nil, err
		}
		buf.Write(data)
		buf.Write([]byte("\n"))
	}

	return buf.Bytes(), nil
}
