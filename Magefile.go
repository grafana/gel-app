//+build mage

package main

import (
	"io/ioutil"
	"os"
	"path"

	// mage:import
	"github.com/grafana/grafana-plugin-sdk-go/build"
	"github.com/grafana/grafana-plugin-sdk-go/build/utils"
)

func CopyArtifacts() error {
	if err := os.MkdirAll("dist", 0775); err != nil {
		return err
	}

	fpaths := []string{"src/plugin.json", "README.md"}
	fis, err := ioutil.ReadDir("img")
	if err != nil {
		return err
	}
	for _, fi := range fis {
		fpaths = append(fpaths, path.Join("img", fi.Name()))
	}

	for _, fpath := range fpaths {
		if err := utils.CopyFile(fpath, path.Join("dist", path.Base(fpath))); err != nil {
			return err
		}
	}

	return nil
}

// Default configures the default target.
var Default = build.BuildAll
