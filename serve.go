package eephttpd

import (

    "strings"
	"net/http"
)

func (f *EepHttpd) ServeHTTP(rw http.ResponseWriter, rq *http.Request) {
    if strings.HasSuffix(rq.URL.Path, ".md") {
        f.HandleMarkdown(rw, rq)
    }
}

func (f *EepHttpd) HandleMarkdown(rw http.ResponseWriter, rq *http.Request) {

}
