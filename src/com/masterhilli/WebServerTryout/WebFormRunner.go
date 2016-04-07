package main

import (
	"./Webserver"
	"os/exec"
	"os"
)

func main() {
	go runMongoDB()
	Webserver.RunWebServer("..\\..\\..\\..")
}

func runMongoDB() {
	cmd := exec.Command(`D:\IDE\mongoDB\32\bin\mongod.exe`, "--dbpath", `D:\IDE\mongoDB\data`	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		panic(err)
	}
	cmd.Wait()
}
