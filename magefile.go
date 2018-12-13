// +build mage

package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
	"github.com/magefile/mage/sh"
)

const (
	appBin = "nicoime.exe"
)

// Default target to run when none is specified
// If not set, running mage will list available targets
var Default = Run

// Manage your deps, or running package managers.
func Setup() error {
	deps := []string{
		"github.com/creasty/go-nkf",
		"github.com/PuerkitoBio/goquery",
	}
	fmt.Println("Installing Deps...")
	for _, dep := range deps {
		cmd := exec.Command("go", "get", "-u", dep)
		err := cmd.Run()
		if err != nil {
			return err
		}
	}
	return nil
}

// Clean clean up after yourself
func Clean() {
	fmt.Println("Clean...")
	os.Remove(appBin)
}

// Build build step that requires additional params, or platform specific steps for example
func Build() error {
	mg.Deps(Clean)
	fmt.Println("Build...")
	return sh.Run("go", "build", "-o", appBin, "-v")

}

// Run execute app
func Run() error {
	mg.Deps(Build)
	fmt.Println("Run...")
	return sh.Run("./" + appBin)
}
