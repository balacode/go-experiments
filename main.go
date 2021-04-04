// -----------------------------------------------------------------------------
// Go Language Experiments                              go-experiments/[main.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package main

import (
	"fmt"
	"strings"
)

var div = strings.Repeat("-", 80)

func main() {
	fmt.Println(div)
	fmt.Println("Running go-experiments...")
	{
		chacha20EncryptionDemo()
		// serverDemo()
		// tlsServerDemo()
	}
	fmt.Println(div)
	fmt.Println("Finished go-experiments")
} //                                                                        main

// end
