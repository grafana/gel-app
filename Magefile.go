//+build mage

package main

import (

	// mage:import
	_ "github.com/grafana/grafana-plugin-sdk-go/build"

	"github.com/magefile/mage/sh"
)

// CheckCheck builds production back-end components.
func CopyArtifacts() error {
	if err := sh.RunV("mkdir", "-p", "dist"); err != nil {
		return err
	}

	if err := sh.RunV("cp", "./README.md", "./dist"); err != nil {
		return err
	}

	if err := sh.RunV("cp", "./src/plugin.json", "./dist"); err != nil {
		return err
	}

	if err := sh.RunV("cp", "./src/pipeline.svg", "./dist"); err != nil {
		return err
	}
	return nil
}

// Default configures the default target.
var Default = CopyArtifacts
