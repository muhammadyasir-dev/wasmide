package runnerservice

import (
	"os/exec"
)

func Execwasm(programminglanguage string) {

	switch programminglanguage {
	case "rust":
		exec.Command("bash", "-c", "cargo", "run")

	case "go":
		exec.Command("bash", "-c", "go", "run", "main.go")

	case "c":
		exec.Command("bash", "-c", "make", "main")
	case "c++":
		exec.Command("bash", "-c", "make", "main")

	}
}
