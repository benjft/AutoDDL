package main

import (
	"fmt"
)

func main() {
	config, _ := LoadConfigFromFiles("appsettings.yaml", "appsettings.development.yaml")
	fmt.Printf("%+v", config)
}
