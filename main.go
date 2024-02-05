package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// CLI -> if http else https -> Done - just simple done project - below is just a map of functions - see TODO idiot
// main -> handleArgs -> main
// switch on sortedArray -> subchecks on sortedArray size
// Either: http or https server

// TODO List TODO
// BUILD AND BUILD
// TODO AFTER THE ABOVE IS COMPLETE built and works no excuses:
// HTTP server
// TODO AFTER THE ABOVE IS COMPLETE built and works no excuses:
// go routine and channels
// TODO AFTER THE ABOVE IS COMPLETE built and works no excuses:
// TLS - regex requred that make sense, validationTLS(), how validateTLS passes data to buildHTTPS()
// TODO AFTER THE ABOVE IS COMPLETE built and works no excuses:
// runHTTPSserver()
// TODO AFTER THE ABOVE IS COMPLETE built and works no excuses:
// Add all the profession stuff
// - DONOT WORRY ABOUT nested regex -> string sorted args oneliners no (5||6)*2 additional variable declarations making that underreadable dense vertically and save some memory

func buildHTTPServer(args []string) error {
	fmt.Println("building HTTP")
	return nil
}
func runHTTPServer() error {
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

// Validates and sorts Args to serverType, interfaceName, interfaceCIDR(retrived by this application), IP, Port, TLS
func handleArgs(args []string) ([]string, error) {
	regexSafeArgs := "#" + strings.Join(args, "#") + "#"
	httpRegex := regexp.MustCompile(`#http#`)
	httpsRegex := regexp.MustCompile(`#https#`)
	ipRegex := regexp.MustCompile(`#\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}#`)
	portRegex := regexp.MustCompile(`#\d{1,5}#`)
	tlsRegex := regexp.MustCompile(`##`)
	sortedArgs := make([]string, len(args)+1)

	matchIP, err := regexp.MatchString(ipRegex.String(), regexSafeArgs)
	checkError(err, 0)
	matchPort, err := regexp.MatchString(portRegex.String(), regexSafeArgs)
	checkError(err, 0)

	if len(args) != 6 {
		matchHTTP, err := regexp.MatchString(httpRegex.String(), regexSafeArgs)
		checkError(err, 0)
		tmpReduced := httpRegex.ReplaceAllString(regexSafeArgs, "")
		tmpReduced = ipRegex.ReplaceAllString(tmpReduced, "")
		tmpReduced = portRegex.ReplaceAllString(tmpReduced, "")
		ifconfigPotentialName := strings.ReplaceAll(tmpReduced, "#", "")
		ifconfig, err := net.InterfaceByName(ifconfigPotentialName)
		checkError(err, 0)
		ifconfigCIDRTmp, err := convIfconfigNameToCIDR(ifconfig)
		checkError(err, 0)
		httpAllMatched := matchHTTP && matchIP && matchPort
		if !httpAllMatched {
			err := fmt.Errorf("Arguments provided are %v: %v", httpAllMatched, args)
			return nil, err
		}
		sortedArgs[0] = strings.ReplaceAll(strings.Join(httpRegex.FindAllString(regexSafeArgs, 1), ""), "#", "")
		sortedArgs[1] = ifconfigPotentialName
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
		tmpReduced := httpRegex.ReplaceAllString(regexSafeArgs, "")
		tmpReduced = ipRegex.ReplaceAllString(tmpReduced, "")
		tmpReduced = portRegex.ReplaceAllString(tmpReduced, "")
		tmpReduced = tlsRegex.ReplaceAllString(tmpReduced, "")
		ifconfigPotentialName := strings.ReplaceAll(tmpReduced, "#", "")
		ifconfig, err := net.InterfaceByName(ifconfigPotentialName)
		checkError(err, 0)
		ifconfigCIDRTmp, err := convIfconfigNameToCIDR(ifconfig)
		checkError(err, 0)
		matchTLS, err := regexp.MatchString(tlsRegex.String(), regexSafeArgs)
		checkError(err, 0)
		httpsAllMatched := matchHTTPS && matchIP && matchPort && matchTLS
		if !httpsAllMatched {
			err := fmt.Errorf("Arguments provided are %v: %v", httpsAllMatched, args)
			return nil, err
		}
		sortedArgs[0] = strings.ReplaceAll(strings.Join(httpRegex.FindAllString(regexSafeArgs, 1), ""), "#", "")
		sortedArgs[1] = ifconfigPotentialName
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

func main() {
	var ipAddress, serverType, tlsInputStr string
	//var netInterface string
	var portInt int
	flag.StringVar(&tlsInputStr, "t", "None", "Provide TLS information delimited by a comma - requires server type: -s https")
	flag.StringVar(&serverType, "s", "http", "Provide a server of type: http, https")
	//flag.StringVar(&netInterface, "e", "lo:", "Provide a Network Interface - required!")
	flag.StringVar(&ipAddress, "i", "127.0.0.1", "Provide a valid IPv4 address - required!")
	flag.IntVar(&portInt, "p", 8443, "Provide a TCP port number - required!")
	flag.Parse()

	args, argsLen := flag.Args(), len(flag.Args())

	if argsLen > 4 {
		err := fmt.Errorf("The number arguments provided was %d", argsLen)
		checkError(err, 0)
		os.Exit(1)
	}

	sortedArgs, err := handleArgs(args)
	checkError(err, 0)

	// len +1 includes hostname checks
	switch sortedArgs[0] {
	case "http":
		err := buildHTTPServer(sortedArgs)
		checkError(err, 0)
		err = runHTTPServer()
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
	return
}
