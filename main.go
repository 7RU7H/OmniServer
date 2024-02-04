package OmniServer

// +build !nolog

import (
        "fmt"
        "log"
        "net"
        "os"
        "strconv"
        "strings"
)

// CLI -> if http else https -> Done - just simple done project - below is just a map of functions - see TODO idiot
// main -> handleArgs -> main 
// switch on sortedArray -> subchecks on sortedArray size
// Either: http or https server

// TODO List TODO
// How is VHOSTing or DNS playing a role? Hosting Configuration;
// Versioning string for Version Flag
// WTF is extra variables for - does the flag var store the flag value or input value
// http or tls command - can you reuse &flagvariable per command? then it a switch on command
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

// Validates and sorts Args to serverType, interfaceName, interfaceCIDR(retrived by this application), IP, Port, (optional Hostname)), TLS
func handleArgs(args []string) error {
        const correctTlsArgsCount int = 3
        var ipAddress, netInterface, serverType, tlsInputStr, hostnamesStr, helpFlag, versionFlag string
        var portInt int
        flag.StringVar(&tlsInputStr, "-t", "None", "Provide TLS information delimited by a comma - requires server type: -s https")
        flag.StringVar(&serverType, "-s", "http", "Provide a server of type: http, https")
        flag.StringVar(&netInterface, "-e", "localhost", "Provide a Network Interface - required!")
        flag.StringVar(&ipAddress, "-i", "127.0.0.1", "Provide a valid IPv4 address - required!")
        flag.StringVar(&hostnamesStr, "-H", "", "Provide DNS or vhosting alias comma delimited: dc.test.org,test.org")
        flag.IntVar(&portInt, "-p", 8443, "Provide a TCP port number - required!")

        flag.StringVar(&helpFlag, "-h", "Help", "Help")
        flag.StringVar(&versionFlag, "-v", "Version", "Version")
        flag.Parse()
        argsLen := len(args)

        if argsLen > 4 {
                flag.Usage()
                os.Exit(1)
        }
     

	// Does this actually make sense more variable does not make sense
        var serverTypeArg, netInterfaceName, netInterfaceCDIR, ipAddressArg, tlsInputStr, hostnameArgStr string
        var portIntArg int
// These checks should check validity 
                flag.Parse(args[1:])
                if flag = flag.Lookup(serverType); flag == nil || serverType != "http" && serverType != "https" {
                        flag.Usage()
                        err := fmt.Errorf("Invalid ServerType flag -s <server type> passed by user")
                        return err
                } else {
                        serverTypeArg = flag.Value.(flag.Getter).Get().(string)
                }
                if flag = flag.Lookup(netInterface); flag == nil || netInterface == "" {
                        flag.Usage()
                        err := fmt.Errorf("Missing Network Interface -e <network interface>")
                        return err
                } else {
                        netInterfaceName = flag.Value.(flag.Getter).Get().(string)
                        ifconfig, err := net.InterfaceByName(netInterfaceArg)
                        util.CheckError(err)
                        ifconfigCIDRTmp, err := util.ConvIfconfigNameToCIDR(ifconfig, netInterfaceName)
                        util.CheckError(err)
                        netInterfaceCDIR = ifconfigCIDRTmp
                }
                if flag = flag.Lookup(portInt); flag == nil {
                        flag.Usage()
                        err := fmt.Errorf("Invalid TCP port number -p <port number>")
                        return err
                } else {
                        portIntArg = flag.Value.(flag.Getter).Get().(int)
                }
                if flag = flag.Lookup(hostnamesStr); flag == nil || hostnamesStr == "" {
                        flag.Usage()
                        err := fmt.Errorf("Missing DNS or vhosting alias -H <alias>")
                        return err
                } else {
                        hostnamesStrArg = flag.Value.(flag.Getter).Get().(string)
                }
                if flag = flag.Lookup(ipAddress); flag == nil || ipAddress == "" {
                        flag.Usage()
                        err := fmt.Errorf("Missing valid IPv4 address -i <ip address>")
                        return err
                } else {
                        ipAddressArg = flag.Value.(flag.Getter).Get().(string)
                        if !util.CheckValidIP(ipAddressArg) {
                                flag.Usage()
                                err := fmt.Errorf("Invalid IP address provided: %s".ipAddressArg)
                                return err
                        }
                }
                if flag = flag.Lookup(tlsInputStr); flag != nil {
                        tlsInputStr = flag.Value.(flag.Getter).Get().(string)
                        if flag = flag.Lookup(serverType); flag == nil && flag.Value.(flag.Getter).Get.(string) != "https" {
                                flag.Usage()
                                err := fmt.Errorf("TLS information provided, but server type -s is not https")
                                return err
                        }
                        if strings.Count(tlsInputStr, ",") != correctTlsArgsCount {
                                flag.Usage()
                                err := fmt.Errorf("TLS information provided does not contain correct number of comma delimited arguments: %s", tlsInputStr)
                                return err
                        }
                        // TLS library checks can occur here!
                } else {
                        tlsStrArg = ""
                }
        case "console":
                consoleFlag = 1
        default:
                flag.Usage()
                err := fmt.Errorf("Bad command given")
                return err
        }

        CheckError(err)
       
        switch consoleFlag {
        case 1:
                server := Server{}
                server.setDefaultServerConfig()
                server.InitServerFromArgs(mc.serverCounter, serverTypeArg, ipAddressArg, hostnameArgStr, netInterfaceName, netInterfaceCDIR, portIntArg, tlsInputStr)
                mc.SelectAction(server, consoleFlag)
        case 2:
                fmt.Println("metahandler.ToConsole()")
        case 0:
                fmt.Println("Error with checkArgs")
        default:
                err := fmt.Errorf("No consoleFlag: %d value passed", consoleFlag)
                return err
        }

        return nil
}

func main() {
       	sortedArgs, err := handleArgs(os.Args[1:])
	checkError(err)
	
	// len +1 includes hostname checks
	switch sortedArgs[0] {
	case "http":
	case "tls":
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

func convIfconfigNameToCIDR(ifconfig *Interface) (string, error) {
        addr, err := ifconfig.Addrs()
        CheckError(err)
        for _, addr := range addrs {
                if ipNet, ok := addr.(*net.IPNet); ok {
                        return ipNet.String(), nil
                }
        }
        return "", fmt.Errorf("No suitable IP address found")
}

