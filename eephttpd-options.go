package eephttpd

import (
	"fmt"
	"path/filepath"
	"strconv"
)

//Option is a EepHttpd Option
type Option func(*EepHttpd) error

//SetFilePath sets the path to save the config file at.
func SetFilePath(s string) func(*EepHttpd) error {
	return func(c *EepHttpd) error {
		c.Forwarder.Config().FilePath = s
		return nil
	}
}

//SetType sets the type of the forwarder server
func SetType(s string) func(*EepHttpd) error {
	return func(c *EepHttpd) error {
		if s == "http" {
			c.Forwarder.Config().Type = s
			return nil
		} else {
			c.Forwarder.Config().Type = "server"
			return nil
		}
	}
}

//SetSigType sets the type of the forwarder server
func SetSigType(s string) func(*EepHttpd) error {
	return func(c *EepHttpd) error {
		if s == "" {
			c.Forwarder.Config().SigType = ""
		} else if s == "DSA_SHA1" {
			c.Forwarder.Config().SigType = "DSA_SHA1"
		} else if s == "ECDSA_SHA256_P256" {
			c.Forwarder.Config().SigType = "ECDSA_SHA256_P256"
		} else if s == "ECDSA_SHA384_P384" {
			c.Forwarder.Config().SigType = "ECDSA_SHA384_P384"
		} else if s == "ECDSA_SHA512_P521" {
			c.Forwarder.Config().SigType = "ECDSA_SHA512_P521"
		} else if s == "EdDSA_SHA512_Ed25519" {
			c.Forwarder.Config().SigType = "EdDSA_SHA512_Ed25519"
		} else {
			c.Forwarder.Config().SigType = "EdDSA_SHA512_Ed25519"
		}
		return nil
	}
}

//SetSaveFile tells the router to save the tunnel's keys long-term
func SetSaveFile(b bool) func(*EepHttpd) error {
	return func(c *EepHttpd) error {
		c.Forwarder.Config().SaveFile = b
		return nil
	}
}

//SetHost sets the host of the service to forward
func SetHost(s string) func(*EepHttpd) error {
	return func(c *EepHttpd) error {
		c.Forwarder.Config().TargetHost = s
		return nil
	}
}

//SetPort sets the port of the service to forward
func SetPort(s string) func(*EepHttpd) error {
	return func(c *EepHttpd) error {
		port, err := strconv.Atoi(s)
		if err != nil {
			return fmt.Errorf("Invalid TCP Server Target Port %s; non-number ", s)
		}
		if port < 65536 && port > -1 {
			c.Forwarder.Config().TargetPort = s
			return nil
		}
		return fmt.Errorf("Invalid port")
	}
}

//SetSAMHost sets the host of the EepHttpd's SAM bridge
func SetSAMHost(s string) func(*EepHttpd) error {
	return func(c *EepHttpd) error {
		c.Forwarder.Config().SamHost = s
		return nil
	}
}

//SetSAMPort sets the port of the EepHttpd's SAM bridge using a string
func SetSAMPort(s string) func(*EepHttpd) error {
	return func(c *EepHttpd) error {
		port, err := strconv.Atoi(s)
		if err != nil {
			return fmt.Errorf("Invalid SAM Port %s; non-number", s)
		}
		if port < 65536 && port > -1 {
			c.Forwarder.Config().SamPort = s
			return nil
		}
		return fmt.Errorf("Invalid port")
	}
}

//SetName sets the host of the EepHttpd's SAM bridge
func SetName(s string) func(*EepHttpd) error {
	return func(c *EepHttpd) error {
		c.Forwarder.Config().TunName = s
		return nil
	}
}

//SetInLength sets the number of hops inbound
func SetInLength(u int) func(*EepHttpd) error {
	return func(c *EepHttpd) error {
		if u < 7 && u >= 0 {
			c.Forwarder.Config().InLength = u
			return nil
		}
		return fmt.Errorf("Invalid inbound tunnel length")
	}
}

//SetOutLength sets the number of hops outbound
func SetOutLength(u int) func(*EepHttpd) error {
	return func(c *EepHttpd) error {
		if u < 7 && u >= 0 {
			c.Forwarder.Config().OutLength = u
			return nil
		}
		return fmt.Errorf("Invalid outbound tunnel length")
	}
}

//SetInVariance sets the variance of a number of hops inbound
func SetInVariance(i int) func(*EepHttpd) error {
	return func(c *EepHttpd) error {
		if i < 7 && i > -7 {
			c.Forwarder.Config().InVariance = i
			return nil
		}
		return fmt.Errorf("Invalid inbound tunnel length")
	}
}

