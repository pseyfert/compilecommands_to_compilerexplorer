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
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Project struct {
	Slot    string
	Day     string
	Project string
	Version string
}

func (p *Project) ConfVersion() string {
	return p.Slot + "/" + p.Day + "/" + p.Version
}

type translationunit struct {
	Builddir string `json:"directory"`
	Command  string `json:"command"`
	File     string `json:"file"`
}

var cmtconfig, nightlyroot string

func main() {
	var p Project
	flag.StringVar(&p.Slot, "slot", "lhcb-head", "nightlies slot (i.e. directory in /cvmfs/lhcbdev.cern.ch/nightlies/)")
	flag.StringVar(&p.Day, "day", "Today", "day/buildID (i.e. subdirectory, such as 'Today', 'Mon', or '2032')")
	flag.StringVar(&p.Project, "project", "Brunel", "project (such as Rec, Brunel, LHCb, Lbcom)")
	flag.StringVar(&p.Version, "version", "HEAD", "version (i.e. the stuff after the underscore like HEAD or 2016-patches)")
	flag.StringVar(&cmtconfig, "cmtconfig", "x86_64+avx2+fma-centos7-gcc7-opt", "platform, like x86_64+avx2+fma-centos7-gcc7-opt or x86_64-centos7-gcc7-opt")
	flag.StringVar(&nightlyroot, "nightly-base", "/cvmfs/lhcbdev.cern.ch/nightlies/", "add the specified directory to the nightly builds search path")
	flag.Parse()
	p.Project = strings.ToUpper(p.Project)

	fmt.Println(colon_separate(parse_and_generate(p, nightlyroot, cmtconfig)))
}

func colon_separate(stringset map[string]bool) string {
	var retval string
	addseparator := false
	for k, _ := range stringset {
		if addseparator {
			retval += ":"
		} else {
			addseparator = true
		}
		retval += k
	}
	return retval
}

func parse_and_generate(p Project, nightlyroot, cmtconfig string) map[string]bool {

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
		fmt.Println(err)
		os.Exit(1)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var db []translationunit
	json.Unmarshal(byteValue, &db)

	stringset := make(map[string]bool)
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
				fmt.Print("could not handle %s\n", inc)
				os.Exit(2)
			}

		}
	}
	return stringset
}
