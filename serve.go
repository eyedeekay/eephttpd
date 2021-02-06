package eephttpd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/d5/tengo/script"
	"github.com/gabriel-vasile/mimetype"
)

func (f *EepHttpd) ProxyRequest(req *http.Request) (*http.Request, error) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	// you can reassign the body if you need to parse it as multipart
	//    req.Body = ioutil.NopCloser(bytes.NewReader(body))

	// create a new url from the raw RequestURI sent by the client
	pp, err := strconv.Atoi(f.SamTracker.Config().TargetPort)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s://%s:%s/%s", "http", f.SamTracker.Config().TargetHost, strconv.Itoa(pp+1), "announce")
	log.Println("handling http tracker request", url)
	proxyReq, err := http.NewRequest(req.Method, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	// We may want to filter some headers, otherwise we could just use a shallow copy
	//
	proxyReq.Header = req.Header
	return proxyReq, nil
}

func (f *EepHttpd) ServeHTTP(rw http.ResponseWriter, rq *http.Request) {
	rp := f.checkURL(rq)
	mtype, err := mimetype.DetectFile(rp)
	if err != nil {
		log.Println("MIME type determination error.", err.Error())
	}
	rw.Header().Set("content-type", mtype.String())
	rw.Header().Set("X-I2P-TORRENTLOCATION", f.magnet)
	defer f.Pull()
	if rp == "announce" {
		client := http.Client{}
		req, err := f.ProxyRequest(rq)
		if err != nil {
			log.Println(err.Error())
			return
		}
		defer req.Body.Close()
		resp, err := client.Do(req)
		if err != nil {
			log.Println(err.Error())
			return
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err.Error())
			return
		}
		fmt.Fprintf(rw, string(body))
		return

	} else {
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
		p = "announce"
		log.Println("URL path", p)
		return p
	}
	if strings.HasSuffix("/"+rq.URL.Path, "/announce") {
		p = "announce"
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
