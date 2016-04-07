package main

import (
	"./Webserver"
	"os/exec"
	"flag"
	. "./Webserver/Logger"
	"os"
)

func main() {
	args := flag.Args()
	pathToResources := "."
	if (len(args) != 0) {
		pathToResources = args[0]
	}
	go runMongoDB()
	Webserver.RunWebServer(pathToResources)
}

func runMongoDB() {
	cmd := exec.Command(`D:\IDE\mongoDB\32\bin\mongod.exe`, "--dbpath", `D:\IDE\mongoDB\data`	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	err := cmd.Start()
	if err != nil {
		LOGGER.Printf("Had to stop the server because: %s")
		return
	}
	cmd.Wait()
}
