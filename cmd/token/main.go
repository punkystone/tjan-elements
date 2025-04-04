//nolint:forbidigo //allow
package main

import (
	"fmt"
	"os"
	"tjan-elements/internal/sheets"
)

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println("no code provided")
		return
	}
	authCode := args[0]
	err := sheets.CreateToken(authCode)
	if err != nil {
		fmt.Println("error creating token:", err)
		return
	}
	fmt.Println("token saved")
}
