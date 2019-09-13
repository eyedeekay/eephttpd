package eephttpd

import (
	"log"
	"net/http"

	"github.com/eyedeekay/sam-forwarder/config"
	"github.com/eyedeekay/sam-forwarder/interface"
	"github.com/eyedeekay/sam-forwarder/tcp"
	"github.com/eyedeekay/sam3/i2pkeys"
	"gitlab.com/golang-commonmark/markdown"
)

//EepHttpd is a structure which automatically configured the forwarding of
//a local service to i2p over the SAM API.
type EepHttpd struct {
	*samforwarder.SAMForwarder
	ServeDir  string
	up        bool
	mark      *markdown.Markdown
}

var err error

func (f *EepHttpd) Config() *i2ptunconf.Conf {
	return f.SAMForwarder.Config()
}

func (f *EepHttpd) ID() string {
	return f.Config().ID()
}

func (f *EepHttpd) Keys() i2pkeys.I2PKeys {
	return f.SAMForwarder.Keys()
}

func (f *EepHttpd) Cleanup() {
	f.SAMForwarder.Cleanup()
}

func (f *EepHttpd) GetType() string {
	return f.SAMForwarder.GetType()
}

/*func (f *EepHttpd) targetForPort443() string {
	if f.TargetForPort443 != "" {
		return "targetForPort.4443=" + f.TargetHost + ":" + f.TargetForPort443
	}
	return ""
}*/

func (f *EepHttpd) Props() map[string]string {
	return f.SAMForwarder.Props()
}

func (f *EepHttpd) Print() string {
	return f.SAMForwarder.Print()
}

func (f *EepHttpd) Search(search string) string {
	return f.SAMForwarder.Search(search)
}

// Target returns the host:port of the local service you want to forward to i2p
func (f *EepHttpd) Target() string {
	return f.SAMForwarder.Target()
}

//Base32 returns the base32 address where the local service is being forwarded
func (f *EepHttpd) Base32() string {
	return f.SAMForwarder.Base32()
}

//Base32Readable returns the base32 address where the local service is being forwarded
func (f *EepHttpd) Base32Readable() string {
	return f.SAMForwarder.Base32Readable()
}

//Base64 returns the base64 address where the local service is being forwarded
func (f *EepHttpd) Base64() string {
	return f.SAMForwarder.Base64()
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
	if f.SAMForwarder.Up() {
		return true
	}
	return false
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
	log.Println("Options loaded", s.Print())
	l, e := s.Load()
	if e != nil {
		return nil, e
	}
	return l.(*EepHttpd), nil
}
