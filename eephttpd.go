package eephttpd

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/eyedeekay/sam-forwarder/config"
	"github.com/eyedeekay/sam-forwarder/interface"
	"github.com/eyedeekay/sam-forwarder/tcp"
	"github.com/eyedeekay/samtracker"
	"github.com/radovskyb/watcher"
	"github.com/sosedoff/gitkit"
	"github.com/xgfone/bt/bencode"
	"github.com/xgfone/bt/metainfo"
	"gitlab.com/golang-commonmark/markdown"

	"gopkg.in/src-d/go-git.v4"
)

//EepHttpd is a structure which automatically configured the forwarding of
//a local service to i2p over the SAM API.
type EepHttpd struct {
	*samtracker.SamTracker
	*gitkit.Server
	*watcher.Watcher
	ServeDir string
	GitRepo  *git.Repository
	GitURL   string
	Hostname string
	IniFile  string
	up       bool
	pulling  bool
	magnet   string
	meta     *metainfo.MetaInfo
	mark     *markdown.Markdown
}

var err error

func (f *EepHttpd) GetType() string {
	return "eephttpd"
}

func (f *EepHttpd) ServeParent() {
	log.Println("Starting eepsite server", f.Base32())
	if err = f.SamTracker.Serve(); err != nil {
		f.Cleanup()
	}
}

func (f *EepHttpd) Target() string {
	pp, _ := strconv.Atoi(f.SamTracker.Config().TargetPort)
	return f.SamTracker.Config().TargetHost + ":" + strconv.Itoa(pp)
}

//Serve starts the SAM connection and and forwards the local host:port to i2p
func (f *EepHttpd) Serve() error {
	go f.ServeParent()
	f.MakeTorrent()
	if f.Up() {
		log.Println("Starting web server", f.Target())
		if err := http.ListenAndServe(f.Target(), f); err != nil {
			return err
		}
	}
	return nil
}

func (f *EepHttpd) Up() bool {
	return f.up
}

//Close shuts the whole thing down.
func (f *EepHttpd) Close() error {
	return f.SamTracker.Close()
}

func (s *EepHttpd) Load() (samtunnel.SAMTunnel, error) {
	if !s.up {
		log.Println("Started putting tunnel up")
	}
	s.Server = gitkit.New(gitkit.Config{
		Dir:        s.ServeDir,
		AutoCreate: true,
		Auth:       true, // Turned off by default
	})

	s.Server.AuthFunc = Never
	//log.Println("Options loaded", e.Print())

	if s.GitURL != "" {
		_, err := os.Stat(filepath.Join(s.ServeDir, ".git"))
		if os.IsNotExist(err) {
			s.GitRepo, err = git.PlainClone(s.ServeDir, false, &git.CloneOptions{
				URL:               s.GitURL,
				RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
			})
			if err != nil {
				return nil, err
			}
		} else {
			s.GitRepo, err = git.PlainOpen(s.ServeDir)
			if err != nil {
				return nil, err
			}
		}
	}
	f, e := s.SamTracker.Load()
	if e != nil {
		return nil, e
	}
	s.SamTracker = f.(*samtracker.SamTracker)

	s.mark = markdown.New(markdown.XHTMLOutput(true))
	s.up = true
	s.Watcher = watcher.New()
	s.Watcher.SetMaxEvents(1)
	s.Watcher.AddRecursive(s.ServeDir)
	go func() {
		for {
			select {
			case event := <-s.Watcher.Event:
				log.Printf("File event %v\n", event)
				err := s.MakeTorrent()
				if err != nil {
					log.Printf("File Watcher Error %e", err)
				}
			case err := <-s.Watcher.Error:
				log.Printf("File Watcher Error %e", err)
			case <-s.Watcher.Closed:
				return
			}
		}
	}()

	go func() {
		s.Watcher.Wait()
		s.Watcher.TriggerEvent(watcher.Create, nil)
		s.Watcher.TriggerEvent(watcher.Remove, nil)
	}()
	err := s.MakeTorrent()
	if err != nil {
		return nil, err
	}

	log.Println("Finished putting tunnel up")
	return s, nil
}

func (e *EepHttpd) HostName() string {
	return e.Base32()
}

func DirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}

