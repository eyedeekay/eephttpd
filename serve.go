package eephttpd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/d5/tengo/script"
)

func (f *EepHttpd) ProxyRequest(req *http.Request) *http.Request {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil
	}

	// you can reassign the body if you need to parse it as multipart
	//    req.Body = ioutil.NopCloser(bytes.NewReader(body))

	// create a new url from the raw RequestURI sent by the client
	url := fmt.Sprintf("%s://%s:%s/%s", "http", f.SamTracker.Config().TargetHost, f.SamTracker.Config().TargetPort, "a")
	log.Println("handling http tracker request", url)
	proxyReq, err := http.NewRequest(req.Method, url, bytes.NewReader(body))

	// We may want to filter some headers, otherwise we could just use a shallow copy
	//
	proxyReq.Header = req.Header
	return proxyReq
}

func (f *EepHttpd) ServeHTTP(rw http.ResponseWriter, rq *http.Request) {
	rp := f.checkURL(rq)
	log.Println("rp", rp)
	rw.Header().Set("X-I2P-TORRENTLOCATION", f.magnet)
	if rp == "a" {
		client := http.Client{}
		req := f.ProxyRequest(rq)
		resp, err := client.Do(req)
		if err != nil {
			return
		}
		resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		fmt.Fprintf(rw, string(body))
		return
	}
	if strings.HasPrefix(rq.Header.Get("User-Agent"), "git") {
		log.Println(rq.Header.Get("User-Agent"))
		f.HandleFile(rw, rq)
		return
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
	if strings.HasSuffix("/"+rq.URL.Path, "/a") {

		p = "a"
		log.Println("URL path", p)
		return p
	}
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
	fp := filepath.Join(f.ServeDir, p)
	return fp
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
	log.Println("Handling Git")
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
