//+build mage

package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/magefile/mage/mg"

	// mage:import
	"github.com/grafana/grafana-plugin-sdk-go/build"
	"github.com/grafana/grafana-plugin-sdk-go/build/utils"
)

const distDir = "dist"
const pluginJsonFileName = "plugin.json"

func CopyArtifacts() error {
	exists, err := utils.Exists(distDir)
	if err != nil {
		return err
	}
	if !exists {
		if err := os.MkdirAll(distDir, 0775); err != nil {
			return err
		}
	}

	fpaths := []string{path.Join("src", pluginJsonFileName), "README.md"}
	fis, err := ioutil.ReadDir("img")
	if err != nil {
		return err
	}
	for _, fi := range fis {
		fpaths = append(fpaths, path.Join("img", fi.Name()))
	}

	for _, fpath := range fpaths {
		if err := utils.CopyFile(fpath, path.Join(distDir, path.Base(fpath))); err != nil {
			return err
		}
	}

	return nil
}

// PluginJSON is a minimal type representation of the plugin.json file found in src/.
type PluginJSON struct {
	ID   string     `json:"id"`
	Info PluginInfo `json:"info"`
}

// PluginInfo is the Info member type for PluginJSON.
type PluginInfo struct {
	Version string `json:"version"`
}

func readPluginJSON() (PluginJSON, error) {
	var pj PluginJSON
	byteValue, err := ioutil.ReadFile(path.Join(distDir, pluginJsonFileName))
	if err != nil {
		return pj, err
	}

	err = json.Unmarshal([]byte(byteValue), &pj)
	return pj, err
}

// Zip builds the plugin zip archive.
func Zip() error {
	mg.Deps(CopyArtifacts, build.BuildAll)

	pluginJson, err := readPluginJSON()
	if err != nil {
		return err
	}

	if pluginJson.ID == "" {
		return fmt.Errorf("can not create zip file because the id property in %v is missing and is needed for the zip file name", pluginJsonFileName)
	}
	if pluginJson.Info.Version == "" {
		return fmt.Errorf("can not create zip file because the info.version property in %v is missing and is needed for the zip file name", pluginJsonFileName)
	}
	pluginName := fmt.Sprintf("%s-%s", pluginJson.ID, pluginJson.Info.Version)
	return makeZip(pluginName)
}

func makeZip(pluginName string) error {
	zipFname := fmt.Sprintf("%s.zip", pluginName)
	log.Printf("Creating zip archive %q", zipFname)
	f, err := os.OpenFile(zipFname, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	zw := zip.NewWriter(f)
	defer zw.Close()

	if err := filepath.Walk(distDir, func(srcPath string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("filepath.Walk failed for %q: %s", srcPath, err)
			return err
		}

		fi, err := os.Lstat(srcPath)
		if err != nil {
			return err
		}
		fh, err := zip.FileInfoHeader(fi)
		if err != nil {
			return err
		}
		fh.Name = path.Join(pluginName, strings.TrimPrefix(srcPath, distDir))
		log.Printf("Adding %q to zip archive as %q", srcPath, fh.Name)

		if info.IsDir() {
			fh.Name += "/"
		}
		fw, err := zw.CreateHeader(fh)
		if err != nil {
			return err
		}
		if info.IsDir() {
			log.Printf("Adding directory to zip archive: %q", fh.Name)
			return nil
		}

		src, err := os.Open(srcPath)
		if err != nil {
			return err
		}
		defer src.Close()
		_, err = io.Copy(fw, src)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

// Default configures the default target.
var Default = Zip
