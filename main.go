package OmniServer

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
// Refactor out structs use a sorted array
// Best CheckError function for all the golangs
// TODO AFTER THE ABOVE IS COMPLETE built and works no excuses:
// HTTP server
// TODO AFTER THE ABOVE IS COMPLETE built and works no excuses:
// go routine and channels
// TODO AFTER THE ABOVE IS COMPLETE built and works no excuses:
// TLS
// TODO AFTER THE ABOVE IS COMPLETE built and works no excuses:
// HTTPS server
// TODO AFTER THE ABOVE IS COMPLETE built and works no excuses:
// Add all the profession

// Validates and sorts Args to serverType, interfaceName, interfaceCIDR(retrived by this application), IP, Port, TLS
func handleArgs(args []string) ([]string, error) {
	regexSafeArgs := "#" + strings.Join(args, "#") + "#"
	httpRegex := `#http#`
	httpsRegex := `#https#`
	interfaceRegex := ``
	ipRegex := `#\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}#`
	portRegex := `#\d{1,5}#`
	tlsRegex := ``

	if len(args) != 6 {
		matchHTTP, err := regexp.MatchString(httpRegex, regexSafeArgs)
		checkError(err)
		matchInterface, err := regexp.MatchString(interfaceRegex, regexSafeArgs)
		checkError(err)
		matchIP, err := regexp.MatchString(ipRegex, regexSafeArgs)
		checkError(err)
		matchPort, err := regexp.MatchString(portRegex, regexSafeArgs)
		checkError(err)
		httpAllMatched := matchHTTP && matchInterface && matchIP && matchPort
		if !httpAllMatched {

		}
		// Get interfaceCIDR

	} else {
		matchHTTPS, err := regexp.MatchString(httpsRegex, regexSafeArgs)
		checkError(err)
		matchInterface, err := regexp.MatchString(interfaceRegex, regexSafeArgs)
		checkError(err)
		matchIP, err := regexp.MatchString(ipRegex, regexSafeArgs)
		checkError(err)
		matchPort, err := regexp.MatchString(portRegex, regexSafeArgs)
		checkError(err)
		matchTLS, err := regexp.MatchString(tlsRegex, regexSafeArgs)
		checkError(err)
		httpsAllMatched := matchHTTPS && matchInterface && matchIP && matchPort && matchTLS
		if !httpsAllMatched {

		}
		// Get interfaceCIDR
		// Validate TLS
	}

	return sortedArgs, nil
}

func main() {
	var ipAddress, netInterface, serverType, tlsInputStr string
	var portInt int
	flag.StringVar(&tlsInputStr, "t", "None", "Provide TLS information delimited by a comma - requires server type: -s https")
	flag.StringVar(&serverType, "s", "http", "Provide a server of type: http, https")
	flag.StringVar(&netInterface, "e", "localhost", "Provide a Network Interface - required!")
	flag.StringVar(&ipAddress, "i", "127.0.0.1", "Provide a valid IPv4 address - required!")
	flag.IntVar(&portInt, "p", 8443, "Provide a TCP port number - required!")
	flag.Parse()

	args, argsLen := flag.Args(), len(flag.Args())

	if argsLen > 4 {
		// Requires
		os.Exit(1)
	}

	sortedArgs, err := handleArgs(args)
	checkError(err)

	// len +1 includes hostname checks
	switch sortedArgs[0] {
	case "http":
	case "https":
	default:
	}

}

func convPortNumber(portNumber int) string {
	builder := strings.Builder{}
	portStr := strconv.Itoa(portNumber)
	builder.WriteString(":" + portStr)
	listeningPort := builder.String()
	builder.Reset()
	return listeningPort
}

// Build and improve
func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func checkFileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		log.Fatal(err)
		return false, err
	}
	if os.IsNotExist(err) {
		log.Fatal("File path does not exist")
		return false, err
	}
	return true, nil
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
	checkError(err)
	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok {
			return ipNet.String(), nil
		}
	}
	return "", fmt.Errorf("No suitable IP address found")
}
