package eephttpd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"

    "github.com/d5/tengo/script"
)

func (f *EepHttpd) ServeHTTP(rw http.ResponseWriter, rq *http.Request) {
	if strings.HasSuffix(rq.URL.Path, ".md") {
		f.HandleMarkdown(rw, rq)
		return
	}

	/* Eventually I'll give this some simple scriptable capabilities.
	    if strings.HasSuffix(rq.URL.Path, ".tengo") {
			r.HandleScript(rw, rq)
		    return
		}
	*/
	f.HandleFile(rw, rq)
}

func (f *EepHttpd) HandleMarkdown(rw http.ResponseWriter, rq *http.Request) {
	path := filepath.Join(f.ServeDir, rq.URL.Path)
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	f.mark.Render(rw, bytes)
}

func (f *EepHttpd) HandleFile(rw http.ResponseWriter, rq *http.Request) {
	path := filepath.Join(f.ServeDir, rq.URL.Path)
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		f.HandleMissing(rw, rq)
	}
	fmt.Fprintf(rw, string(bytes))
}

func (f *EepHttpd) HandleMissing(rw http.ResponseWriter, rq *http.Request) {
	path := filepath.Join(f.ServeDir, rq.URL.Path)
	fmt.Fprintf(rw, "ERROR %s NOT FOUND", path)
}
