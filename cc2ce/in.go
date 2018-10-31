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

// This file contains the parsing of a compile_commands.json database to
// extract a union (over all translation units) of all include paths given to
// the compiler.

package cc2ce

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type translationunit struct {
	Builddir string `json:"directory"`
	Command  string `json:"command"`
	File     string `json:"file"`
}

// ParseJsonByBytes parses json provided as []byte (and is called by
// ParseJsonByFilename). It collects all include paths given in the form
// `-Isomepath` and `-isystem somepath`.
// The return is a quasi-'set' of strings: a map string -> bool.  All bools are
// true. An include path is present in the compile_commands if and only if it
// is present as key in the map.
func ParseJsonByBytes(inFileContent []byte) (map[string]bool, error) {
	stringset := make(map[string]bool)
	var db []translationunit
	json.Unmarshal(inFileContent, &db)

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
			stringset[inc] = true
		}
	}
	return stringset, nil
}

// ParseJsonByFilename opens a compile_commands.json file and passes it to
// ParseJsonByBytes to get the union of all include paths. If the argument ends
// on "compile_commands.json", it is assumed to be the path to the
// compile_commands.json file. Otherwise, it is assumed to be the directory
// containing the compile_commands.json file.
// The return is a quasi-'set' of strings: a map string -> bool.  All bools are
// true. An include path is present in the compile_commands if and only if it
// is present as key in the map.
func ParseJsonByFilename(inFileName string) (map[string]bool, error) {
	stringset := make(map[string]bool)

	if !strings.HasSuffix(inFileName, "compile_commands.json") {
		inFileName = filepath.Join(inFileName, "compile_commands.json")
	}
	jsonFile, err := os.Open(inFileName)
	if err != nil {
		return stringset, err
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

	stringset, err = ParseJsonByBytes(byteValue)
	return stringset, err
}
