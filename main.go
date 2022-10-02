package main

import (
	"fmt"

	"github.com/irfansofyana/tetengger/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil && err.Error() != "" {
		fmt.Println(err)
	}
}
