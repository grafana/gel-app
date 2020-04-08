//+build mage

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"

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

func readPluginJSON() (map[string]interface{}, error) {
	byteValue, err := ioutil.ReadFile(path.Join("src", "plugin.json"))
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal([]byte(byteValue), &result)
	return result, err
}

func MakePluginZip() error {
	pluginJson, err := readPluginJSON()
	if err != nil {
		return err
	}

	zip := fmt.Sprintf("%s-%s.zip", pluginJson["id"], pluginJson["version"])

	// TODO: pick final name based on the contents of plugin.json
	if err := sh.RunV("zip", "-r", zip, "dist"); err != nil {
		return err
	}

	return nil
}

// Default configures the default target.
var Default = CopyArtifacts
