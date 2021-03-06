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
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type JsonTranslationunit struct {
	Builddir string `json:"directory"` // working dir, necessary for relative paths
	Command  string `json:"command"`   // contains the compiler call
	File     string `json:"file"`      // input file

	Arguments []string `json:"arguments"` // alternative to 'command' (list of strings rather than string) FIXME: handle
	Output    string   `json:"output"`    // optional, unused
}

// IncludesFromJsonByBytes parses json provided as []byte (and is called by
// ParseJsonByFilename). It collects all include paths given in the form
// `-Isomepath` and `-isystem somepath`.
//
// The return is a quasi-'set' of strings: a map string -> bool.  All bools are
// true. An include path is present in the compile_commands if and only if it
// is present as key in the map.
//
// When the turnAbsolute option is true, relative paths get turned into
// absolute paths by using the specified working directory from the json.
// Otherwise, no path manipulation is done.
func IncludesFromJsonByBytes(inFileContent []byte, turnAbsolute bool) (map[string]bool, error) {
	db, err := JsonTUsByBytes(inFileContent)

	if nil != err {
		return make(map[string]bool), err
	}

	return IncludesFromJsonByDB(db, turnAbsolute)
}

func IncludesFromJsonByDB(db []JsonTranslationunit, turnAbsolute bool) (map[string]bool, error) {
	stringset := make(map[string]bool)
	for _, tu := range db {
		words := strings.Fields(tu.Command)
		for j, w := range words {
			inc := ""
			if w[0:2] == "-I" {
				inc = w[2:len(w)]
			}
			if w == "-isystem" {
				inc = words[j+1]
			}
			if inc != "" {
				if !filepath.IsAbs(inc) && turnAbsolute {
					inc = filepath.Join(tu.Builddir, inc)
				}
				stringset[inc] = true
			}
		}
	}
	return stringset, nil
}

// Attempt to get compiler options from the compile_commands.json. On a pure
// luck based approach, the compile command of the first translation unit is
// used to extract -W, -m, -f, -p, -std, -O, and -D settings.  These may well
// differ from one translation unit to the other.
//
// The -D options are filtered based on what I found not useful in LHCb
// projects.
func OptionsFromJsonByBytes(inFileContent []byte, skippackagenameversion bool) (string, error) {
	db, err := JsonTUsByBytes(inFileContent)
	if nil != err {
		return "", err
	}
	optionsstring, err := OptionsFromJsonByDB(db, skippackagenameversion)
	return optionsstring, err
}

func OptionsFromJsonByDB(db []JsonTranslationunit, skippackagenameversion bool) (string, error) {
	var b bytes.Buffer
	for _, tu := range db {
		words := strings.Fields(tu.Command)
		for _, w := range words {
			if strings.HasPrefix(w, "-D") {
				if strings.HasSuffix(w, "EXPORTS") {
					continue
				} else if strings.HasPrefix(w, "-DPACKAGE_NAME") {
					if !skippackagenameversion {
						b.WriteString("-DPACKAGE_NAME=\"CompilerExplorer\"")
					}
				} else if strings.HasPrefix(w, "-DPACKAGE_VERSION") {
					if !skippackagenameversion {
						b.WriteString("-DPACKAGE_VERSION=\"v0r0\"")
					}
				} else if w == "-DGAUDI_LINKER_LIBRARY" {
					continue
				} else {
					// In the .json I often see -Dsomevar=\\\"someval\\\"
					// For the .properties this needs to be -Dsomevar="someval" with all backslashes gone
					// NB: in the past I used 7 backslashes instead of 3. At the time of writing this comment
					//     I have to change 7 to 3 to fix production. I did not verify if this is a change
					//     in go (encoding/json?). According to the previous comment, this is not a change
					//     in cmake (observations match).
					b.WriteString(strings.Replace(w, "\\\"", "\"", 2))
				}
			} else if strings.HasPrefix(w, "-p") {
				b.WriteString(w)
			} else if strings.HasPrefix(w, "-O") {
				b.WriteString(w)
			} else if strings.HasPrefix(w, "-m") {
				b.WriteString(w)
			} else if strings.HasPrefix(w, "-f") {
				b.WriteString(w)
			} else if strings.HasPrefix(w, "-W") {
				b.WriteString(w)
			} else if strings.HasPrefix(w, "-std") {
				b.WriteString(w)
			} else {
				continue
			}
			b.WriteString(" ")
		}
		return b.String(), nil
	}
	return "", fmt.Errorf("no translation units found")
}

// ParseJsonByFilename opens a compile_commands.json file and passes it to
// IncludesFromJsonByBytes to get the union of all include paths. If the
// argument ends on "compile_commands.json", it is assumed to be the path to
// the compile_commands.json file. Otherwise, it is assumed to be the directory
// containing the compile_commands.json file.
//
// The return is a quasi-'set' of strings: a map string -> bool.  All bools are
// true. An include path is present in the compile_commands if and only if it
// is present as key in the map.
//
// When the turnAbsolute option is true, relative paths get turned into
// absolute paths by using the specified working directory from the json.
// Otherwise, no path manipulation is done.
func ParseJsonByFilename(inFileName string, turnAbsolute bool) (map[string]bool, error) {
	stringset := make(map[string]bool)
	db, err := JsonTUsByFilename(inFileName)
	if nil != err {
		return stringset, err
	}
	return IncludesFromJsonByDB(db, turnAbsolute)
}

func BytesFromFilename(inFileName string) ([]byte, error) {
	if !strings.HasSuffix(inFileName, "compile_commands.json") {
		inFileName = filepath.Join(inFileName, "compile_commands.json")
	}
	jsonFile, err := os.Open(inFileName)
	if err != nil {
		retval := make([]byte, 0)
		return retval, err
	}
	defer jsonFile.Close()
	byteValue, err := ioutil.ReadAll(jsonFile)

	return byteValue, err
}

func JsonTUsByBytes(inFileContent []byte) ([]JsonTranslationunit, error) {
	var db []JsonTranslationunit
	json.Unmarshal(inFileContent, &db)

	return db, nil
}

func JsonTUsByFilename(inFileName string) ([]JsonTranslationunit, error) {
	inFileContent, err := BytesFromFilename(inFileName)
	if nil != err {
		return make([]JsonTranslationunit, 0), err
	}
	db, err := JsonTUsByBytes(inFileContent)
	return db, err
}
