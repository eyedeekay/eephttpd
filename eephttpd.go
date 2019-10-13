package eephttpd

import (
	"log"
	"net/http"

	"github.com/eyedeekay/sam-forwarder/interface"
	"github.com/eyedeekay/sam-forwarder/tcp"
	"gitlab.com/golang-commonmark/markdown"
)

//EepHttpd is a structure which automatically configured the forwarding of
//a local service to i2p over the SAM API.
type EepHttpd struct {
	*samforwarder.SAMForwarder
	ServeDir string
	up       bool
	mark     *markdown.Markdown
}

var err error

func (f *EepHttpd) GetType() string {
	return "eephttpd"
}

func (f *EepHttpd) ServeParent() {
	log.Println("Starting eepsite server", f.Base32())
	if err = f.SAMForwarder.Serve(); err != nil {
		f.Cleanup()
	}
}

//Serve starts the SAM connection and and forwards the local host:port to i2p
func (f *EepHttpd) Serve() error {
	go f.ServeParent()
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
	return f.SAMForwarder.Close()
}

func (s *EepHttpd) Load() (samtunnel.SAMTunnel, error) {
	if !s.up {
		log.Println("Started putting tunnel up")
	}
	f, e := s.SAMForwarder.Load()
	if e != nil {
		return nil, e
	}
	s.SAMForwarder = f.(*samforwarder.SAMForwarder)
	s.mark = markdown.New(markdown.XHTMLOutput(true))
	s.up = true
	log.Println("Finished putting tunnel up")
	return s, nil
}

//NewEepHttpd makes a new SAM forwarder with default options, accepts host:port arguments
func NewEepHttpd(host, port string) (*EepHttpd, error) {
	return NewEepHttpdFromOptions(SetHost(host), SetPort(port))
}

//NewEepHttpdFromOptions makes a new SAM forwarder with default options, accepts host:port arguments
func NewEepHttpdFromOptions(opts ...func(*EepHttpd) error) (*EepHttpd, error) {
	var s EepHttpd
	s.SAMForwarder = &samforwarder.SAMForwarder{}
	log.Println("Initializing eephttpd")
	for _, o := range opts {
		if err := o(&s); err != nil {
			return nil, err
		}
	}
	s.SAMForwarder.Config().SaveFile = true
	l, e := s.Load()
    //log.Println("Options loaded", s.Print())
	if e != nil {
		return nil, e
	}
	return l.(*EepHttpd), nil
}
