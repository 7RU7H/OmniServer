ckage OmniServer

import (
        "fmt"
        "log"
        "net/http"
        "os"
        "strconv"
        "strings"
        "time"

	"github.com/7ru7h/OmniServer/util"
)


// Control
// Hold all IDs
// Hold all PIDs
// Hold a pointer to all servers, consoles

func CreateServer(s *Server) (error)  {
        if CheckAvaliableIDs(s.ServerID) || CheckAvaliableIDs() {
                // ID in use
        }
        // ServerType == Integer reference for each - decimalise as in 0 - 9 is debug; 10 is webserver, 20 proxy, 30 capture - 11 is then an option for feature extension of a webserver
        switch s.ServerType {
                case 10: // HTTP Server
                        web.CreateWebServer()
                case 11: // HTTPS Server
                        // Handle TLS certificate generation, custom usage
                        tls.manageTLSCertInit() // pass ??.TLSInfo ->
                        web.CreateWebServer()
                default:
                        if s.ServerType <= 9 {  // Debug ServerType value
		        } else {
                                err := fmt.Errorf("Incorrect ServerType %d", s.ServerType)
                                return err
                        }
        }
        return nil
}



//
func StartServer(s *Server) (error)  {
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
               * return err
        } else if err != nil {
                fmt.Printf("Error listening for %s: %s\n", ServerID, err)
                log.Fatal("Error listening for %s - ID %d: %s\n", ServerID, err)
                return err
        } else {
                log.Printf("%s is listening...\n", ServerID)
                return err
        }
        return nil
}





// Pause server, retain memory and does not deallocate
func (s *Server) PauseServer() (error)  {
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

func (s *Server) StopServer() (error)  {
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


func GracefulExit() error {
        // For all Server, Console IDs kill each PID        
        return nil
}


func SelectAction(s *Server, actionFlag int) error {
        switch actionFlag {
        case 1:
                // StartConsole
        case 2:        
                // OpenConsole
        case 3:
                // StopConsole
        case 4:
                StartServer(s)
        case 5:
                StopServer(s)
        case 6:
                PauseServer(s)
        case 7:
                CreateServer(s)        
        default:
                err := fmt.Errorf("Incorrect actionFlag %d provide to SelectAction")
                return err
        }
        
        err := GracefulExit()
        util.CheckError(err)
        return nil
}
