package OmniServer

import (
	""
	"os"
)

type Server struct {
        ServerType int // Integer reference for each - decimalise as in 0 - 9 is (default=0, debug; 10 is webserver, 20 proxy, 30 capture - 11 is then an option for feature extension of a webserver
        ServerID int // 0 ID is temporary ID from intialiseDefault , negative digits are stopped server IDs
        ServerWithCtx *http.Server
        Ctx Context
        CancelCtx CancelFunc
        Mux *ServerMux
        ServerInfo ServerInfo
        TLSInfo TLSInfo
        NewProc bool
        ProcInfo ProcInfo
}

type ProcInfo struct {
        PID int // Set -1 by default
        UID int // os.Geteuid used by default 
}

type ServerInfo struct {
        Status int // 0 ID is set for all initialing servers till checks
        Hostnames []string
        TotalHostnames int
        ListeningPort int
        ServerAddr string
}

type TLSInfo struct {
        ServerCertPath string
        ServerKeyPath string
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
        s.ServerInfo.ListeningPort = -1
        s.ServerAddr = ""
        s.TLSInfo.ServerCertPath = ""
        s.TLSInfo.ServerKeyPath = ""
        s.TLSInfo.CertExpiryDays = -1
        return nil
}

func (s *Server) InitServerStruct(hasTLS, hasHosts, newProc bool, argsServerInfo, fromArgsTlsInfo []string) (error) {
	//  
	//
	//
	tls := TLSInfo{}
	if hasTLS {
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
	} else { 
		tls.ServerCertPath = "none"
		tls.ServerKeyPath = "none"
        	tls.CertExpiryDays = -1
		s.TLSInfo = tls
	}
	
        //EvaluateHostnames return len(arr)
	ServerInfo {
		Status = 0,
                Hostnames = append(),
                TotalHostnames = len()-1
		// hostnames =  
		// func () if !hasHosts { hostname = "" } else { hostnameList := fromArgsServerInfo[INDEX] }
		// 
	}
	
}


func checkArgs(args []string) (int, error) {
        if len(args) != 0 {

        } else {
        // No args provided
        return 0, err
        }
}

// Remember to reread oldmain.go !!
// logging has to done somewhere
func HandleArgs(args []string) error {
        consoleFlag, err := checkArgs(args)
        util.CheckError(err)
        switch consoleFlag {
        case 1:       
                server := Server{}
                server.setDefaultServerConfig()
                server.InitServerFromArgs(args)
                metahandler.SelectAction(server, consoleFlag)
        case 2: 
                fmt.Println("metahandler.RunConsole(args)")
        case 0: 
                fmt.Println("Error with checkArgs")
        default:
                fmt.Println("no value passed")
                return err
        }

        return nil
}
