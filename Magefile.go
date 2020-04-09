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

const distDir = "dist"
const pluginJson = "plugin.json"

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

	fpaths := []string{path.Join("src", pluginJson), "README.md"}
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

func readPluginJSON() (map[string]interface{}, error) {
	byteValue, err := ioutil.ReadFile(path.Join(distDir, pluginJson))
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
