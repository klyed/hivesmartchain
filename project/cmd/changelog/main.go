package main

import (
	"fmt"

	"github.com/KLYE-Dev/HSC-MAIN/project"
)

func main() {
	fmt.Println(project.History.MustChangelog())
}
