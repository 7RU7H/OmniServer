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
	httpRegex := regexp.MustCompile(`#http#`)
	httpsRegex := regexp.MustCompile(`#https#`)
	interfaceRegex := regexp.MustCompile(`##`)
	ipRegex := regexp.MustCompile(`#\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}#`)
	portRegex := regexp.MustCompile(`#\d{1,5}#`)
	tlsRegex := regexp.MustCompile(`##`)
	sortedArgs := make([]string, len(args)+1)
	if len(args) != 6 {
		matchHTTP, err := regexp.MatchString(httpRegex.String(), regexSafeArgs)
		checkError(err)
		matchInterface, err := regexp.MatchString(interfaceRegex.String(), regexSafeArgs)
		checkError(err)
		matchIP, err := regexp.MatchString(ipRegex.String(), regexSafeArgs)
		checkError(err)
		matchPort, err := regexp.MatchString(portRegex.String(), regexSafeArgs)
		checkError(err)
		httpAllMatched := matchHTTP && matchInterface && matchIP && matchPort
		if !httpAllMatched {
			err := fmt.Errorf("Arguments provided are %v: %v", httpAllMatched, args)
			return nil, err
		}
		sortedArgs[0] = strings.ReplaceAll(strings.Join(httpRegex.FindAllString(regexSafeArgs, 1), ""), "#", "")
		sortedArgs[1] = strings.ReplaceAll(strings.Join(interfaceRegex.FindAllString(regexSafeArgs, 1), ""), "#", "")
		//sortedArgs[2] = convIfconfigNameToCIDR() // needs inferface from sortedArgs[1]
		sortedArgs[3] = strings.ReplaceAll(strings.Join(ipRegex.FindAllString(regexSafeArgs, 1), ""), "#", "")
		sortedArgs[4] = strings.ReplaceAll(strings.Join(portRegex.FindAllString(regexSafeArgs, 1), ""), "#", "")

		if !(checkValidIP(sortedArgs[3]) && checkValidPort(sortedArgs[4])) {
			err := fmt.Errorf("Invalid IP and Port combination: %s and %s", sortedArgs[3], sortedArgs[4])
			return nil, err
		}
		checkError(err)
		return sortedArgs, nil

	} else {
		matchHTTPS, err := regexp.MatchString(httpsRegex.String(), regexSafeArgs)
		checkError(err)
		matchInterface, err := regexp.MatchString(interfaceRegex.String(), regexSafeArgs)
		checkError(err)
		matchIP, err := regexp.MatchString(ipRegex.String(), regexSafeArgs)
		checkError(err)
		matchPort, err := regexp.MatchString(portRegex.String(), regexSafeArgs)
		checkError(err)
		matchTLS, err := regexp.MatchString(tlsRegex.String(), regexSafeArgs)
		checkError(err)
		httpsAllMatched := matchHTTPS && matchInterface && matchIP && matchPort && matchTLS
		if !httpsAllMatched {
			err := fmt.Errorf("Arguments provided are %v: %v", httpsAllMatched, args)
			return nil, err
		}
		// Get interfaceCIDR
		// Validate TLS
		sortedArgs[0] = strings.ReplaceAll(strings.Join(httpRegex.FindAllString(regexSafeArgs, 1), ""), "#", "")
		sortedArgs[1] = strings.ReplaceAll(strings.Join(interfaceRegex.FindAllString(regexSafeArgs, 1), ""), "#", "")
		//sortedArgs[2] = convIfconfigNameToCIDR() // needs inferface from sortedArgs[1]
		sortedArgs[3] = strings.ReplaceAll(strings.Join(ipRegex.FindAllString(regexSafeArgs, 1), ""), "#", "")
		sortedArgs[4] = strings.ReplaceAll(strings.Join(portRegex.FindAllString(regexSafeArgs, 1), ""), "#", "")
		// sortedArgs[5] = TLS
		if !(checkValidIP(sortedArgs[3]) && checkValidPort(sortedArgs[4])) {
			err := fmt.Errorf("Invalid IP and Port combination: %s and %s", sortedArgs[3], sortedArgs[4])
			return nil, err
		}
		checkError(err)
		return sortedArgs, nil
	}
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

func prependColonToPortNumber(port string) string {
	builder := strings.Builder{}
	builder.WriteString(":" + port)
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

func checkValidPort(portStr string) bool {
	portInt, err := strconv.Atoi(strings.ReplaceAll(portStr, ":", ""))
	checkError(err)
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
	checkError(err)
	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok {
			return ipNet.String(), nil
		}
	}
	return "", fmt.Errorf("No suitable IP address found")
}
