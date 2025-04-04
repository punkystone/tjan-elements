//nolint:forbidigo //allow
package main

import (
	"fmt"
	"tjan-elements/internal/sheets"
)

func main() {
	authURL, err := sheets.CreateAuthURL()
	if err != nil {
		fmt.Println("Error creating auth URL:", err)
		return
	}
	fmt.Println(authURL)
}
