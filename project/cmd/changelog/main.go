package main

import (
	"fmt"

	"github.com/klyed/hivesmartchain/project"
)

func main() {
	fmt.Println(project.History.MustChangelog())
}
