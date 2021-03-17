package main

import (
	"fmt"

	"github.com/klye-dev/hsc-main/project"
)

func main() {
	fmt.Println(project.History.CurrentVersion().String())
}
