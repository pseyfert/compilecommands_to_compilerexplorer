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
	"strings"
)

type translationunit struct {
	Builddir string `json:"directory"`
	Command  string `json:"command"`
	File     string `json:"file"`
}

func main() {
	var slot, day, project, version, cmtconfig string
	flag.StringVar(&slot, "slot", "lhcb-head", "nightlies slot (i.e. directory in /cvmfs/lhcbdev.cern.ch/nightlies/)")
	flag.StringVar(&day, "day", "Today", "day/buildID (i.e. subdirectory, such as 'Today', 'Mon', or '2032')")
	flag.StringVar(&project, "project", "Brunel", "project (such as Rec, Brunel, LHCb, Lbcom)")
	flag.StringVar(&version, "version", "HEAD", "version (i.e. the stuff after the underscore like HEAD or 2016-patches)")
	flag.StringVar(&cmtconfig, "cmtconfig", "x86_64+avx2+fma-centos7-gcc7-opt", "platform, like x86_64+avx2+fma-centos7-gcc7-opt or x86_64-centos7-gcc7-opt")
	flag.Parse()
	project = strings.ToUpper(project)

	jsonFile, err := os.Open("/cvmfs/lhcbdev.cern.ch/nightlies/" + slot + "/" + day + "/" + project + "/" + project + "_" + version + "/InstallArea/" + cmtconfig + "/compile_commands.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var db []translationunit
	json.Unmarshal(byteValue, &db)

	stringset := make(map[string]bool)
	stringset["/cvmfs/lhcbdev.cern.ch/nightlies/"+slot+"/"+day+"/"+project+"/"+project+"_"+version+"/InstallArea/"+cmtconfig+"/include"] = true

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
				stringset[strings.Replace(inc, "/workspace/build/", "/cvmfs/lhcbdev.cern.ch/nightlies/"+slot+"/"+day+"/", 1)] = true
			} else if strings.HasPrefix(inc, "/workspace/build/"+project+"/"+project+"_"+version) {
				// should be fine, I hope
			} else if inc != "" {
				fmt.Print("could not handle %s\n", inc)
				os.Exit(2)
			}

		}
	}

	addseparator := false
	for k, _ := range stringset {
		if addseparator {
			fmt.Print(":")
		} else {
			addseparator = true
		}
		fmt.Print(k)
	}

}
