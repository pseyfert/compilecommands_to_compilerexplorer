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

	"github.com/pseyfert/compilecommands_to_compilerexplorer/cc2ce"
)

func main() {
	var lib cc2ce.Library
	var dbpath string
	flag.StringVar(&dbpath, "p", ".", "Compilation database path")
	flag.StringVar(&lib.LibraryName, "l", "local", "Name of library to display in CE")
	flag.StringVar(&lib.LibraryUrl, "u", "", "URL to link from CE")
	flag.StringVar(&lib.LibraryVersion, "version", "master", "version information to display in CE")
	flag.Parse()
	var err error
	lib.Paths, err = cc2ce.ParseJsonByFilename(dbpath, true)
	if err != nil {
		log.Printf("Could not read compile_commands.json: %v", err)
		os.Exit(1)
	}
	err = cc2ce.WriteSingleLibraryAndVersion(lib)
	if err != nil {
		os.Exit(5)
	}
}
