package main

import (
	"github.com/jaytaylor/goose-cli/cmd"
)

func main() {
	//if len(os.Args) < 2 {
	//	fmt.Fprintf(os.Stderr, "missing required parameter: url\n")
	//	os.Exit(1)
	//}

	cmd.Execute()
}
