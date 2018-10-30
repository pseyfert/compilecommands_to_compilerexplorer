/*
 * Copyright (C) 2018  CERN for the benefit of the LHCb collaboration
 * Author: Paul Seyfert <pseyfert@cern.ch>
 *
 * This software is distributed under the terms of the GNU General Public
 * Licence version 3 (GPL Version 3), copied verbatim in the file "LICENSE".
 *
 * In applying this licence, CERN does not waive the privileges and immunities
 * granted to it by virtue of its status as an Intergovernmental Organization
 * or submit itself to any jurisdiction.
 */

package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type translationunit struct {
	Builddir string `json:"directory"`
	Command  string `json:"command"`
	File     string `json:"file"`
}

type ParseError struct {
	ReadError *os.PathError
	message   string
}

func (e *ParseError) Error() string {
	if e.ReadError != nil {
		return e.ReadError.Error()
	}
	return "Undefined Error"
}

func parse_and_generate(p Project, nightlyroot, cmtconfig string) (map[string]bool, *ParseError) {
	stringset := make(map[string]bool)

	installarea := filepath.Join(
		nightlyroot,
		p.Slot,
		p.Day,
		p.Project,
		p.Project+"_"+p.Version,
		"InstallArea",
		cmtconfig)

	jsonFile, err := os.Open(filepath.Join(installarea, "compile_commands.json"))

	if err != nil {
		if patherr, ok := err.(*os.PathError); ok {
			return stringset, &ParseError{ReadError: patherr, message: ""}
		} else {
			return stringset, &ParseError{ReadError: nil, message: err.Error()}
		}
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var db []translationunit
	json.Unmarshal(byteValue, &db)

	stringset[filepath.Join(installarea, "/include")] = true

	for _, tu := range db {
		words := strings.Fields(tu.Command)
		for j, w := range words {
			var inc string
			if w[0:2] == "-I" {
				inc = w[2:len(w)]
			}
			if w == "-isystem" {
				inc = words[j+1]
			}
			if strings.HasPrefix(inc, "/cvmfs") {
				stringset[inc] = true
			} else if strings.Contains(inc, "InstallArea") {
				stringset[strings.Replace(inc, "/workspace/build/", filepath.Join(nightlyroot, p.Slot, p.Day), 1)] = true
			} else if strings.HasPrefix(inc, filepath.Join("/workspace/build", p.Project, p.Project+"_"+p.Version)) {
				// should be fine, I hope
			} else if inc != "" {
				log.Printf("could not handle compiler argument %s", inc)
				os.Exit(2)
			}

		}
	}
	return stringset, nil
}
