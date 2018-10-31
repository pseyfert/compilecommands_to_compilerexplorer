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

package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

// filter_LHCb_includes removes or manipulates include paths from a
// map[string]bool that need special treatment in the setup of the LHCb build
// servers:
//  * Include paths from /cvmfs get accepted
//  * Includes that look like they are (in the) the source directory of the
//    current project get ignored
//  * The header paths installed by the current project get added
//  * Include paths from the current workspace that look like install
//    directories of dependencies (built by the same slot) get manipulated to
//    their expected cvmfs deployment destination
func filter_LHCb_includes(unfiltered map[string]bool, p Project) (map[string]bool, error) {
	filtered := make(map[string]bool)
	// add the deployed install area of the current project
	filtered[filepath.Join(installarea(p), "/include")] = true
	for inc, boolean := range unfiltered {
		if !boolean {
			// this is unexpected input
			return make(map[string]bool), fmt.Errorf("filter_LHCb_includes unexpected input: false flagged include path")
		}
		if strings.HasPrefix(inc, "/cvmfs") {
			// accept paths from cvmfs
			filtered[inc] = true
		} else if strings.Contains(inc, "InstallArea") {
			// this looks like the installation area of a dependency project
			// replace /workspace/build/... by something like
			// /cvmfs/lhcbdev.cern.ch/nightlies/lhcb-head/Tue/...
			// where ... looks like GAUDI/GAUDI_master/InstallArea/x86_64+avx2+fma-centos7-gcc7-opt/include
			filtered[strings.Replace(inc, "/workspace/build/", filepath.Join(nightlyroot, p.Slot, p.Day)+"/", 1)] = true
		} else if strings.HasPrefix(inc, filepath.Join("/workspace/build", p.Project, p.Project+"_"+p.Version)) {
			// should be the source of the current project and is skipped
		} else if inc != "" {
			// includes which are none of the above are unexpected
			return make(map[string]bool), fmt.Errorf("Unexpected include path for LHCb nightly treatment: %s", inc)
		}
	}
	return filtered, nil
}

func installarea(p Project) string {
	return filepath.Join(
		nightlyroot,
		p.Slot,
		p.Day,
		p.Project,
		p.Project+"_"+p.Version,
		"InstallArea",
		cmtconfig)
}

func parse_and_generate(p Project, nightlyroot, cmtconfig string) (map[string]bool, error) {
	stringset := make(map[string]bool)

	unfiltered, err := ParseJsonByFilename(installarea(p))
	if err != nil {
		return stringset, err
	}

	filtered, err := filter_LHCb_includes(unfiltered, p)
	if err != nil {
		return stringset, err
	}

	return filtered, nil
}
