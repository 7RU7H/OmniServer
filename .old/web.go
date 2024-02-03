package OmniServer

import (
        "net/http"
        "strings"
        "context"
        "os"
)


// Webserver 
// Specifics of http or https (tls) are seperated to each

// Web server functionality that inclusive to both http,https,etc is here

// BIG REWRITE:
        // Seperate TLS and nonTLS as functions that then are the below go routine,
        // For ServerID, serverCertPath, serverKeyPath
        go func(server *http.Server, ctx Context, cancelCtx CancelFunc) error {

                // ListenAndServer either HTTP or HTTPS server
                // HTTP
                // WILL NOT REQUIRE ListenAndServer() as Contexts will be used was just an idiot
                // go ListenAndServerWebServer()
                if !isTLS {
                // goroutine this function
                        // Better error handling - account for contexts, go routines, when it should return exit out of this block
                        err := server.ListenAndServe()
                        if errors.Is(err, http.ErrServerClosed) {
                                        fmt.Printf("%s closed\n", ServerID, err)
                                        log.Fatal("%s closed\n", ServerID, err)
                                        return err
                        } else if err != nil {
                                        fmt.Printf("Error listening for %s: %s\n", ServerID, err)
                                        log.Fatal("Error listening for %s: %s\n", ServerID, err)
                                        return err
                        } else {
                                        log.Printf("%s is listening...\n", ServerID)
                                        return err
                        }
                        cancelCtx()
                } else {
                        // If HTTPS server
                        //serverStartTime := time.Now()
                        err := http.ListenAndServeTLS(listeningPort, serverCertPath, serverKeyPath, nil)
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
                        cancelCtx()
                }

                return nil
        }()


// Upload file - filename
func UploadFileHandler(w http.ResponseWriter, r *http.Request) error {

        // Parse our multipart form, 10 << 20 specifies a maximum
        // upload of 10 MB files.
        //r.ParseMultipartForm(10 << 20)

        // Get filename from body of r.Body

        // FormFile returns the first file for the given key `myFile`
        // it also returns the FileHeader so we can get the Filename,
        // the Header and the size of the file

        log.Printf("/upload/%s - Upload requested by ...", handler.Filename)
        file, handler, err := r.FormFile()
        if err != nil {
                // Error retrieving file of filename
                return err
        }
        startTime := time.Now()
        defer file.Close()
        //log.Print("",  ) File upload request success
        //log.Print("",  ) File upload INFO:
        log.Printf("Uploaded File: %+v\n", handler.Filename)
        lof.Printf("File Size: %+v\n", handler.Size)
        log.Printf("MIME Header: %+v\n", handler.Header)
        //log.Print("",  ) File upload request success

        // Create a temporary file within our temp-images directory that conforms to a naming scheme
        tempFile, err := os.TempFile(tmpUploadDir, "tmp-")
        if err != nil {
                // Error creating temporary file
                return err
        }
        defer tempFile.Close()

        // read all of the contents of our uploaded file into a byte array
        fileBytes, err := os.ReadFile(file)
        if err != nil {
                // Failed to read file being uploaded to byte array
                return err
        }
        // write this byte array to our temporary file
        tempFile.Write(fileBytes)
        endTime := time.Now()
        //Return that we have successfully uploaded our file!
        log.Printf("Successfully Uploaded File - %s \n", handler.Filename)
        fmt.Fprintf(w, "Successfully Uploaded File - %s \n", handler.Filename)
        return nil
}

// Download file - filename
func DownloadFileHandler(w http.ResponseWriter, r *http.Request) error {

        // client := Headers - IP User-Agent
        // requestedFileToDownload :=

        if !util.checkFileExists(requestedFileToDownload) {
                w.WriteHeader(404)
                w.Write([]byte("404\n"))
                log.Fatal("Failed to Download file: %s - requested by: %s, %s", requestedFiletoDownload, clientIP, clientUA)
                // Fail to download file error
                return err
        } else {
                startTime := time.Now()
                log.Printf("Downloading file at:  %s - requested by: %s, %s", requestedFiletoDownload, clientIP, clientUA)
        }
        endTime := time.Now()
        log.Printf("Successfully Download File - %s by %s\n", handler.Filename, clientIP, clientUA)
        return nil
}

func SaveReqBodyFileHandler(w http.ResponseWriter, r *http.Request) error {
        builder := strings.Builder()
        startTime := time.Now()
        builder.WriteString(os.TempDir() + "/" + strings.ReplaceAll(r.RemoteAddr, ".", "-") + "-T-" + strconv.Itoa(int(time.Now().Unix())))
        filepath :=     builder.Write()
        err := os.Create(filepath,  0644)
        if err != nil {
                log.Fatal(err)
                return err
        }
        err := io.Copy(filepath, r.Body)
        if err != nil {
                log.Fatal(err)
                return err
        }
        defer f.Close()
        endTime := time.Now()
        // Log file and time
        builder.Flush()
        return nil
}
