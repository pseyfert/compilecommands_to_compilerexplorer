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
	"fmt"
	"log"
	"os"
	"strings"
)

type Project struct {
	Slot       string
	Day        string
	Project    string
	Version    string
	IncludeMap map[string]bool
}

func (p *Project) ConfVersion() string {
	return p.Slot + "/" + p.Day + "/" + p.Version
}

var cmtconfig, nightlyroot string

func cli() {
	var p Project
	flag.StringVar(&p.Slot, "slot", "lhcb-head", "nightlies slot (i.e. directory in /cvmfs/lhcbdev.cern.ch/nightlies/)")
	flag.StringVar(&p.Day, "day", "Today", "day/buildID (i.e. subdirectory, such as 'Today', 'Mon', or '2032')")
	flag.StringVar(&p.Project, "project", "Brunel", "project (such as Rec, Brunel, LHCb, Lbcom)")
	flag.StringVar(&p.Version, "version", "HEAD", "version (i.e. the stuff after the underscore like HEAD or 2016-patches)")
	flag.StringVar(&cmtconfig, "cmtconfig", "x86_64+avx2+fma-centos7-gcc7-opt", "platform, like x86_64+avx2+fma-centos7-gcc7-opt or x86_64-centos7-gcc7-opt")
	flag.StringVar(&nightlyroot, "nightly-base", "/cvmfs/lhcbdev.cern.ch/nightlies/", "add the specified directory to the nightly builds search path")
	flag.Parse()
	p.Project = strings.ToUpper(p.Project)
	incs, err := parse_and_generate(p, nightlyroot, cmtconfig)
	if err != nil {
		log.Printf("couldn't read json: %v", err)
		os.Exit(1)
	}

	p.IncludeMap = incs

	fmt.Println(colonSeparateMap(p.IncludeMap))
	create([]Project{p})
}

func main() {
	days := []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun", "Today"}
	slots := []string{"lhcb-head", "lhcb-gaudi-head"}
	top_projects := []string{"Brunel", "Gaudi"}
	versions := []string{"HEAD"}
	flag.StringVar(&cmtconfig, "cmtconfig", "x86_64+avx2+fma-centos7-gcc7-opt", "platform, like x86_64+avx2+fma-centos7-gcc7-opt or x86_64-centos7-gcc7-opt")
	flag.StringVar(&nightlyroot, "nightly-base", "/cvmfs/lhcbdev.cern.ch/nightlies/", "add the specified directory to the nightly builds search path")
	flag.Parse()

	projects := []Project{}

	for _, day := range days {
		for _, slot := range slots {
			for _, top_project := range top_projects {
				for _, version := range versions {
					var p Project
					p.Slot = slot
					p.Day = day
					p.Project = strings.ToUpper(top_project)
					p.Version = version
					incs, err := parse_and_generate(p, nightlyroot, cmtconfig)
					if err != nil {
						if err.ReadError != nil {
							log.Printf("read error: %v", err.ReadError)
							// probably this slot doesn't exist on cvmfs (not set up for publication, or build failed)
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
	create(projects)
}