func (e *EepHttpd) MakeTorrent() error {
	e.meta = &metainfo.MetaInfo{}
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	defer os.Chdir(wd)
	err = os.Chdir(e.ServeDir)
	if err != nil {
		return err
	}
	size, err := DirSize(e.ServeDir)
	if err != nil {
		return err
	}
	pieceLen := (size / 30000) / 2
	if pieceLen < 25600 {
		pieceLen = 25600
	}
	log.Println("Calculating optimal piece length", size, pieceLen)
	info, err := metainfo.NewInfoFromFilePath(".", int64(pieceLen))
	if err != nil {
		return err
	}
	info.Name = e.HostName()
	log.Println("Generating torrent:", info.CountPieces(), "pieces")

	//	e.meta.SetDefaults()
	e.meta.InfoBytes, err = bencode.EncodeBytes(info)
	if err != nil {
		return err
	}
	e.meta.URLList = metainfo.URLList{"http://" + e.HostName()}
	e.meta.Announce = "http://" + e.HostName() + "/a"
	e.meta.AnnounceList = metainfo.AnnounceList{[]string{"http://" + e.HostName() + "/a", "http://w7tpbzncbcocrqtwwm3nezhnnsw4ozadvi2hmvzdhrqzfxfum7wa.b32.i2p/a"}}
	e.meta.CreatedBy = "eephttpd"
	e.magnet = e.meta.Magnet(e.HostName(), e.meta.InfoHash()).String()
	return nil
}

func (e *EepHttpd) Print() string {
	pr := e.Config().Print()
	pr += "servedir=" + e.ServeDir + "\n"
	pr += "gitrepo=" + e.GitURL + "\n"
	return pr
}

func (e *EepHttpd) Save() error {
	if e.IniFile == "none" {
		e.IniFile = "eephttpd.ini"
	}
	return ioutil.WriteFile(e.IniFile, []byte(e.Print()), 0644)
}

func (e *EepHttpd) noPull() {
	e.pulling = false
}

func (e *EepHttpd) ResetGit() error {
	if e.GitURL != "" {
		defer e.MakeTorrent()
		log.Println("Resetting git repository to", e.GitURL)
		os.RemoveAll(filepath.Join(e.ServeDir))
		os.Mkdir(filepath.Join(e.ServeDir), 0755)
		_, err := os.Stat(filepath.Join(e.ServeDir, ".git"))
		if os.IsNotExist(err) {
			e.GitRepo, err = git.PlainClone(e.ServeDir, false, &git.CloneOptions{
				URL:               e.GitURL,
				RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
			})
			if err != nil {
				return err
			}
		} else {
			e.GitRepo, err = git.PlainOpen(e.ServeDir)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (e *EepHttpd) Pull() error {
	if e.pulling {
		return nil
	} else {
		e.pulling = true
		defer e.noPull()
	}

	if e.GitURL != "" {
		defer e.MakeTorrent()
		_, err := os.Stat(filepath.Join(e.ServeDir, ".git"))
		if os.IsNotExist(err) {
			e.GitRepo, err = git.PlainClone(e.ServeDir, false, &git.CloneOptions{
				URL:               e.GitURL,
				RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
			})
			if err != nil {
				return err
			}
		} else {
			e.GitRepo, err = git.PlainOpen(e.ServeDir)
			if err != nil {
				return err
			}
		}
	}

	if e.GitURL != "" {
		if e.GitRepo != nil {
			w, err := e.GitRepo.Worktree()
			if err != nil {
				return err
			}
			err = w.Checkout(&git.CheckoutOptions{})
			if err != nil {
				return err
			}
			err = w.Pull(&git.PullOptions{RemoteName: "origin"})
			if err != nil {
				return err
			}
			return nil
		} else {
			e.GitRepo, err = git.PlainOpen(e.ServeDir)
			if err != nil {
				return err
			}
			w, err := e.GitRepo.Worktree()
			if err != nil {
				return err
			}
			err = w.Checkout(&git.CheckoutOptions{})
			if err != nil {
				return err
			}
			err = w.Pull(&git.PullOptions{RemoteName: "origin"})
			if err != nil {
				return err
			}
			return nil
		}
		return nil
	}
	return nil
}

//NewEepHttpd makes a new SAM forwarder with default options, accepts host:port arguments
func NewEepHttpd(host, port string) (*EepHttpd, error) {
	return NewEepHttpdFromOptions(SetHost(host), SetPort(port))
}

func Never(gitkit.Credential, *gitkit.Request) (bool, error) {
	return false, nil
}

//NewEepHttpdFromOptions makes a new SAM forwarder with default options, accepts host:port arguments
func NewEepHttpdFromOptions(opts ...func(*EepHttpd) error) (*EepHttpd, error) {
	var s EepHttpd
	s.SamTracker = &samtracker.SamTracker{
		SAMForwarder: &samforwarder.SAMForwarder{
			Conf: &i2ptunconf.Conf{},
		},
	}
	s.Server = &gitkit.Server{}
	log.Println("Initializing eephttpd")
	for _, o := range opts {
		if err := o(&s); err != nil {
			return nil, err
		}
	}
	s.SamTracker.Config().SaveFile = true
	pp, _ := strconv.Atoi(s.SamTracker.Config().TargetPort)
	s.SamTracker.InitTarget(s.SamTracker.Config().TargetHost + ":" + strconv.Itoa(pp+1))
	//	s.tracker.SamTracker = e.SamTracker
	l, e := s.Load()
	if e != nil {
		return nil, e
	}
	return l.(*EepHttpd), nil
}
