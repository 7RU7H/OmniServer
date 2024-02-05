package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// CLI -> if http else https -> Done - just simple done project - below is just a map of functions - see TODO idiot
// main -> handleArgs -> main
// switch on sortedArray -> subchecks on sortedArray size
// Either: http or https server

// TODO List TODO
// HTTP server
// Error ids and code - id for where in the source for no lost in the src and code for switch case fatal or not
// TODO AFTER THE ABOVE IS COMPLETE built and works no excuses:
// Application end - start time printTotalRuntime()
// go routine and channels
// TODO AFTER THE ABOVE IS COMPLETE built and works no excuses:
// TLS - regex requred that make sense, validationTLS(), how validateTLS passes data to buildHTTPS()
// TODO AFTER THE ABOVE IS COMPLETE built and works no excuses:
// runHTTPSserver()
// TODO AFTER THE ABOVE IS COMPLETE built and works no excuses:
// Add all the profession stuff
// - Make authentication actually work without any other packages
// - Colourful Banner!
// - DONOT WORRY ABOUT nested regex -> string sorted args oneliners no (5||6)*2 additional variable declarations making that underreadable dense vertically and save some memory

// Download file - filename
func DownloadFileHandler(w http.ResponseWriter, r *http.Request) error {
	// client := Headers - IP User-Agent

	requestedURL := ""
	requestedFile := ""
	clientIP := ""
	clientMAC := ""
	clientUA := ""
	exists, err := checkFileExists(requestedFile)
	if !exists {
		http.NotFound(w, r)
		log.Println("Failed to Download file: %v from %v - requested by: %v, %v, %v", requestedFile, requestedURL, clientIP, clientMAC, clientUA)
		return err
	} else {
		startTime := time.Now()
		http.ServeFile(w, r, requestedFile)
		log.Printf("Downloading file at: %v from: %v - requested by: %v, %v, %v\n", requestedFile, requestedURL, clientIP, clientMAC, clientUA)
	}
	endTime := time.Now()
	log.Printf("Successfully Downloaded File - %v by %v - %v\n", requestedFile, clientIP, clientMAC)
	return nil
}

// TODO
func UploadFileHandler(w http.ResponseWriter, r *http.Request) error {

	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	//r.ParseMultipartForm(10 << 20)

	// Get filename from body of r.Body

	// FormFile returns the first file for the given key `myFile`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file

	log.Printf("/upload/%s - Upload requested by ...", filename)
	file, handler, err := r.FormFile()
	if err != nil {
		// Error retrieving file of filename
		return err
	}
	startTime := time.Now()
	defer file.Close()
	//log.Print("",  ) File upload request success
	//log.Print("",  ) File upload INFO:
	log.Printf("Uploaded File: %+v\n", filename)
	lof.Printf("File Size: %+v\n", fileSize)
	log.Printf("MIME Header: %+v\n", mimeHeader)

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
	log.Printf("Successfully Uploaded File - %s \n", filename)
	fmt.Fprintf(w, "Successfully Uploaded File - %s \n", filename)
	return nil
}

func saveReqBodyFileHandler(r *http.Request) error {
	builder := strings.Builder()
	startTime := time.Now()
	builder.WriteString(os.TempDir())
	builder.WriteString("/")
	builder.WriteString(strings.ReplaceAll(r.RemoteAddr, ".", "-"))
	builder.WriteString("-T-")
	builder.WriteString(strconv.Itoa(int(time.Now().Unix())))
	filepath := builder.String()
	err := os.Create(filepath, 0644)
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
	builder.Flush()
	endTime := time.Now()
	log.Println("Entire process of file creation for file upload: %v - took from: %v till %v", filepath, startTime, endTime)
	return nil
}

func createDefaultWebServerMux() *ServerMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/upload", uploadFileHandler())
	mux.HandleFunc("/download", downloadFileHandler())
	mux.HandleFunc("/saveReqBody", saveReqBodyFileHandler())
	return mux
}

func initServerContext(lportString, keyServerAddr string, srvMux *http.ServerMux) (*http.Server, context.Context, context.CancelFunc, error) {
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

func buildHTTPServer(args []string) (*http.Server, context.Context, context.CancelFunc, error) {
	log.Println("--- Building HTTP Server ---")
	mux := createDefaultWebServerMux()
	log.Println("Mux for %v created", args)
	httpServer, ctx, cancelCtx, err := initServerContext(args[4], args[3], mux)
	checkError(err, 0)
	log.Println("--- Server Built for %v created ---", httpServer)
	return httpServer, ctx, cancelCtx, nil
}
func runHTTPServer(*http.Server, context.Context, context.CancelFunc) error {
	fmt.Println("running HTTP")
	return nil
}
func validateTLS(args string) (string, error) {
	fmt.Println("validating TLS")
	return "TLS incoming", nil
}

func buildHTTPSServer(nontlsArgs []string, tls string) error {
	fmt.Println("building HTTPS")
	return nil
}
func runHTTPSServer() error {
	fmt.Println("running HTTPS")
	return nil
}

func gracefulExit() error {
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
func checkError(err error, errorCode int) {
	if err != nil {
		log.Fatal(err, " - Error code: ", errorCode)
	}
}

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
func printTotalRuntime(appStartTime time.Time) {
	appTerminateTime := time.Now()
	totalRuntime := 0
	// Do the mathematics idiot
	log.Printf("Application started: %v - Terminated: %v - Runtime: %v\n", appStartTime, appTerminateTime, totalRuntime)
}

func main() {
	var ipAddress, serverType, tlsInputStr, netInterface string
	var portInt int
	flag.StringVar(&tlsInputStr, "t", "None", "Provide TLS information delimited by a comma - requires server type: -s https")
	flag.StringVar(&serverType, "s", "http", "Provide a server of type: http, https")
	flag.StringVar(&netInterface, "e", "lo", "Provide a Network Interface - required!")
	flag.StringVar(&ipAddress, "i", "127.0.0.1", "Provide a valid IPv4 address - required!")
	flag.IntVar(&portInt, "p", 8443, "Provide a TCP port number - required!")
	flag.Parse()
	// Banner !!

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
		httpServer, ctx, cancelCtx, err := buildHTTPServer(sortedArgs)
		checkError(err, 0)
		err = runHTTPServer(httpServer, ctx, cancelCtx)
		checkError(err, 0)
		err = gracefulExit()
		checkError(err, 0)
		break
	case "https":
		tlsReq, err := validateTLS(sortedArgs[5])
		checkError(err, 0)
		err = buildHTTPSServer(sortedArgs[0:4], tlsReq)
		checkError(err, 0)
		err = runHTTPSServer()
		checkError(err, 0)
		err = gracefulExit()
		checkError(err, 0)
		break

	default:
		err := fmt.Errorf("Invalid Sorted arguments and proof at index 0: %s", sortedArgs[0])
		checkError(err, 0)
	}

	printTotalRuntime(appStartTime)
	return
}
