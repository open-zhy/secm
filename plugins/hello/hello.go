package main

import (
	"fmt"
)

// Init is exported and will be called when the plugin is loaded
func Init() {
	fmt.Printf("plugins:hello > Hello, World!\n")
}
