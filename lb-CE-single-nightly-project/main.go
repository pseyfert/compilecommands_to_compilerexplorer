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

	"github.com/pseyfert/compilecommands_to_compilerexplorer/cc2ce"
	"github.com/pseyfert/compilecommands_to_compilerexplorer/cc2ce4lhcb"
)

func main() {
	var p cc2ce4lhcb.Project
	var conffilename string
	flag.StringVar(&p.Slot, "slot", "lhcb-head", "nightlies slot (i.e. directory in /cvmfs/lhcbdev.cern.ch/nightlies/)")
	flag.StringVar(&p.Day, "day", "Today", "day/buildID (i.e. subdirectory, such as 'Today', 'Mon', or '2032')")
	flag.StringVar(&p.Project, "project", "Brunel", "project (such as Rec, Brunel, LHCb, Lbcom)")
	flag.StringVar(&p.Version, "version", "HEAD", "version (i.e. the stuff after the underscore like HEAD or 2016-patches)")
	flag.StringVar(&cc2ce4lhcb.Cmtconfig, "cmtconfig", "x86_64+avx2+fma-centos7-gcc7-opt", "platform, like x86_64+avx2+fma-centos7-gcc7-opt or x86_64-centos7-gcc7-opt")
	flag.StringVar(&cc2ce4lhcb.Nightlyroot, "nightly-base", "/cvmfs/lhcbdev.cern.ch/nightlies/", "add the specified directory to the nightly builds search path")
	flag.StringVar(&conffilename, "o", "./c++.local.properties", "output filename")
	flag.BoolVar(&cc2ce4lhcb.Released, "R", false, "look for released projects")
	flag.Parse()
	incs, err := cc2ce4lhcb.Parse_and_generate(p, cc2ce4lhcb.Nightlyroot, cc2ce4lhcb.Cmtconfig)
	if err != nil {
		log.Printf("couldn't read json: %v", err)
		os.Exit(1)
	}

	p.IncludeMap = incs

	fmt.Println(cc2ce.ColonSeparateMap(p.IncludeMap))
	cc2ce4lhcb.Create([]cc2ce4lhcb.Project{p}, conffilename)
}
