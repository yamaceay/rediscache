package main

import (
	"fmt"
	"main/modes"
)

func main() {
	settings := modes.ReadSettings()
	if err := modes.StartServer(settings); err != nil {
		panic(fmt.Errorf("server ended: %s", err))
	}
}
