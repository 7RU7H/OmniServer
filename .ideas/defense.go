package 

import (
	"runtime"
)


switch runtime.GOOS {
case "android":
		 defendAndroid()
case "darwin":
		 defendDarwin()
case "dragonfly":
		 defendDragonfly()
case "freebsd":
		 defendFreebsd()
case "linux":
		 defendLinux()
case "nacl":
		 defendNacl()
case "netbsd":
		 defendNetbsd()
case "openbsd":
		 defendOpenbsd()
case "plan9":
		 defendPlan9()
case "solaris":
		 defendSolaris()
case "windows":
		 defendWindows()
default:
	err := fmt.Errorf("Unknown runtime.GOOS: %s", runtime.GOOS)
}

// [cmd : flags]string
// cmd := exec.Command([0], [1])
// err := cmd.Start()

// Filesystem, Process, Kernel, Memory, PortSwitching

// chattr *nix;  windows: attrib +h