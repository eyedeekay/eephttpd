package main

import (
	"crypto/tls"
	"flag"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
)

import (
	"github.com/eyedeekay/eephttpd"
	"github.com/eyedeekay/sam-forwarder/config"
)

var cfg = &tls.Config{
	MinVersion:               tls.VersionTLS12,
	CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
	PreferServerCipherSuites: true,
	CipherSuites: []uint16{
		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
		tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_RSA_WITH_AES_256_CBC_SHA,
	},
}

func wwwdir() string {
	path := "./www"
	if runtime.GOOS == "windows" {
		user, err := user.Current()
		if err != nil {
			return "./www"
		}
		path = filepath.Join(user.HomeDir, "My Documents", "www")
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0755)
	}
	return path
}

var (
	host               = flag.String("a", "127.0.0.1", "hostname to serve on")
	port               = flag.String("p", "7880", "port to serve locally on")
	samhost            = flag.String("sh", "127.0.0.1", "sam host to connect to")
	samport            = flag.String("sp", "7656", "sam port to connect to")
	directory          = flag.String("d", wwwdir(), "the directory of static files to host(default"+wwwdir()+")")
	giturl             = flag.String("b", "", "URL of a git repository to build populate the static directory with(optional)")
	usei2p             = flag.Bool("i", true, "save i2p keys(and thus destinations) across reboots")
	servicename        = flag.String("n", "eephttpd", "name to give the tunnel(default eephttpd)")
	useCompression     = flag.Bool("g", true, "Uze gzip(true or false)")
	accessListType     = flag.String("l", "none", "Type of access list to use, can be \"whitelist\" \"blacklist\" or \"none\".")
	encryptLeaseSet    = flag.Bool("c", false, "Use an encrypted leaseset(true or false)")
	allowZeroHop       = flag.Bool("z", false, "Allow zero-hop, non-anonymous tunnels(true or false)")
	reduceIdle         = flag.Bool("r", false, "Reduce tunnel quantity when idle(true or false)")
	reduceIdleTime     = flag.Int("rt", 600000, "Reduce tunnel quantity after X (milliseconds)")
	reduceIdleQuantity = flag.Int("rc", 3, "Reduce idle tunnel quantity to X (0 to 5)")
	inLength           = flag.Int("il", 3, "Set inbound tunnel length(0 to 7)")
	outLength          = flag.Int("ol", 3, "Set outbound tunnel length(0 to 7)")
	inQuantity         = flag.Int("iq", 2, "Set inbound tunnel quantity(0 to 15)")
	outQuantity        = flag.Int("oq", 2, "Set outbound tunnel quantity(0 to 15)")
	inVariance         = flag.Int("iv", 0, "Set inbound tunnel length variance(-7 to 7)")
	outVariance        = flag.Int("ov", 0, "Set outbound tunnel length variance(-7 to 7)")
	inBackupQuantity   = flag.Int("ib", 1, "Set inbound tunnel backup quantity(0 to 5)")
	outBackupQuantity  = flag.Int("ob", 1, "Set outbound tunnel backup quantity(0 to 5)")
	iniFile            = flag.String("f", "none", "Use an ini file for configuration")
	useTLS             = flag.Bool("t", false, "Generate or use an existing TLS certificate")
	certFile           = flag.String("m", "cert", "Certificate name to use")
)

var eepsite *eephttpd.EepHttpd

func main() {
	flag.Parse()
	var err error
	go UiMain()
	config := i2ptunconf.NewI2PBlankTunConf()
	if *iniFile != "none" {
		var err error
		config, err = i2ptunconf.NewI2PTunConf(*iniFile)
		if err != nil {
			log.Fatal(err)
		}
	}
	config.TargetHost = config.GetHost(*host, "127.0.0.1")
	config.TargetPort = config.GetPort(*port, "7880")
	config.SaveFile = config.GetSaveFile(*usei2p, true)
	config.SamHost = config.GetSAMHost(*samhost, "127.0.0.1")
	config.SamPort = config.GetSAMPort(*samport, "7656")
	config.TunName = config.GetKeys(*servicename, "eephttpd")
	config.InLength = config.GetInLength(*inLength, 3)
	config.OutLength = config.GetOutLength(*outLength, 3)
	config.InVariance = config.GetInVariance(*inVariance, 0)
	config.OutVariance = config.GetOutVariance(*outVariance, 0)
	config.InQuantity = config.GetInQuantity(*inQuantity, 2)
	config.OutQuantity = config.GetOutQuantity(*outQuantity, 2)
	config.InBackupQuantity = config.GetInBackups(*inBackupQuantity, 1)
	config.OutBackupQuantity = config.GetOutBackups(*outBackupQuantity, 1)
	config.EncryptLeaseSet = config.GetEncryptLeaseset(*encryptLeaseSet, false)
	config.InAllowZeroHop = config.GetInAllowZeroHop(*allowZeroHop, false)
	config.OutAllowZeroHop = config.GetOutAllowZeroHop(*allowZeroHop, false)
	config.UseCompression = config.GetUseCompression(*useCompression, true)
	config.ReduceIdle = config.GetReduceOnIdle(*reduceIdle, true)
	config.ReduceIdleTime = config.GetReduceIdleTime(*reduceIdleTime, 600000)
	config.ReduceIdleQuantity = config.GetReduceIdleQuantity(*reduceIdleQuantity, 2)
	config.AccessListType = config.GetAccessListType(*accessListType, "none")
	config.Type = config.GetTypes(false, false, false, "server")

	testdir, ok := config.Get("servedir")
	if ok {
		*directory = testdir
	}
	testgit, ok := config.Get("gitrepo")
	if ok {
		*giturl = testgit
	}

	eepsite, err = eephttpd.NewEepHttpdFromOptions(
		eephttpd.SetType(config.Type),
		eephttpd.SetSAMHost(config.SamHost),
		eephttpd.SetSAMPort(config.SamPort),
		eephttpd.SetHost(config.TargetHost),
		eephttpd.SetPort(config.TargetPort),
		eephttpd.SetSaveFile(config.SaveFile),
		eephttpd.SetName(config.TunName),
		eephttpd.SetInLength(config.InLength),
		eephttpd.SetOutLength(config.OutLength),
		eephttpd.SetInVariance(config.InVariance),
		eephttpd.SetOutVariance(config.OutVariance),
		eephttpd.SetInQuantity(config.InQuantity),
		eephttpd.SetOutQuantity(config.OutQuantity),
		eephttpd.SetInBackups(config.InBackupQuantity),
		eephttpd.SetOutBackups(config.OutBackupQuantity),
		eephttpd.SetEncrypt(config.EncryptLeaseSet),
		eephttpd.SetAllowZeroIn(config.InAllowZeroHop),
		eephttpd.SetAllowZeroOut(config.OutAllowZeroHop),
		eephttpd.SetCompress(config.UseCompression),
		eephttpd.SetReduceIdle(config.ReduceIdle),
		eephttpd.SetReduceIdleTimeMs(config.ReduceIdleTime),
		eephttpd.SetReduceIdleQuantity(config.ReduceIdleQuantity),
		eephttpd.SetAccessListType(config.AccessListType),
		eephttpd.SetAccessList(config.AccessList),
		eephttpd.SetServeDir(*directory),
		eephttpd.SetGitURL(*giturl),
		eephttpd.SetINIFile(*iniFile),
	)
	if err != nil {
		log.Fatal(err)
	}

	if eepsite != nil {
		log.Println("Starting server")
		if err = eepsite.Serve(); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Println("Unable to start, eepsite was", eepsite)
	}
}
