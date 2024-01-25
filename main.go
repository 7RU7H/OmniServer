package OmniServer

import (
	"os"
)

func main() {
	os.Exit(cli.HandleAll(os.Args[1:]))
}
