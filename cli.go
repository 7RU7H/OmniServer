package OmniServer

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"

	"OmniServer/util.go"
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
	s.ServerType, s.ServerID = 0
	s.NewProc = false
	s.ProcInfo.PID = -1
	s.ProcInfo.UID = os.Geteuid()
	s.ServerInfo.Status = 0
	s.ServerInfo.Hostnames = make([]string, 0, 0)
	s.ServerInfo.TotalHostnames = -1
	s.ServerInfo.ListeningPort = "-1"
	s.ServerInfo.ServerAddr = ""
	s.TLSInfo.ServerCertPath = ""
	s.TLSInfo.ServerKeyPath = ""
	s.TLSInfo.CertExpiryDays = -1
	return nil
}

func (s *Server) InitServerFromArgs(mhAssignedServerID int, ipAddressArg, hostnameStr, ifconfigArgStr, portArg, tlsArgsStr string) (err error) {
        s.ServerInfo.Status = 1 // Initial phase checks
	s.ServerID = mhAssignedServerID
        s.ServerInfo.ServerAddr = util.CheckValidIP(ipAddressArg)
	s.ServerInfo.ListeningPort = util.ConvertPortNumber(portArg) // Needs ":" prepended to the digits
	s.convertServerTypeStrToInt()
	hostnameSlice := strings.SplitN(hostnamesArgStr, ",", -1)
	s.ServerInfo.Hostnames = append(hostnameSlice)
	s.ServerInfo.TotalHostnames = len(hostnameSlice) - 1
	ifconfig, err := net.InterfaceByName(ifconfigArgStr)
	util.CheckError(err)
	s.ServerInfo.ifconfigName = ifconfigArgStr
	util.CheckError(err)
	ifconfigCIDRTmp, err := util.ConvIfconfigNameToCIDR(ifconfig, ifconfigArgStr)
	CheckError(err)
	ifconfigCIDR = ifconfigCIDRTmp

	if tlsArgsStr != "" {
		// if len(tlsArgsStr) != 0 {}
		// tlsArgs := string.SplitN(tlsArgStr, ",", -1)
		// if (isTLS == true && len(tlsArgs) != 3) { } // chill out for a second
	}
	return nil

}

func (s *Server) InitTLSFromArgs() {
	// TLS must be converted to slice
	tls := TLSInfo{}
	checkCertPath, err := util.CheckFileExists(fromArgsTlsInfo[1])
	if !checkCertPath {
		//
		return err
	} else {
		tls.ServerCertPath = fromArgsTlsInfo[1]
	}
	checkKeyPath, err := util.CheckFileExists(fromArgsTlsInfo[2])
	if !checkKeyPath {
		//
		return err
	} else {
		tls.ServerKeyPath = fromArgsTlsInfo[2]

	}
	tls.CertExpiryDays = fromArgsTlsInfo[3]
	s.TLSInfo = tls
}

func (s *Server) convertServerTypeStrToInt(userArg string) error {
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
	
        // Server Comand and flags
	var ipAddress, netInterface, serverType, tlsInputStr, hostnamesStr string
	var port int
	serverCommand := flag.NewFlagSet("server", flag.ExitOnError)
	serverCommand.StringVar(&tlsInputStr, "-t", "None", "Provide TLS information delimited by a comma - requires server type: -s https")
	serverCommand.StringVar(&serverType, "-s", "http", "Provide a server of type: http, https")
	serverCommand.StringVar(&netInterface, "-e", "localhost", "Provide a Network Interface - required!")
	serverCommand.StringVar(&ipAddress, "-i", "127.0.0.1", "Provide a valid IPv4 address - required!")
	serverCommand.StringVar(&hostnamesStr, "-H", "", "Provide DNS or vhosting alias comma delimited: dc.test.org,test.org")
	serverCommand.IntVar(&port, "-p", 8443, "Provide a TCP port number - required!")

	// Console command and flags
	consoleCommand := flag.NewFlagSet("console", flag.ExitOnError)
	// Utility flag - comment for aiding eyes
	var helpFlag, versionFlag string
	flag.StringVar(&helpFlag, "-h", "Help", "Help")
	flag.StringVar(&versionFlag, "-v", "Version", "Version")

	if err := flag.Parse(args); err != nil {
		return err
	}
	argsLen := len(args)

	if argsLen > 1 {
		flag.Usage()
		os.Exit(1)
	}

	if argsLen == 1 {
		if flag.Lookup(helpFlag) != nil {
			flag.Usage()
			os.Exit(1)
		}
		if flag.Lookup(versionFlag) != nil {
			flag.Usage()
			os.Exit(1)
		}
		if flag.Lookup(consoleCommand) != nil {
			consoleFlag = 1
		}
	}
	if flag.Lookup(serverCommand) == nil {
		err := fmt.Errorf("Invalid command passed by user: %s", args)
		return err
	}
 
        
        ifconfigArgStr :=
        portArgInt :=
        hostnameStr :=
        ipAddress := 
        tlsArgsStr :=


	// serverType
	// ipAddress
	// netinterface
	util.CheckValidInterface(flag.Lookup(netInterface))
	// port
	// tlsInputStr
	// hostnamesStr


	util.CheckError(err)
	mhAssignedServerIDCounter := metahandler.InitSelf()
	switch consoleFlag {
	case 1:
		server := Server{}
		server.setDefaultServerConfig()
		server.InitServerFromArgs(mhAssignedServerID, ipAddressArg, hostnameStr, ifconfigArgStr, portArg, tlsArgsStr)
		metahandler.SelectAction(server, consoleFlag)
	case 2:
		fmt.Println("metahandler.ToConsole()")
	case 0:
		fmt.Println("Error with checkArgs")
	default:
		fmt.Println("no value passed")
		return err
	}

	return nil
}
