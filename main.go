package OmniServer

import (
	"os"
	
	"github.com/7ru7h/OmniServer/metahandler.go"
)

func main() {
	os.Exit(metahandler.HandleAll(os.Args[1:]))
}