//SetOutVariance sets the variance of a number of hops outbound
func SetOutVariance(i int) func(*EepHttpd) error {
	return func(c *EepHttpd) error {
		if i < 7 && i > -7 {
			c.Forwarder.Config().OutVariance = i
			return nil
		}
		return fmt.Errorf("Invalid outbound tunnel variance")
	}
}

//SetInQuantity sets the inbound tunnel quantity
func SetInQuantity(u int) func(*EepHttpd) error {
	return func(c *EepHttpd) error {
		if u <= 16 && u > 0 {
			c.Forwarder.Config().InQuantity = u
			return nil
		}
		return fmt.Errorf("Invalid inbound tunnel quantity")
	}
}

//SetOutQuantity sets the outbound tunnel quantity
func SetOutQuantity(u int) func(*EepHttpd) error {
	return func(c *EepHttpd) error {
		if u <= 16 && u > 0 {
			c.Forwarder.Config().OutQuantity = u
			return nil
		}
		return fmt.Errorf("Invalid outbound tunnel quantity")
	}
}

//SetInBackups sets the inbound tunnel backups
func SetInBackups(u int) func(*EepHttpd) error {
	return func(c *EepHttpd) error {
		if u < 6 && u >= 0 {
			c.Forwarder.Config().InBackupQuantity = u
			return nil
		}
		return fmt.Errorf("Invalid inbound tunnel backup quantity")
	}
}

//SetOutBackups sets the inbound tunnel backups
func SetOutBackups(u int) func(*EepHttpd) error {
	return func(c *EepHttpd) error {
		if u < 6 && u >= 0 {
			c.Forwarder.Config().OutBackupQuantity = u
			return nil
		}
		return fmt.Errorf("Invalid outbound tunnel backup quantity")
	}
}

//SetEncrypt tells the router to use an encrypted leaseset
func SetEncrypt(b bool) func(*EepHttpd) error {
	return func(c *EepHttpd) error {
		if b {
			c.Forwarder.Config().EncryptLeaseSet = true
			return nil
		}
		c.Forwarder.Config().EncryptLeaseSet = false
		return nil
	}
}

//SetServeDir sets the path to the directory you want to serve
func SetServeDir(s string) func(*EepHttpd) error {
	return func(c *EepHttpd) error {
		var err error
		c.ServeDir, err = filepath.Abs(s)
		if err != nil {
			return err
		}
		return nil
	}
}

//SetLeaseSetKey sets the host of the EepHttpd's SAM bridge
func SetLeaseSetKey(s string) func(*EepHttpd) error {
	return func(c *EepHttpd) error {
		c.Forwarder.Config().LeaseSetKey = s
		return nil
	}
}

//SetLeaseSetPrivateKey sets the host of the EepHttpd's SAM bridge
func SetLeaseSetPrivateKey(s string) func(*EepHttpd) error {
	return func(c *EepHttpd) error {
		c.Forwarder.Config().LeaseSetPrivateKey = s
		return nil
	}
}

//SetLeaseSetPrivateSigningKey sets the host of the EepHttpd's SAM bridge
func SetLeaseSetPrivateSigningKey(s string) func(*EepHttpd) error {
	return func(c *EepHttpd) error {
		c.Forwarder.Config().LeaseSetPrivateSigningKey = s
		return nil
	}
}

//SetMessageReliability sets the host of the EepHttpd's SAM bridge
func SetMessageReliability(s string) func(*EepHttpd) error {
	return func(c *EepHttpd) error {
		c.Forwarder.Config().MessageReliability = s
		return nil
	}
}

//SetAllowZeroIn tells the tunnel to accept zero-hop peers
func SetAllowZeroIn(b bool) func(*EepHttpd) error {
	return func(c *EepHttpd) error {
		if b {
			c.Forwarder.Config().InAllowZeroHop = true
			return nil
		}
		c.Forwarder.Config().InAllowZeroHop = false
		return nil
	}
}

//SetAllowZeroOut tells the tunnel to accept zero-hop peers
func SetAllowZeroOut(b bool) func(*EepHttpd) error {
	return func(c *EepHttpd) error {
		if b {
			c.Forwarder.Config().OutAllowZeroHop = true
			return nil
		}
		c.Forwarder.Config().OutAllowZeroHop = false
		return nil
	}
}

//SetCompress tells clients to use compression
func SetCompress(b bool) func(*EepHttpd) error {
	return func(c *EepHttpd) error {
		if b {
			c.Forwarder.Config().UseCompression = true
			return nil
		}
		c.Forwarder.Config().UseCompression = false
		return nil
	}
}

//SetFastRecieve tells clients to use compression
func SetFastRecieve(b bool) func(*EepHttpd) error {
	return func(c *EepHttpd) error {
		if b {
			c.Forwarder.Config().FastRecieve = true
			return nil
		}
		c.Forwarder.Config().FastRecieve = false
		return nil
	}
}

