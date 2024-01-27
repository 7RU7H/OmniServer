package OmniServer

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
)

type Server struct {
	ServerType    int // Integer reference for each - decimalise as in 0 - 9 is (default=0, debug; 10 is webserver, 20 proxy, 30 capture - 11 is then an option for feature extension of a webserver
	ServerID      int // 0 ID is temporary ID from intialiseDefault , negative digits are stopped server IDs
	ServerWithCtx *http.Server
	Ctx           Context
	CancelCtx     CancelFunc
	Mux           *ServerMux
	ServerInfo    ServerInfo
	TLSInfo       TLSInfo
	NewProc       bool
	ProcInfo      ProcInfo
}

type ProcInfo struct {
	PID int // Set -1 by default
	UID int // os.Geteuid used by default
}

type ServerInfo struct {
	Status         int      // Track lifecycle: 0:default;1:initialisation;2:metahandling; - nice for checking
	Hostnames      []string // append always
	TotalHostnames int
	ListeningPort  string
	ServerAddr     string
	ifconfigName   string
	ifconfigCIDR   string
}

// OLDMAIN AND TLS.go HAS MORE PLEASE READ
type TLSInfo struct {
	ServerCertPath string
	ServerKeyPath  string
	CertExpiryDays int
}

func (s *Server) setDefaultServerConfig() error {
	s.ServerType, s.ServerID, s.ServerInfo.Status = 0, 0, 0
	s.NewProc = false
	s.ProcInfo.UID = os.Geteuid()
	s.ServerInfo.Hostnames = make([]string, 0, 0)
	s.ServerInfo.TotalHostnames, s.ProcInfo.PID, s.TLSInfo.CertExpiryDays = -1, -1, -1
	s.ServerInfo.ListeningPort = "-1"
	s.TLSInfo.ServerCertPath, s.ServerInfo.ServerAddr, s.TLSInfo.ServerKeyPath = "", "", ""
	return nil
}

func (s *Server) InitServerFromArgs(mhAssignedServerID int, serverTypeArg ipAddressArg, hostnameArgStr, ifconfigName, ifconfigCIDR, portArg, tlsArgsStr string) (err error) {
	s.ServerID = mhAssignedServerID
	s.convServerTypeItoa(serverTypeArg)
	s.ServerInfo.Status = 1 // Initial phase checks
	s.ServerInfo.ServerAddr = util.CheckValidIP(ipAddressArg)
	s.ServerInfo.ListeningPort = util.ConvertPortNumber(portArg) // Needs ":" prepended to the digits
	hostnameSlice := strings.SplitN(hostnamesArgStr, ",", -1)
	s.ServerInfo.Hostnames = append(hostnameSlice)
	s.ServerInfo.TotalHostnames = len(hostnameSlice) - 1
	s.ServerInfo.ifconfigName = ifconfigName
	s.ServerInfo.ifconfigCIDR = ifconfigCIDR

	if tlsArgsStr != "" {
		err := s.InitTLSFromArgs(tlsArgsStr)
		util.CheckError(err)
		// if len(tlsArgsStr) != 0 {}
		// tlsArgs := string.SplitN(tlsArgStr, ",", -1)
		// if (isTLS == true && len(tlsArgs) != 3) { } // chill out for a second
	}
	return nil

}

func (s *Server) InitTLSFromArgs(tlsArgStr string) error {
	fromArgsTlsInfo := strings.Split(tlsArgStr, ",")
	if len(fromArgsTlsInfo) != 3 {
		err := fmt.Errorf("Size of fromArgsTlsInfo - sliced from tlsArgStr: %s, contains %d, of the required 3 arguments", tlsArgStr, len(fromArgsTlsInfo)-1)
	}
	tls := TLSInfo{}
	checkCertPath, err := util.CheckFileExists(fromArgsTlsInfo[0])
	if !checkCertPath {
		//
		return err
	} else {
		tls.ServerCertPath = fromArgsTlsInfo[0]
	}
	checkKeyPath, err := util.CheckFileExists(fromArgsTlsInfo[1])
	if !checkKeyPath {
		//
		return err
	} else {
		tls.ServerKeyPath = fromArgsTlsInfo[1]

	}
	// Default days
	// strconv.AtoI
	tls.CertExpiryDays = fromArgsTlsInfo[2]
	s.TLSInfo = tls
	return nil
}

func (s *Server) convServerTypeItoa(userArg string) error {
	switch userArg {
	case "http":
		s.ServerType = 10
	case "https":
		s.ServerType = 15
	default:
		err := fmt.Errorf("Invalid server type provided by user: %s", userArg)
		return err
	}
	return nil
}

