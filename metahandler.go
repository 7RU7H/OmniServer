package omniServer

import (
        "fmt"
        "log"
        "net/http"
        "os"
        "strconv"
        "strings"
        "time"

        "github.com/7ru7h/Shadow/omniServer/tls.go"
	"github.com/7ru7h/Shadow/omniServer/util.go"
)

		
type Server struct {
        ServerType int // Integer reference for each - decimalise as in 0 - 9 is debug; 10 is webserver, 20 proxy, 30 capture - 11 is then an option for feature extension of a webserver
        ServerID int // 0 ID is temporary ID till checks, negative digits are stopped server IDs
        ServerWithCtx *http.Server
        Ctx Context
        CancelCtx CancelFunc
        Mux *ServerMux
        ServerInfo ServerInfo
        TLSInfo TLSInfo
        NewProc bool
        ProcInfo ProcInfo
}

// For if Server is required to be run as a new process
type ProcInfo struct {
        PID int
        UID int
}

type ServerInfo struct {
        Status int
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

// 0 ID is set for all initialing servers till checks
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
	
	s.ServerID = 0

        //EvaluateHostnames return len(arr)
	ServerInfo {
		Status = 0,
                Hostnames = ,
                TotalHostnames = len()-1
		// hostnames =  
		// func () if !hasHosts { hostname = "" } else { hostnameList := fromArgsServerInfo[INDEX] }
		// 
	}
	
}

func (s *Server) CreateServer() (error)  {
        if CheckAvaliableIDs(s.ServerID) || CheckAvaliableIDs() {
                // ID in use
        }
        // ServerType == Integer reference for each - decimalise as in 0 - 9 is debug; 10 is webserver, 20 proxy, 30 capture - 11 is then an option for feature extension of a webserver
        switch s.ServerType {
                case 10: // HTTP Server
                        s.CreateWebServer()
                case 11: // HTTPS Server
                        // Handle TLS certificate generation, custom usage
                        tls.manageTLSCertInit() // pass ??.TLSInfo ->
                        s.CreateWebServer()
                default:
                        if s.ServerType <= 9 {  // Debug ServerType value
		        }
                        // Incorrect s.ServerType

        }

}



//
func (s *Server) StartServer() (error)  {
        if !CheckAvaliableIDs(s.ServerID) {

                // Error no server ID to
        }

        if !s.NewProc {


        } else {
        // Create new process
        TestProcInfo := ProcInfo{}
        // Check errors or assign
        s.ProcInfo = TestProcInfo
        }

        if errors.Is(err, http.ErrServerClosed) {
                fmt.Printf("%s closed\n", ServerID, err)
                log.Fatal("%s closed\n", ServerID, err)
                return err
        } else if err != nil {
                fmt.Printf("Error listening for %s: %s\n", ServerID, err)
                log.Fatal("Error listening for %s - ID %d: %s\n", ServerID, err)
                return err
        } else {
                log.Printf("%s is listening...\n", ServerID)
                return err
        }

}

// ConfigServer?
// OPcode = modify Info,  
func ( *) ConfigServer(s *Server) (error)  {
	switch s.ServerType {
	case 10: // DefaultWebServer
		// By s.ServerID, OPcode?
	default:
		if s.ServerType < 10 {
		// Debug ServerType value
		}

		// Invalid server type
	}
}



// Pause server, retain memory and does not deallocate
func (s *Server) StopServer() (error)  {
        if !CheckAvaliableIDs(s.ServerID) {
                // Error no server ID to
        }
}

// What does restart mean and why? - Recreate Context and reassign memory etc
func (s *Server) RestartServer() (error)  {
        if !CheckAvaliableIDs(s.ServerID) {

                // Error no server ID to
        }
}

// CloseServer
func (s *Server) CloseServer() (error)  {
        if !CheckAvaliableIDs(s.ServerID) {
                // Error no server ID to

        }

        // Context termination
        s.CancelCtx()
        <-s.Ctx.Done()
        ServerTerminationTime := time.Now()
        // Checks on termination

        return ServerTerminationTime, time.Now()
}


// manager/handler
// server 




func CheckArgs(args []string) error {

}

// Remember to reread oldmain.go !!
// logging has to done somewhere
func HandleAll(args []string) error {
       	err := checkArgs(args)
	util.CheckError(err)
	server := Server{}
	server.InitServerStruct() // WTF are these args
	// Selection


        // GracefulExit()
	
	return nil
}
