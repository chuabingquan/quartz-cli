package cmd

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

const (
	templateDir = "templates/nodejs"
)

// Project ...
type Project struct {
	AppName      string
	AbsolutePath string
}

type config struct {
	Name     string   `json:"name"`
	Timezone string   `json:"timezone"`
	Schedule []string `json:"schedule"`
}

type pkg struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Main    string `json:"main"`
	License string `json:"license"`
}

// Create ...
func (p *Project) Create() error {
	if _, err := os.Stat(p.AbsolutePath); os.IsNotExist(err) {
		if err := os.Mkdir(p.AbsolutePath, 0754); err != nil {
			return err
		}
	}

	appFile := `'use strict';

module.exports = ctx => {
	// Start writing cron job logic here.
};`

	err := ioutil.WriteFile(p.AbsolutePath+"/app.js", []byte(appFile), os.ModePerm)
	if err != nil {
		return err
	}

	projectConfig := config{
		Name:     p.AppName,
		Timezone: "",
		Schedule: []string{},
	}

	configFile, err := json.MarshalIndent(&projectConfig, "", "\t")
	err = ioutil.WriteFile(p.AbsolutePath+"/config.json", configFile, os.ModePerm)
	if err != nil {
		return err
	}

	projectPkg := pkg{
		Name:    p.AppName,
		Version: "1.0.0",
		Main:    "app.js",
		License: "MIT",
	}

	pkgFile, err := json.MarshalIndent(&projectPkg, "", "\t")
	err = ioutil.WriteFile(p.AbsolutePath+"/package.json", pkgFile, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}