// Remember to reread oldmain.go !!
// logging has to done somewhere
func HandleArgs(args []string) error {
	var consoleFlag int
	const correctTlsArgsCount int = 3

	// Server Comand and flags
	var ipAddress, netInterface, serverType, tlsInputStr, hostnamesStr string
	var portInt int
	serverCommand := flag.NewFlagSet("server", flag.ExitOnError)
	serverCommand.StringVar(&tlsInputStr, "-t", "None", "Provide TLS information delimited by a comma - requires server type: -s https")
	serverCommand.StringVar(&serverType, "-s", "http", "Provide a server of type: http, https")
	serverCommand.StringVar(&netInterface, "-e", "localhost", "Provide a Network Interface - required!")
	serverCommand.StringVar(&ipAddress, "-i", "127.0.0.1", "Provide a valid IPv4 address - required!")
	serverCommand.StringVar(&hostnamesStr, "-H", "", "Provide DNS or vhosting alias comma delimited: dc.test.org,test.org")
	serverCommand.IntVar(&portInt, "-p", 8443, "Provide a TCP port number - required!")

	// Console command and flags
	consoleCommand := flag.NewFlagSet("console", flag.ExitOnError)
	// Utility flag - comment for aiding eyes
	var helpFlag, versionFlag string
	flag.StringVar(&helpFlag, "-h", "Help", "Help")
	flag.StringVar(&versionFlag, "-v", "Version", "Version")

	flag.Parse()
	argsLen := len(args)

	if argsLen > 1 {
		flag.Usage()
		os.Exit(1)
	}

	if argsLen < 2 {
		if flag.Lookup(helpFlag) != nil {
			flag.Usage()
			os.Exit(1)
		}
		if flag.Lookup(versionFlag) != nil {
			flag.Usage()
			os.Exit(1)
		}
	} else {
		fmt.Printf("Invalid command - either server or console")
		flag.Usage()
		os.Exit(1)
	}
	var serverTypeArg, netInterfaceName, netInterfaceCDIR, ipAddressArg, tlsInputStr, hostnameArgStr string
	var portIntArg int
	switch args[1] {
	case "server":
		consoleFlag = 0
		serverCommand.Parse(args[2:])
		if flag = flag.Lookup(serverType); flag == nil || serverType != "http" && serverType != "https" {
			flag.Usage()
			err := fmt.Errorf("Invalid ServerType flag -s <server type> passed by user")
			return err
		} else {
			serverTypeArg = flag.Value.(flag.Getter).Get().(string)
		}
		if flag = flag.Lookup(netInterface); flag == nil || netInterface == "" {
			flag.Usage()
			err := fmt.Errorf("Missing Network Interface -e <network interface>")
			return err
		} else {
			netInterfaceName = flag.Value.(flag.Getter).Get().(string)
			ifconfig, err := net.InterfaceByName(netInterfaceArg)
			util.CheckError(err)
			ifconfigCIDRTmp, err := util.ConvIfconfigNameToCIDR(ifconfig, netInterfaceName)
			util.CheckError(err)
			netInterfaceCDIR = ifconfigCIDRTmp
		}
		if flag = flag.Lookup(portInt); flag == nil {
			flag.Usage()
			err := fmt.Errorf("Invalid TCP port number -p <port number>")
			return err
		} else {
			portIntArg = flag.Value.(flag.Getter).Get().(int)
		}
		if flag = flag.Lookup(hostnamesStr); flag == nil || hostnamesStr == "" {
			flag.Usage()
			err := fmt.Errorf("Missing DNS or vhosting alias -H <alias>")
			return err
		} else {
			hostnamesStrArg = flag.Value.(flag.Getter).Get().(string)
		}
		if flag = flag.Lookup(ipAddress); flag == nil || ipAddress == "" {
			flag.Usage()
			err := fmt.Errorf("Missing valid IPv4 address -i <ip address>")
			return err
		} else {
			ipAddressArg = flag.Value.(flag.Getter).Get().(string)
			if !util.CheckValidIP(ipAddressArg) {
				flag.Usage()
				err := fmt.Errorf("Invalid IP address provided: %s".ipAddressArg)
				return err
			}
		}
		if flag = flag.Lookup(tlsInputStr); flag != nil {
			tlsInputStr = flag.Value.(flag.Getter).Get().(string)
			if flag = flag.Lookup(serverType); flag == nil && flag.Value.(flag.Getter).Get.(string) != "https" {
				flag.Usage()
				err := fmt.Errorf("TLS information provided, but server type -s is not https")
				return err
			}
			if strings.Count(tlsInputStr, ",") != correctTlsArgsCount {
				flag.Usage()
				err := fmt.Errorf("TLS information provided does not contain correct number of comma delimited arguments: %s", tlsInputStr)
				return err
			}
			// TLS library checks can occur here!
		} else {
			tlsStrArg = ""
		}
	case "console":
		consoleFlag = 1
	default:
		flag.Usage()
		err := fmt.Errorf("Bad command given")
		return err
	}

	util.CheckError(err)
	mc := MetaControl{
		allPIDs:       make([]int, 0, 0),
		allServerPtr:  make([]*Server, 0),
		allServerIDs:  make([]int, 0, 0),
		serverCounter: 0,
	}
	switch consoleFlag {
	case 1:
		server := Server{}
		server.setDefaultServerConfig()
		server.InitServerFromArgs(mc.serverCounter, serverTypeArg, ipAddressArg, hostnameArgStr, netInterfaceName, netInterfaceCDIR, portIntArg, tlsInputStr)
		mc.SelectAction(server, consoleFlag)
	case 2:
		fmt.Println("metahandler.ToConsole()")
	case 0:
		fmt.Println("Error with checkArgs")
	default:
		err := fmt.Errorf("No consoleFlag: %d value passed", consoleFlag)
		return err
	}

	return nil
}
