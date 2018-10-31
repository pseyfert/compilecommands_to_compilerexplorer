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
	"fmt"
	"log"
	"os"
	"strings"

	write "github.com/google/go-write"
)

func create(ps []Project) {
	if len(ps) == 0 {
		log.Print("no project?")
		os.Exit(8)
	}
	unique_project_names := make(map[string]bool)
	for _, p := range ps {
		unique_project_names[strings.ToLower(p.Project)] = true
	}
	project_names := colonSeparateMap(unique_project_names)

	f, err := write.TempFile("", "./c++.local.properties")
	if err != nil {
		log.Print("Couldn't create tempfile for output writing: %v", err)
		os.Exit(5)
	}
	defer f.Cleanup()

	if _, err := fmt.Fprintf(f, "libs=%s\n", project_names); err != nil {
		log.Printf("assembling projects %v to c++.local.properties: %v", project_names, err)
		os.Exit(5)
	}

	// EXAMPLE:
	// ```
	// libs=moore:brunel
	// libs.moore.name=MOORE
	// libs.moore.versions=v30r0
	// libs.moore.url=https://google.com/sorry
	// libs.moore.versions.v30r0.version=v30r0
	// libs.moore.versions.v30r0.path=/cvmfs/lhcb.cern.ch/lib/lhcb/MOORE/MOORE_v30r0/InstallArea/x86_64-centos7-gcc7-opt/include:/cvmfs/lhcb.cern.ch/lib/lcg/releases/LCG_93/Python/2.7.13/x86_64-centos7-gcc7-opt/include/python2.7:/cvmfs/lhcb.cern.ch/lib/lcg/releases/LCG_93/cppgsl/b07383ea/x86_64-centos7-gcc7-opt:/cvmfs/lhcb.cern.ch/lib/lcg/releases/LCG_93/vdt/0.3.9/x86_64-centos7-gcc7-opt/include:/cvmfs/lhcb.cern.ch/lib/lcg/releases/LCG_93/clhep/2.4.0.1/x86_64-centos7-gcc7-opt/include:/cvmfs/lhcb.cern.ch/lib/lcg/releases/LCG_93/GSL/2.1/x86_64-centos7-gcc7-opt/include:/cvmfs/lhcb.cern.ch/lib/lcg/releases/LCG_93/rangev3/0.3.0/x86_64-centos7-gcc7-opt/include:/cvmfs/lhcb.cern.ch/lib/lcg/releases/LCG_93/AIDA/3.2.1/x86_64-centos7-gcc7-opt/src/cpp:/cvmfs/lhcb.cern.ch/lib/lcg/releases/LCG_93/tbb/2018_U1/x86_64-centos7-gcc7-opt/include:/cvmfs/lhcb.cern.ch/lib/lcg/releases/LCG_93/ROOT/6.12.06/x86_64-centos7-gcc7-opt/include:/cvmfs/lhcb.cern.ch/lib/lcg/releases/LCG_93/Boost/1.66.0/x86_64-centos7-gcc7-opt/include:/cvmfs/lhcb.cern.ch/lib/lhcb/HLT/HLT_v30r0/InstallArea/x86_64-centos7-gcc7-opt/include:/cvmfs/lhcb.cern.ch/lib/lhcb/PHYS/PHYS_v30r0/InstallArea/x86_64-centos7-gcc7-opt/include:/cvmfs/lhcb.cern.ch/lib/lhcb/REC/REC_v30r0/InstallArea/x86_64-centos7-gcc7-opt/include:/cvmfs/lhcb.cern.ch/lib/lhcb/LBCOM/LBCOM_v30r0/InstallArea/x86_64-centos7-gcc7-opt/include:/cvmfs/lhcb.cern.ch/lib/lhcb/LHCB/LHCB_v50r0/InstallArea/x86_64-centos7-gcc7-opt/include:/cvmfs/lhcb.cern.ch/lib/lhcb/GAUDI/GAUDI_v30r2/InstallArea/x86_64-centos7-gcc7-opt/include
	// ```

	setup_project_names := make(map[string]bool)

	for i, p := range ps {
		if _, found := setup_project_names[p.Project]; !found {
			var versions []string
			for j := i; j < len(ps); j++ {
				if ps[j].Project == p.Project {
					versions = append(versions, ps[j].ConfVersion())
				}
			}
			if _, err := fmt.Fprintf(f, "libs.%s.name=%s\n", strings.ToLower(p.Project), strings.ToLower(p.Project)); err != nil {
				log.Printf("writing project name for %s to c++.local.properties: %v", p.Project, err)
				os.Exit(5)
			}
			if _, err := fmt.Fprintf(f, "libs.%s.url=https://lhcb-nightlies.cern.ch/nightly/summary/\n", strings.ToLower(p.Project)); err != nil {
				log.Printf("writing project url for %s to c++.local.properties: %v", p.Project, err)
				os.Exit(5)
			}

			output_versions := colonSeparateArray(versions)
			if _, err := fmt.Fprintf(f, "libs.%s.versions=%s\n", strings.ToLower(p.Project), output_versions); err != nil {
				log.Printf("writing project %s versions %v to c++.local.properties: %v", strings.ToLower(p.Project), output_versions, err)
				os.Exit(5)
			}

		}
		setup_project_names[p.Project] = true
	}

	for _, p := range ps {
		pr := func(s1, s2 string) {
			if _, err := fmt.Fprintf(f, "libs.%s.versions.%s.%s=%s\n", strings.ToLower(p.Project), p.ConfVersion(), s1, s2); err != nil {
				log.Printf("adding configuration %s=%s to %s/%s: %v", s1, s2, p.Project, p.ConfVersion(), err)
				os.Exit(5)
			}
		}
		pr("version", p.ConfVersion())
		pr("path", colonSeparateMap(p.IncludeMap))
	}
	if err := f.CloseAtomicallyReplace(); err != nil {
		log.Printf("writing c++.local.properties: %v", err)
		os.Exit(6)
	}
}
