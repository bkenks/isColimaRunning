package main

import (
	"os"
	"os/exec"
)

const colimaPath = "/opt/homebrew/bin" // wherever `which colima` points

func colimaCommand(args ...string) *exec.Cmd {
	cmd := exec.Command(colimaPath+"/colima", args...)
	cmd.Env = append(os.Environ(), "PATH="+colimaPath+":/usr/bin:/bin")
	return cmd
}

func isColimaRunning() bool {
	return colimaCommand("status").Run() == nil
}

func startColima() error { return colimaCommand("start").Run() }
func stopColima() error  { return colimaCommand("stop").Run() }
