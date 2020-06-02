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
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	write "github.com/google/renameio"
	"github.com/pseyfert/compilecommands_to_compilerexplorer/cc2ce"
)

type CompilerConfig struct {
	Exe      string
	Name     string
	ConfName string
	Options  string
}

func CompilerFromJsonByDB(db []cc2ce.JsonTranslationunit) (string, error) {
	var b bytes.Buffer
	for _, tu := range db {
		words := strings.Fields(tu.Command)
		for i, w := range words {
			if strings.HasPrefix(w, "-") || strings.HasSuffix(w, ".cpp") {
				break
			}
			if i != 0 {
				b.WriteString(" ")
			}
			b.WriteString(w)
		}
		return b.String(), nil
	}
	return "", fmt.Errorf("no translation units found")
}

func main() {
	var lib cc2ce.Library
	var dbpath string
	flag.StringVar(&dbpath, "p", ".", "Compilation database path")
	flag.StringVar(&lib.LibraryName, "l", "local", "Name of library to display in CE")
	flag.StringVar(&lib.LibraryUrl, "u", "", "URL to link from CE")
	flag.StringVar(&lib.LibraryVersion, "version", "master", "version information to display in CE")
	ofname := flag.String("o", "./c++.local.properties", "output file with CE configuration")
	flag.Parse()
	var err error
	turnAbsolute := true
	db, err := cc2ce.JsonTUsByFilename(dbpath)
	if err != nil {
		log.Printf("Could not read compile_commands.json: %v", err)
		os.Exit(1)
	}
	lib.Paths, err = cc2ce.IncludesFromJsonByDB(db, turnAbsolute)
	if err != nil {
		log.Printf("reading of include paths failed: %v", err)
		os.Exit(1)
	}

	var compiler CompilerConfig
	compiler.Exe, err = CompilerFromJsonByDB(db)
	compiler.Options, err = cc2ce.OptionsFromJsonByDB(db, false)
	compiler.Name = "hardcoded"
	compiler.ConfName = "hardcoded"
	if err != nil {
		log.Printf("Error obtaining compiler options: %v", err)
		os.Exit(1)
	}

	f, err := write.TempFile("", *ofname)
	if err != nil {
		log.Printf("Couldn't create tempfile for output writing: %v", err)
		os.Exit(5)
	}
	err = cc2ce.WriteSingleLibraryAndVersionToFile(lib, f)
	if err != nil {
		log.Printf("Error writing library config: %v", err)
		f.Cleanup()
		os.Exit(5)
	}
	err = WriteConfig([]CompilerConfig{compiler}, f)
	if err != nil {
		log.Printf("Error writing compiler config: %v", err)
		f.Cleanup()
		os.Exit(5)
	}
	if err := f.CloseAtomicallyReplace(); err != nil {
		log.Printf("writing %s failed: %v", *ofname, err)
		f.Cleanup()
		os.Exit(5)
	}
	f.Cleanup()
	os.Exit(0)
}

func WriteConfig(confs []CompilerConfig, f *write.PendingFile) error {
	fmt.Fprintf(f, "compilers=&autogen")
	{
		var b bytes.Buffer
		addseparator := false
		for _, c := range confs {
			if addseparator {
				b.WriteString(":")
			} else {
				addseparator = true
			}
			b.WriteString(c.ConfName)
		}
		if _, err := fmt.Fprintf(f, "group.autogen.compilers=%s\n", b.String()); err != nil {
			log.Printf("Error writing to config: %v", err)
			return err
		}
	}
	if _, err := fmt.Fprint(f, "group.autogen.groupName=auto-generated compiler settings\n"); err != nil {
		log.Printf("Error writing to config: %v", err)
		return err
	}
	compiler_writer := func(c CompilerConfig) error {
		if _, err := fmt.Fprintf(f, "compiler.%s.name=%s\n", c.ConfName, c.Name); err != nil {
			log.Printf("Error writing to config: %v", err)
			return err
		}
		if _, err := fmt.Fprintf(f, "compiler.%s.exe=%s\n", c.ConfName, c.Exe); err != nil {
			log.Printf("Error writing to config: %v", err)
			return err
		}
		if _, err := fmt.Fprintf(f, "compiler.%s.options=%s\n", c.ConfName, c.Options); err != nil {
			log.Printf("Error writing to config: %v", err)
			return err
		}
		return nil
	}
	for _, c := range confs {
		if err := compiler_writer(c); err != nil {
			return err
		}
	}
	return nil
}
