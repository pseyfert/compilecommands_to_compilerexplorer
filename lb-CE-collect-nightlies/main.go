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
	"flag"
	"log"
	"os"

	"github.com/pseyfert/compilecommands_to_compilerexplorer/cc2ce4lhcb"
)

func main() {
	days := []string{"latest", "Today", "Yesterday", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}
	slots := []string{"lhcb-head", "lhcb-gaudi-head"}
	top_projects := []string{"Brunel", "Gaudi"}
	versions := []string{"HEAD", "master"}
	var conffilename string
	flag.StringVar(&cc2ce4lhcb.Cmtconfig, "cmtconfig", "x86_64+avx2+fma-centos7-gcc8-opt", "platform, like x86_64+avx2+fma-centos7-gcc7-opt or x86_64-centos7-gcc7-opt")
	flag.StringVar(&cc2ce4lhcb.Nightlyroot, "nightly-base", "/cvmfs/lhcbdev.cern.ch/nightlies/", "add the specified directory to the nightly builds search path")
	flag.StringVar(&conffilename, "o", "./c++.local.properties", "output filename")
	flag.Parse()

	projects := []cc2ce4lhcb.Project{}

	for _, slot := range slots {
		for _, day := range days {
			for _, top_project := range top_projects {
				for _, version := range versions {
					var p cc2ce4lhcb.Project
					p.Slot = slot
					p.Day = day
					p.Project = top_project
					p.Version = version
					incs, err := cc2ce4lhcb.Parse_and_generate(p, cc2ce4lhcb.Nightlyroot, cc2ce4lhcb.Cmtconfig)
					if err != nil {
						if os.IsNotExist(err) {
							log.Printf("configuration doesn't exist: %v", err)
							// this slot doesn't exist on cvmfs (not set up for publication, or build failed)
							// just skip
						} else {
							log.Printf("%v", err)
							os.Exit(7)
						}
					} else {
						p.IncludeMap = incs
						projects = append(projects, p)
					}
				}
			}
		}
	}
	cc2ce4lhcb.Create(projects, conffilename)
}
