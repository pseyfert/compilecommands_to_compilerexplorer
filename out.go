package main

import (
	"fmt"
	write "github.com/google/go-write"
	"os"
	"strings"
)

func create(p Project) {
	f, err := write.TempFile("", "./c++.local.properties")
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(5)
	}
	defer f.Cleanup()
	// EXAMPLE:
	// ```
	// libs=moore:brunel
	// libs.moore.name=MOORE
	// libs.moore.versions=v30r0
	// libs.moore.url=https://google.com/sorry
	// libs.moore.versions.v30r0.version=v30r0
	// libs.moore.versions.v30r0.path=/cvmfs/lhcb.cern.ch/lib/lhcb/MOORE/MOORE_v30r0/InstallArea/x86_64-centos7-gcc7-opt/include:/cvmfs/lhcb.cern.ch/lib/lcg/releases/LCG_93/Python/2.7.13/x86_64-centos7-gcc7-opt/include/python2.7:/cvmfs/lhcb.cern.ch/lib/lcg/releases/LCG_93/cppgsl/b07383ea/x86_64-centos7-gcc7-opt:/cvmfs/lhcb.cern.ch/lib/lcg/releases/LCG_93/vdt/0.3.9/x86_64-centos7-gcc7-opt/include:/cvmfs/lhcb.cern.ch/lib/lcg/releases/LCG_93/clhep/2.4.0.1/x86_64-centos7-gcc7-opt/include:/cvmfs/lhcb.cern.ch/lib/lcg/releases/LCG_93/GSL/2.1/x86_64-centos7-gcc7-opt/include:/cvmfs/lhcb.cern.ch/lib/lcg/releases/LCG_93/rangev3/0.3.0/x86_64-centos7-gcc7-opt/include:/cvmfs/lhcb.cern.ch/lib/lcg/releases/LCG_93/AIDA/3.2.1/x86_64-centos7-gcc7-opt/src/cpp:/cvmfs/lhcb.cern.ch/lib/lcg/releases/LCG_93/tbb/2018_U1/x86_64-centos7-gcc7-opt/include:/cvmfs/lhcb.cern.ch/lib/lcg/releases/LCG_93/ROOT/6.12.06/x86_64-centos7-gcc7-opt/include:/cvmfs/lhcb.cern.ch/lib/lcg/releases/LCG_93/Boost/1.66.0/x86_64-centos7-gcc7-opt/include:/cvmfs/lhcb.cern.ch/lib/lhcb/HLT/HLT_v30r0/InstallArea/x86_64-centos7-gcc7-opt/include:/cvmfs/lhcb.cern.ch/lib/lhcb/PHYS/PHYS_v30r0/InstallArea/x86_64-centos7-gcc7-opt/include:/cvmfs/lhcb.cern.ch/lib/lhcb/REC/REC_v30r0/InstallArea/x86_64-centos7-gcc7-opt/include:/cvmfs/lhcb.cern.ch/lib/lhcb/LBCOM/LBCOM_v30r0/InstallArea/x86_64-centos7-gcc7-opt/include:/cvmfs/lhcb.cern.ch/lib/lhcb/LHCB/LHCB_v50r0/InstallArea/x86_64-centos7-gcc7-opt/include:/cvmfs/lhcb.cern.ch/lib/lhcb/GAUDI/GAUDI_v30r2/InstallArea/x86_64-centos7-gcc7-opt/include
	// ```
	pr := func(s string) {
		if _, err := fmt.Fprintf(f, "libs.%s.%s\n", strings.ToLower(p.Project), s); err != nil {
			fmt.Print(err.Error())
			os.Exit(5)
		}
	}
	if _, err := fmt.Fprintf(f, "libs=%s\n", strings.ToLower(p.Project)); err != nil {
		fmt.Print(err.Error())
		os.Exit(5)
	}
	pr("name=" + strings.ToUpper(p.Project))
	pr("versions=" + p.ConfVersion())
	pr("versions." + p.ConfVersion() + ".version=" + p.ConfVersion())
	pr("versions." + p.ConfVersion() + ".path=" + colon_separate(parse_and_generate(p, nightlyroot, cmtconfig)))

	if err := f.CloseAtomicallyReplace(); err != nil {
		fmt.Print(err.Error())
		os.Exit(6)
	}
}
