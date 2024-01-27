package OmniServer

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/7ru7h/OmniServer/util"
)

type MetaControl struct {
	allPIDs       []int
	allServerIDs  []int
	allServerPtr  []*Server
	serverCounter int
}

func (m *MetaControl) AddServer(s *Server) error {
	i, err := m.FindServerID(s.ServerID, s.ServerStatus)
	CheckError(err)
	if i != -1 {
		CheckError(err)
		return err
	}
	m.serverCounter = +1
	s.ServerID = m.serverCounter

	m.allServerPtr = append(m.allServerPtr, &s)
	return nil
}

func (m *MetaControl) CheckServerExists(s *Server) error {
	i, err := m.FindServerID(s.ServerID, s.ServerStatus)
	CheckError(err)
	if i != -1 {
		CheckError(err)
		return err
	}
	m.FindPID()
	// Compare Server at &i and s
	//m.allServerPtr[i]
}

func (m *MetaControl) DeleteServer(s *Server) error {
	m.FindPID(s.ProcInfo.PID)

	i, err := m.FindServerID(s.ServerID, s.ServerStatus)
	if i == -9999 {
		CheckError(err)
		return err
	}
	delete(m.ServerID[i])
	s.ServerID = 0
	return nil
}

func (m *MetaControl) FindServerID(id int) (int, error) {
	for i, value := range m.allServerIDs {
		if value == id {
			return i, nil
		}
	}
	err := fmt.Errorf("Server ID: %d not found", id)
	return -9999, err
}

// Check is bad:
// Find & New & Delete
func (m *MetaControl) FindPID(pid int) error {
	// ps aux && handle grep
	//
}

// AddConsole
// CheckConsoleExists
// DeleteConsole
//
// Console Update
//
// func CreateConsole() error {}
// func StartConsole() error {}
// func OpenConsole() error {}
// func StopConsole() error {}
//

func (m *MetaControl) StartServer(s *Server) error {
	i, err := m.FindServerID(s.ServerID)
	if i == -9999 {
		err := fmt.Errorf("No ServerID found for %d", s.ServerID)
		return err
	}
	// update ServerStatus to pending
	//

	if !s.NewProc {

	} else {
		// Create new process
		TestProcInfo := ProcInfo{}
		// Check errors or assign
		s.ProcInfo = TestProcInfo
	}

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("%s closed\n", s.ServerID, err)
		log.Fatal("%s closed\n", s.ServerID, err)
		return err
	} else if err != nil {
		fmt.Printf("Error listening for %s: %s\n", s.ServerID, err)
		log.Fatal("Error listening for %s - ID %d: %s\n", s.ServerID, err)
		return err
	} else {
		log.Printf("%s is listening...\n", s.ServerID)
		return err
	}
}

// Pause server, retain memory and does not deallocate
func (m *MetaControl) PauseServer(s *Server) error {
	err := m.CheckServerExists(s)
	CheckError(err)
	//
	return nil
}

func (m *MetaControl) ResumeServer(s *Server) error {
	err := m.CheckServerExists(s)
	CheckError(err)
	//
	return nil
}

func (m *MetaControl) StopServer(s *Server) (time.Time, time.Time, error) {
	err := m.CheckServerExists(s)
	CheckError(err)

	// Context termination
	s.CancelCtx()
	<-s.Ctx.Done()
	ServerTerminationTime := time.Now()
	// Checks on termination

	return ServerTerminationTime, time.Now(), nil
}

// What does restart mean and why? - Recreate Context and reassign memory etc
func RestartServer(s *Server) error {

	return nil
}

func (m *MetaControl) CreateServer(s *Server) error {
	m.AddServer()
	// ServerType == Integer reference for each - decimalise as in 0 - 9 is debug; 10 is webserver, 20 proxy, 30 capture - 11 is then an option for feature extension of a webserver
	switch s.ServerType {
	case 10: // HTTP Server
		web.CreateWebServer()
	case 11: // HTTPS Server
		// Handle TLS certificate generation, custom usage
		tls.manageTLSCertInit() // pass ??.TLSInfo ->
		web.CreateWebServer()
	default:
		if s.ServerType <= 9 { // Debug ServerType value
		} else {
			err := fmt.Errorf("Incorrect ServerType %d", s.ServerType)
			return err
		}
	}
	return nil
}

func GracefulExit() error {
	// For all Server, Console IDs kill each PID
	return nil
}

func (mc *MetaControl) SelectAction(s *Server, actionFlag int) error {
	// is metacontrol the best name?
	switch actionFlag {
	case 1:
		fmt.Println("here for the switch case to work %d - is the wrong number", actionFlag)
		// StartConsole
	case 2:
		fmt.Println("here for the switch case to work %d - is the wrong number", actionFlag)
		// OpenConsole
	case 3:
		fmt.Println("here for the switch case to work %d - is the wrong number", actionFlag)
		// StopConsole
	case 4:
		fmt.Println("here for the switch case to work %d - is the wrong number", actionFlag)
		// CreateConsole
	case 5:
		mc.StartServer(s)
	case 6:
		mc.StopServer(s)
	case 7:
		mc.PauseServer(s)
	case 8:
		mc.ResumeServer(s)
	case 9:
		mc.CreateServer(s)
	case 10:
		mc.RestartServer(s)

	default:
		err := fmt.Errorf("Incorrect actionFlag %d provide to SelectAction")
		return err
	}

	err := GracefulExit()
	util.CheckError(err)
	return nil
}