//SetReduceIdle tells the connection to reduce it's tunnels during extended idle time.
func SetReduceIdle(b bool) func(*EepHttpd) error {
	return func(c *EepHttpd) error {
		if b {
			c.Forwarder.Config().ReduceIdle = true
			return nil
		}
		c.Forwarder.Config().ReduceIdle = false
		return nil
	}
}

//SetReduceIdleTime sets the time to wait before reducing tunnels to idle levels
func SetReduceIdleTime(u int) func(*EepHttpd) error {
	return func(c *EepHttpd) error {
		c.Forwarder.Config().ReduceIdleTime = 300000
		if u >= 6 {
			c.Forwarder.Config().ReduceIdleTime = (u * 60) * 1000
			return nil
		}
		return fmt.Errorf("Invalid reduce idle timeout(Measured in minutes) %v", u)
	}
}

//SetReduceIdleTimeMs sets the time to wait before reducing tunnels to idle levels in milliseconds
func SetReduceIdleTimeMs(u int) func(*EepHttpd) error {
	return func(c *EepHttpd) error {
		c.Forwarder.Config().ReduceIdleTime = 300000
		if u >= 300000 {
			c.Forwarder.Config().ReduceIdleTime = u
			return nil
		}
		return fmt.Errorf("Invalid reduce idle timeout(Measured in milliseconds) %v", u)
	}
}

//SetReduceIdleQuantity sets minimum number of tunnels to reduce to during idle time
func SetReduceIdleQuantity(u int) func(*EepHttpd) error {
	return func(c *EepHttpd) error {
		if u < 5 {
			c.Forwarder.Config().ReduceIdleQuantity = u
			return nil
		}
		return fmt.Errorf("Invalid reduce tunnel quantity")
	}
}

//SetCloseIdle tells the connection to close it's tunnels during extended idle time.
func SetCloseIdle(b bool) func(*EepHttpd) error {
	return func(c *EepHttpd) error {
		if b {
			c.Forwarder.Config().CloseIdle = true
			return nil
		}
		c.Forwarder.Config().CloseIdle = false
		return nil
	}
}

//SetCloseIdleTime sets the time to wait before closing tunnels to idle levels
func SetCloseIdleTime(u int) func(*EepHttpd) error {
	return func(c *EepHttpd) error {
		c.Forwarder.Config().CloseIdleTime = 300000
		if u >= 6 {
			c.Forwarder.Config().CloseIdleTime = (u * 60) * 1000
			return nil
		}
		return fmt.Errorf("Invalid close idle timeout(Measured in minutes) %v", u)
	}
}

//SetCloseIdleTimeMs sets the time to wait before closing tunnels to idle levels in milliseconds
func SetCloseIdleTimeMs(u int) func(*EepHttpd) error {
	return func(c *EepHttpd) error {
		c.Forwarder.Config().CloseIdleTime = 300000
		if u >= 300000 {
			c.Forwarder.Config().CloseIdleTime = u
			return nil
		}
		return fmt.Errorf("Invalid close idle timeout(Measured in milliseconds) %v", u)
	}
}

//SetAccessListType tells the system to treat the accessList as a whitelist
func SetAccessListType(s string) func(*EepHttpd) error {
	return func(c *EepHttpd) error {
		if s == "whitelist" {
			c.Forwarder.Config().AccessListType = "whitelist"
			return nil
		} else if s == "blacklist" {
			c.Forwarder.Config().AccessListType = "blacklist"
			return nil
		} else if s == "none" {
			c.Forwarder.Config().AccessListType = ""
			return nil
		} else if s == "" {
			c.Forwarder.Config().AccessListType = ""
			return nil
		}
		return fmt.Errorf("Invalid Access list type(whitelist, blacklist, none)")
	}
}

//SetAccessList tells the system to treat the accessList as a whitelist
func SetAccessList(s []string) func(*EepHttpd) error {
	return func(c *EepHttpd) error {
		if len(s) > 0 {
			for _, a := range s {
				c.Forwarder.Config().AccessList = append(c.Forwarder.Config().AccessList, a)
			}
			return nil
		}
		return nil
	}
}

//SetTargetForPort sets the port of the EepHttpd's SAM bridge using a string
/*func SetTargetForPort443(s string) func(*EepHttpd) error {
	return func(c *EepHttpd) error {
		port, err := strconv.Atoi(s)
		if err != nil {
			return fmt.Errorf("Invalid Target Port %s; non-number ", s)
		}
		if port < 65536 && port > -1 {
			c.Forwarder.Config().TargetForPort443 = s
			return nil
		}
		return fmt.Errorf("Invalid port")
	}
}
*/

//SetKeyFile sets
func SetKeyFile(s string) func(*EepHttpd) error {
	return func(c *EepHttpd) error {
		c.Forwarder.Config().KeyFilePath = s
		return nil
	}
}
