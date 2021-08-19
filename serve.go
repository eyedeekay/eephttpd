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
	url := fmt.Sprintf("%s://%s:%s%s", "http", f.SamTracker.Config().TargetHost, strconv.Itoa(pp+1), req.URL.Path)
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
	//	Content-Security-Policy: default-src 'self' trusted.com *.trusted.com
	csp := []string{"default-src: 'self' " + f.Base32(), "script-src: 'self' " + f.Base32()}
	rw.Header().Set("Content-Security-Policy", strings.Join(csp, "; "))
	rp := f.checkURL(rq)
	mtype, err := mimetype.DetectFile(rp)
	if err != nil {
		log.Println("MIME type determination error.", err.Error())
	} else {
		log.Println("MIME type detected", mtype.String())
	}
	rw.Header().Set("Content-Type", mtype.String())
	if strings.HasSuffix(rq.URL.Path, ".css") {
		rw.Header().Set("Content-Type", "text/css")
	}
	if strings.HasSuffix(rq.URL.Path, ".js") {
		rw.Header().Set("Content-Type", "text/javascript")
	}
	if strings.HasSuffix(rq.URL.Path, ".md") {
		rw.Header().Set("Content-Type", "text/html")
	}
	if strings.HasSuffix(rp, ".html") {
		rw.Header().Set("Content-Type", "text/html")
	}
	rw.Header().Set("X-I2P-TORRENTLOCATION", f.magnet)
	defer f.Pull()
	if rp == "torrent" {
		f.HandleTorrent(rw, rq)
	} else if rp == "announce" {
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
	f.PullFeeds()
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
	if strings.HasSuffix(rq.URL.Path, "eephttpd.torrent") {
		p = "torrent"
		return p
	}
	if strings.HasSuffix("/"+rq.URL.Path, "/a") {
		p = "announce"
		return p
	}
	if strings.HasPrefix("/"+rq.URL.Path, "/announce") {
		p = "announce"
		return p
	}
	if strings.HasSuffix("/"+rq.URL.Path, "/a/s") {
		p = "scrape"
		return p
	}
	if "/"+rq.URL.Path == "/s" {
		p = "scrape"
		return p
	}
	if strings.HasSuffix("/"+rq.URL.Path, "/announce/s") {
		p = "scrape"
		return p
	}
	if strings.HasPrefix("/"+rq.URL.Path, "/scrape") {
		p = "scrape"
		return p
	}
	if strings.HasPrefix("/"+rq.URL.Path, "/a/scrape") {
		p = "scrape"
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
	if IsDirectory(p) {
		p = filepath.Join(rq.URL.Path, "index.html")
	}
	fp := filepath.Join(f.ServeDir, p)
	return fp
}

func IsDirectory(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
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
	rw.Write(bytes)
	//	fmt.Fprintf(rw, string(bytes))
}

func (e *EepHttpd) HandleTorrent(rw http.ResponseWriter, rq *http.Request) {
	e.meta.Write(rw)
}

func (f *EepHttpd) HandleMissing(rw http.ResponseWriter, rq *http.Request) {
	path := f.checkURL(rq)
	fmt.Fprintf(rw, "ERROR %s NOT FOUND", strings.Replace(path, f.ServeDir, "", -1))
}
