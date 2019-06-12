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

// This file contains the special treatment of LHCb specific include paths and
// how I want to treat them in our compiler explorer instance.
// Non-LHCb users should not need anything from this file.

package cc2ce4lhcb

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/pseyfert/compilecommands_to_compilerexplorer/cc2ce"
)

// Filter_LHCb_public_includes removes or manipulates include paths from a
// map[string]bool that need special treatment in the setup of the LHCb build
// servers:
//  * Include paths from /cvmfs get accepted
//  * Includes that look like they are (in the) the source directory of the
//    current project get ignored
//  * The header paths installed by the current project get added
//  * Include paths from the current workspace that look like install
//    directories of dependencies (built by the same slot) get manipulated to
//    their expected cvmfs deployment destination
func Filter_LHCb_public_includes(unfiltered map[string]bool, p Project) (map[string]bool, error) {
	filtered, err := Filter_LHCb_includes(unfiltered, p, false)
	return filtered, err
}

func Filter_LHCb_includes(unfiltered map[string]bool, p Project, keep_local_includes bool) (map[string]bool, error) {
	filtered := make(map[string]bool)
	// add the deployed install area of the current project
	filtered[filepath.Join(Installarea(p), "/include")] = true
	for inc, boolean := range unfiltered {
		if !boolean {
			// this is unexpected input
			return make(map[string]bool), fmt.Errorf("Filter_LHCb_all_includes unexpected input: false flagged include path")
		}
		if strings.HasPrefix(inc, "/cvmfs") {
			// accept paths from cvmfs
			filtered[inc] = true
		} else if strings.Contains(inc, "InstallArea") {
			// this looks like the installation area of a dependency project
			// replace /workspace/build/... by something like
			// /cvmfs/lhcbdev.cern.ch/nightlies/lhcb-head/Tue/...
			// where ... looks like GAUDI/GAUDI_master/InstallArea/x86_64+avx2+fma-centos7-gcc7-opt/include
			filtered[strings.Replace(inc, "/workspace/build/", p.Buildarea()+"/", 1)] = true
		} else if strings.HasPrefix(inc, filepath.Join("/workspace/build", p.ProjectareaInBuildarea_new())) {
			// should be the source of the current project (in a new - i.e. nightlies - build area)
			if keep_local_includes {
				filtered[strings.Replace(inc, "/workspace/build/", p.Buildarea()+"/", 1)] = true
			}
		} else if strings.HasPrefix(inc, filepath.Join("/workspace/build", p.ProjectareaInBuildarea_old())) {
			// should be the source of the current project (in an old - i.e. old released - build area)
			if keep_local_includes {
				filtered[strings.Replace(inc, "/workspace/build/", p.Buildarea()+"/", 1)] = true
			}
		} else if inc != "" {
			// includes which are none of the above are unexpected
			return make(map[string]bool), fmt.Errorf("Unexpected include path for LHCb nightly treatment: %s", inc)
		}
	}
	return filtered, nil
}

func (p *Project) CE_config_name() string {
	return strings.ToLower(p.Project)
}

func (p *Project) Buildarea() string {
	if Released {
		return "/cvmfs/lhcb.cern.ch/lib/lhcb"
	}
	return filepath.Join(
		Nightlyroot,
		p.Slot,
		p.Day)
}

func (p *Project) ProjectareaInBuildarea_new() string {
	return p.Project
}

func (p *Project) ProjectareaInBuildarea_old() string {
	return filepath.Join(
		strings.ToUpper(p.Project),
		strings.ToUpper(p.Project)+"_"+p.Version)
}

func Installarea(p Project) string {
	if Released {
		return filepath.Join(
			"/cvmfs/lhcb.cern.ch/lib/lhcb",
			strings.ToUpper(p.Project),
			strings.ToUpper(p.Project)+"_"+p.Version,
			"InstallArea",
			Cmtconfig)
	}
	return filepath.Join(
		Nightlyroot,
		p.Slot,
		p.Day,
		p.Project,
		"InstallArea",
		Cmtconfig)
}

func (p *Project) GenerateIncludes() error {
	incs, err := Parse_and_generate(*p, Nightlyroot, Cmtconfig)
	if err != nil {
		return err
	}
	p.IncludeMap = incs
	return nil
}

func Parse_and_generate(p Project, nightlyroot, cmtconfig string) (map[string]bool, error) {
	stringset := make(map[string]bool)

	unfiltered, err := cc2ce.ParseJsonByFilename(Installarea(p), false)
	if err != nil {
		return stringset, err
	}

	filtered, err := Filter_LHCb_public_includes(unfiltered, p)
	if err != nil {
		return stringset, err
	}

	return filtered, nil
}

// Wrapper of what should become one version of a library in Compiler-Explorer.
// Given the installation of nightlies on cvmfs, this is defined by the
// architecture, slot, day (or build), project name and version.
//
// * Version is usually HEAD.
// * Project must be all upper case
// * Day is the number of the build as string, or the shorthand symlink name (e.g. "Today")
// * Slot is the slot of the nightly build system (e.g. lhcb-head or lhcb-gaudi-head)
// * IncludeMap is the quasi-set of all include paths (the installed ones and the dependencies)
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

// The current platform, e.g. "x86_64-centos7-gcc7-opt"
var Cmtconfig string

// The installation of nightlies, i.e. "/cvmfs/lhcbdev.cern.ch/nightlies"
var Nightlyroot string

var Released bool
