//+build mage

package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"github.com/magefile/mage/mg"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	// mage:import
	"github.com/grafana/grafana-plugin-sdk-go/build"
	"github.com/grafana/grafana-plugin-sdk-go/build/utils"
)

func CopyArtifacts() error {
	// TODO: Only try to remove if dist exists, check error
	_ = os.RemoveAll("dist")
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

func readPluginJSON() (map[string]interface{}, error) {
	byteValue, err := ioutil.ReadFile(path.Join("dist", "plugin.json"))
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal([]byte(byteValue), &result)
	return result, err
}

// Zip builds the plugin zip archive.
func Zip() error {
	mg.Deps(CopyArtifacts, build.BuildAll)

	pluginJson, err := readPluginJSON()
	if err != nil {
		return err
	}

	pluginName := fmt.Sprintf("%s-%s", pluginJson["id"], pluginJson["version"])
	zipFname := fmt.Sprintf("%s.zip", pluginName)
	log.Printf("Creating zip archive %q", zipFname)
	f, err := os.OpenFile(zipFname, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	zw := zip.NewWriter(f)
	defer zw.Close()

	if err := filepath.Walk("dist", func(srcPath string, info os.FileInfo, err error) error {
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
		fh.Name = path.Join(pluginName, strings.TrimPrefix(srcPath, "dist"))
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
