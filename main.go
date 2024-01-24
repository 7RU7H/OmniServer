package OmniServer

import (
	"os"

	"github.com/7ru7h/OmniServer/cli.go"
)

func main() {
	os.Exit(cli.HandleAll(os.Args[1:]))
}
