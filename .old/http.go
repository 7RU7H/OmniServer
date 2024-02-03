package OmniServer

import (
	"context"
	"net"
	"net/http"

	"github.com/7ru7h/OmniServer/web.go"
)

// http specifics go here

func CreateWebServer(s *Server) error {
	// Define Mux first to then pass it in Context Creation
	s.Mux = CreateDefaultWebServerMux()
	// Context creation
	// Assigned to a struct!
	s.ServerWithCtx, s.Ctx, s.CancelCtx = InitServerContext(s.ServerInfo.ListeningPort, s.ServerInfo.ServerAddr, s.Mux)
	return nil
}

// Mux is a multiplexer to handle routes for Webserver
func CreateDefaultWebServerMux() *ServerMux {
	mux := http.NewServeMux()
	// Setup routes
	mux.HandleFunc("/upload", web.UploadFileHandler())
	mux.HandleFunc("/download", web.DownloadFileHandler())
	mux.HandleFunc("/saveReqBody", web.SaveReqBodyFileHandler())
	return mux
}

func InitServerContext(lportString, keyServerAddr string, srvMux *ServerMux) (*http.Server, Context, CancelFunc, error) {
	ctx, cancelCtx := context.WithCancel(context.Background())
	server := &http.Server{
		Addr:    lportString,
		Handler: srvMux,
		BaseContext: func(l net.Listener) context.Context {
			ctx = context.WithValue(ctx, keyServerAddr, l.Addr().String())
			return ctx
		},
	}
	return server, ctx, cancelCtx, nil
}
