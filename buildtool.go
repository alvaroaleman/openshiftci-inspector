// +build ignore

package main

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

func usage() {
	println("Usage: go run buildtool.go [COMMAND]")
	println("")
	println("Possible commands are:")
	println("  api: to generate the API client for the frontend")
	println("  export-api: to generate the swagger.json file from the backend only")
	println("  import-api: to generate the frontend library from swagger.json only")
	println("  frontend: to generate the frontend assets for embedding")
}

func runInFrontend(args ...string) {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = "frontend"
	if err := cmd.Run(); err != nil {
		log.Fatalf("error while running %s: %v", strings.Join(args, " "), err)
	}
}

func run(args ...string) {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("error while running %s: %v", strings.Join(args, " "), err)
	}
}

func exportAPI() {
	run("goswagger", "generate", "spec", "-o", "swagger.json")
}

func importAPI() {
	runInFrontend("npm", "install")
	runInFrontend("npm", "run", "client")
}

func api() {
	exportAPI()
	importAPI()
}

func frontend() {
	runInFrontend("npm", "install")
	runInFrontend("npm", "run", "build")
}

func main() {
	if len(os.Args) != 2 {
		usage()
		os.Exit(1)
	}
	switch os.Args[1] {
	case "-h":
		usage()
	case "--help":
		usage()
	case "api":
		api()
	case "export-api":
		exportAPI()
	case "import-api":
		importAPI()
	case "frontend":
		frontend()
	default:
		usage()
		os.Exit(1)
	}
}
