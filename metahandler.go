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

func SelectAction(s *Server) {

}
