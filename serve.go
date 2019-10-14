package eephttpd

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/d5/tengo/script"
)

func (f *EepHttpd) ServeHTTP(rw http.ResponseWriter, rq *http.Request) {
	rp := f.checkURL(rq)
	if strings.HasPrefix(rq.Header.Get("X-User-Agent"), "git") {
		f.HandleGit(rw, rq)
	}
	if strings.HasSuffix(rp, ".md") {
		f.HandleMarkdown(rw, rq)
		return
	}
	if strings.HasSuffix(rp, ".tengo") {
		f.HandleScript(rw, rq)
		return
	}
	f.HandleFile(rw, rq)
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	if info != nil {
		return !info.IsDir()
	}
	return false
}

func (f *EepHttpd) checkURL(rq *http.Request) string {
	p := rq.URL.Path
	if strings.HasSuffix(rq.URL.Path, "/") {
		p = filepath.Join(rq.URL.Path, "index.html")
	}
	if !FileExists(filepath.Join(f.ServeDir, p)) {
		p = filepath.Join(rq.URL.Path, "README.md")
	}
	if !FileExists(filepath.Join(f.ServeDir, p)) {
		p = filepath.Join(rq.URL.Path, "index.tengo")
		if !FileExists(filepath.Join(f.ServeDir, p)) {
			p = rq.URL.Path
		}
	}
	log.Println(p)
	return filepath.Join(f.ServeDir, p)
}

func (f *EepHttpd) HandleScript(rw http.ResponseWriter, rq *http.Request) {
	path := f.checkURL(rq)
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println(err)
		return
	}
	scr := script.New(bytes)
	com, err := scr.Compile()
	if err != nil {
		log.Println(err)
		panic(err)
	}
	if err := com.Run(); err != nil {
		log.Println(err)
		panic(err)
	}
	response := com.Get("response")
	fmt.Fprintf(rw, response.String())
}

func (f *EepHttpd) HandleMarkdown(rw http.ResponseWriter, rq *http.Request) {
	path := f.checkURL(rq)
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	f.mark.Render(rw, bytes)
}

func (f *EepHttpd) HandleGit(rw http.ResponseWriter, rq *http.Request) {
    f.Server.ServeHTTP(rw, rq)
}

func (f *EepHttpd) HandleFile(rw http.ResponseWriter, rq *http.Request) {
	path := f.checkURL(rq)
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		f.HandleMissing(rw, rq)
	}
	fmt.Fprintf(rw, string(bytes))
}

func (f *EepHttpd) HandleMissing(rw http.ResponseWriter, rq *http.Request) {
	path := f.checkURL(rq)
	fmt.Fprintf(rw, "ERROR %s NOT FOUND", strings.Replace(path, f.ServeDir, "", -1))
}
