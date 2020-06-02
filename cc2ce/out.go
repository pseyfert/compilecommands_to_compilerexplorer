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

package cc2ce

import (
	"fmt"
	"log"
	"strings"

	write "github.com/google/renameio"
)

type Library struct {
	LibraryName    string
	LibraryVersion string
	LibraryUrl     string
	Paths          map[string]bool
}

func WriteSingleLibraryAndVersionToFile(lib Library, f *write.PendingFile) error {
	if _, err := fmt.Fprintf(f, "libs=%s\n", strings.ToLower(lib.LibraryName)); err != nil {
		log.Printf("writing to c++.local.properties failed: %v", err)
		return err
	}

	print_lib := func(key, val string) error {
		if _, err := fmt.Fprintf(f, "libs.%s.%s=%s\n", strings.ToLower(lib.LibraryName), key, val); err != nil {
			log.Printf("writing to c++.local.properties failed: %v", err)
			return err
		}
		return nil
	}

	err := print_lib("name", lib.LibraryName)
	if err != nil {
		return err
	}
	if lib.LibraryUrl != "" {
		err = print_lib("url", lib.LibraryUrl)
		if err != nil {
			return err
		}
	}
	err = print_lib("versions", lib.LibraryVersion)
	if err != nil {
		return err
	}

	print_lib_ver := func(key, val string) error {
		if _, err := fmt.Fprintf(f, "libs.%s.versions.%s.%s=%s\n", strings.ToLower(lib.LibraryName), lib.LibraryVersion, key, val); err != nil {
			log.Printf("writing to c++.local.properties failed: %v", err)
			return err
		}
		return nil
	}
	err = print_lib_ver("version", lib.LibraryVersion)
	if err != nil {
		return err
	}
	err = print_lib_ver("path", ColonSeparateMap(lib.Paths))
	if err != nil {
		return err
	}
	return nil
}

func WriteSingleLibraryAndVersion(lib Library) error {
	f, err := write.TempFile("", "./c++.local.properties")
	if err != nil {
		log.Printf("Couldn't create tempfile for output writing: %v", err)
		return err
	}
	defer f.Cleanup()
	err = WriteSingleLibraryAndVersionToFile(lib, f)
	if err != nil {
		return err
	}
	if err := f.CloseAtomicallyReplace(); err != nil {
		log.Printf("writing c++.local.properties failed: %v", err)
		return err
	}
	return nil
}
