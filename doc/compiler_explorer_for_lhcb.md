Title: Compiler Explorer for LHCb
Date: 2018-10-31 22:48
Category: Computer
Tags: web, computing, physics-software
Authors: Paul Seyfert
Summary: I started running Compiler-Explorer (aka godbolt) for the LHCb software stack. Considering a typical compiler call is 2k characters, that needed some setup. The most relevant part was setting up all the `-I` such that one can `#include "Event/Track.h"` our standard data structure headers.

## Compiler Explorer

You might be familiar with [Compiler Explorer](https://godbolt.org). A
website where you can enter some C++ (or other languages) code and with
then compiles your code and shows the assembly, or preprocessor output,
or intermediate compiler results. This is to enable programmers to
understand better what the compiler is doing, the impact of compiler
options, compare different implementations (if they result in different
code) …

I started a VM on our openstack instance and allowed a few colleagues to use
the site.

## The first problem I'm trying to solve (I got greedy, less relevant ones further down)

One thing I didn't like is that on the website I could use none of the headers
of our LHCb software stack because the website offers many popular C++
libraries as selectable include paths, but not our stack.

### What Compiler Explorer provides

The compiler explorer code is hosted on github and pretty easy to get running
(a decent node.js version provided).

Adding libraries to a local instance is rather straight forward. The
documentation says one should look at `etc/config/c++.amazon.properties` how it
works and then add the corresponding lines to a
`etc/config/c++.<something>.properties` file. `<something>` in my case is
`defaults` or `local`. The difference being, the latter is gitignore'd.

They look like this
```
libs=moore:local
libs.local.name=LOCAL
libs.local.versions=current
libs.local.url=http://127.0.0.1
libs.local.versions.current.name=LOCAL_current
libs.local.versions.current.path=/home/pseyfert/.local/include

libs.moore.name=MOORE
libs.moore.versions=v30r0
libs.moore.url=https://google.com/sorry
libs.moore.versions.v30r0.name=Moore_v30r0
libs.moore.versions.v30r0.path=/cvmfs/lhcb.cern.ch/lib/lhcb/MOORE/MOORE_v30r0/InstallArea/x86_64-centos7-gcc7-opt/include
```

When these files get updated, Compiler Explorer automatically restarts.

### What needs additions from my side

The next problem are transitive dependencies. I.e. I might add the installed
header dirs of the Brunel project to my instance, but basically all headers
would include something from the Rec, LHCb, or Gaudi projects, which then
depend on our Boost version, VDT, ROOT, the Microsoft guideline support library
(GSL), the GNU scientific library (GSL, I'm not kidding we have a naming
collision here), and so on (and all in the right version).

My point is, it's not too much to ask a user "If you use Rec, ROOT, and Boost,
then click all three of them", but a transitive include would turn this into a
frustrating try-and-fix-compilation-error cycle.

One way to get these right is taking a look at the
`compile_commands.json` database as written by cmake which contains
effectively a list of all compilations without any build framework
variables:

```json
[
{
  "directory": "/workspace/build/BRUNEL/BRUNEL_HEAD/build/Rec/BrunelCache",
  "command": "/afs/cern.ch/work/m/marcocle/workspace/LbScripts/LbUtils/scripts/lcg-g++-7.3.0  -DBOOST_FILESYSTEM_VERSION=3 -DBOOST_SPIRIT_USE_PHOENIX_V3 -DBrunel_FunctorCache_EXPORTS -DGAUDI_V20_COMPAT -DPACKAGE_NAME=\\\"BrunelCache\\\" -DPACKAGE_VERSION=\\\"HEAD\\\" -D_GNU_SOURCE -Df2cFortran -Dlinux -Dunix -I/workspace/build/BRUNEL/BRUNEL_HEAD/build/Rec/BrunelCache -I/workspace/build/BRUNEL/BRUNEL_HEAD/Rec/BrunelCache -I/cvmfs/lhcb.cern.ch/lib/lcg/releases/LCG_93/Python/2.7.13/x86_64-centos7-gcc7-opt/include/python2.7 -I/cvmfs/lhcb.cern.ch/lib/lcg/releases/LCG_93/cppgsl/b07383ea/x86_64-centos7-gcc7-opt -I/cvmfs/lhcb.cern.ch/lib/lcg/releases/LCG_93/vdt/0.3.9/x86_64-centos7-gcc7-opt/include -I/cvmfs/lhcb.cern.ch/lib/lcg/releases/LCG_93/clhep/2.4.0.1/x86_64-centos7-gcc7-opt/include -I/cvmfs/lhcb.cern.ch/lib/lcg/releases/LCG_93/GSL/2.1/x86_64-centos7-gcc7-opt/include -I/cvmfs/lhcb.cern.ch/lib/lcg/releases/LCG_93/rangev3/0.3.0/x86_64-centos7-gcc7-opt/include -I/cvmfs/lhcb.cern.ch/lib/lcg/releases/LCG_93/AIDA/3.2.1/x86_64-centos7-gcc7-opt/src/cpp -I/cvmfs/lhcb.cern.ch/lib/lcg/releases/LCG_93/tbb/2018_U1/x86_64-centos7-gcc7-opt/include -isystem /cvmfs/lhcb.cern.ch/lib/lcg/releases/LCG_93/ROOT/6.12.06/x86_64-centos7-gcc7-opt/include -isystem /cvmfs/lhcb.cern.ch/lib/lcg/releases/LCG_93/Boost/1.66.0/x86_64-centos7-gcc7-opt/include -I/workspace/build/BRUNEL/BRUNEL_HEAD -I/workspace/build/BRUNEL/BRUNEL_HEAD/build/include -I/workspace/build/REC/REC_HEAD/InstallArea/x86_64+avx2+fma-centos7-gcc7-opt/include -I/workspace/build/LBCOM/LBCOM_HEAD/InstallArea/x86_64+avx2+fma-centos7-gcc7-opt/include -I/workspace/build/LHCB/LHCB_HEAD/InstallArea/x86_64+avx2+fma-centos7-gcc7-opt/include -I/workspace/build/GAUDI/GAUDI_master/InstallArea/x86_64+avx2+fma-centos7-gcc7-opt/include   -mavx2 -mfma -fmessage-length=0 -pipe -Wall -Wextra -Werror=return-type -pthread -pedantic -Wwrite-strings -Wpointer-arith -Woverloaded-virtual -Wsuggest-override -std=c++17 -fdiagnostics-color -O3 -DNDEBUG -fPIC   -o CMakeFiles/Brunel_FunctorCache.dir/Brunel_FunctorCache_srcs/FUNCTORS_Hlt2HltFactory_0001.cpp.o -c /workspace/build/BRUNEL/BRUNEL_HEAD/build/Rec/BrunelCache/Brunel_FunctorCache_srcs/FUNCTORS_Hlt2HltFactory_0001.cpp",
  "file": "/workspace/build/BRUNEL/BRUNEL_HEAD/build/Rec/BrunelCache/Brunel_FunctorCache_srcs/FUNCTORS_Hlt2HltFactory_0001.cpp"
},
{
  "directory": "/workspace/build/BRUNEL/BRUNEL_HEAD/build/Rec/BrunelCache",
  "command": "/afs/cern.ch/work/m/marcocle/workspace/LbScripts/LbUtils/scripts/lcg-g++-7.3.0  -DBOOST_FILESYSTEM_VERSION=3 -DBOOST_SPIRIT_USE_PHOENIX_V3 -DBrunel_FunctorCache_EXPORTS -DGAUDI_V20_COMPAT -DPACKAGE_NAME=\\\"BrunelCache\\\" -DPACKAGE_VERSION=\\\"HEAD\\\" -D_GNU_SOURCE -Df2cFortran -Dlinux -Dunix -I/workspace/build/BRUNEL/BRUNEL_HEAD/build/Rec/BrunelCache -I/workspace/build/BRUNEL/BRUNEL_HEAD/Rec/BrunelCache -I/cvmfs/lhcb.cern.ch/lib/lcg/releases/LCG_93/Python/2.7.13/x86_64-centos7-gcc7-opt/include/python2.7 -I/cvmfs/lhcb.cern.ch/lib/lcg/releases/LCG_93/cppgsl/b07383ea/x86_64-centos7-gcc7-opt -I/cvmfs/lhcb.cern.ch/lib/lcg/releases/LCG_93/vdt/0.3.9/x86_64-centos7-gcc7-opt/include -I/cvmfs/lhcb.cern.ch/lib/lcg/releases/LCG_93/clhep/2.4.0.1/x86_64-centos7-gcc7-opt/include -I/cvmfs/lhcb.cern.ch/lib/lcg/releases/LCG_93/GSL/2.1/x86_64-centos7-gcc7-opt/include -I/cvmfs/lhcb.cern.ch/lib/lcg/releases/LCG_93/rangev3/0.3.0/x86_64-centos7-gcc7-opt/include -I/cvmfs/lhcb.cern.ch/lib/lcg/releases/LCG_93/AIDA/3.2.1/x86_64-centos7-gcc7-opt/src/cpp -I/cvmfs/lhcb.cern.ch/lib/lcg/releases/LCG_93/tbb/2018_U1/x86_64-centos7-gcc7-opt/include -isystem /cvmfs/lhcb.cern.ch/lib/lcg/releases/LCG_93/ROOT/6.12.06/x86_64-centos7-gcc7-opt/include -isystem /cvmfs/lhcb.cern.ch/lib/lcg/releases/LCG_93/Boost/1.66.0/x86_64-centos7-gcc7-opt/include -I/workspace/build/BRUNEL/BRUNEL_HEAD -I/workspace/build/BRUNEL/BRUNEL_HEAD/build/include -I/workspace/build/REC/REC_HEAD/InstallArea/x86_64+avx2+fma-centos7-gcc7-opt/include -I/workspace/build/LBCOM/LBCOM_HEAD/InstallArea/x86_64+avx2+fma-centos7-gcc7-opt/include -I/workspace/build/LHCB/LHCB_HEAD/InstallArea/x86_64+avx2+fma-centos7-gcc7-opt/include -I/workspace/build/GAUDI/GAUDI_master/InstallArea/x86_64+avx2+fma-centos7-gcc7-opt/include   -mavx2 -mfma -fmessage-length=0 -pipe -Wall -Wextra -Werror=return-type -pthread -pedantic -Wwrite-strings -Wpointer-arith -Woverloaded-virtual -Wsuggest-override -std=c++17 -fdiagnostics-color -O3 -DNDEBUG -fPIC   -o CMakeFiles/Brunel_FunctorCache.dir/Brunel_FunctorCache_srcs/FUNCTORS_Hlt1HltFactory_0001.cpp.o -c /workspace/build/BRUNEL/BRUNEL_HEAD/build/Rec/BrunelCache/Brunel_FunctorCache_srcs/FUNCTORS_Hlt1HltFactory_0001.cpp",
  "file": "/workspace/build/BRUNEL/BRUNEL_HEAD/build/Rec/BrunelCache/Brunel_FunctorCache_srcs/FUNCTORS_Hlt1HltFactory_0001.cpp"
},
...
```

These are not relocatable (or … at least not without some work). The
`/workspace/build` paths are all local on the build machine and can be source
directories, or install directories of dependencies, or build directories (for
copied or generated files). The `/workspace/build/.../InstallArea/` paths get
deployed to cvmfs on `/cvmfs/lhcbdev.cern.ch/nightlies/...` (where the `...` is
not the same in the two paths) in case of the nightly build. Released versions
end up on `/cvmfs/lhcb.cern.ch/lib/lhcb/...` but I don't care about them at the
moment.

I work with the following assumptions:

 * If a user picks the Gaudi project to work with, they intend to work with
   code that is intended to go into Gaudi. So they may rely on public headers
   of "packages" in Gaudi. Therefore: the installed `include` directory needs to
   be in the include paths of Compiler Explorer.
 * With the same argument, they need to get all dependencies of Gaudi.
 * I do not filter inter-project dependencies: so external dependencies of all
   Gaudi packages get added to the path, regardless if the user's package has
   that dependency or not.
 * I do not filter inter-project dependencies: so even public Gaudi headers of
   all Gaudi packages are available, even if the user could not access them in
   real Gaudi (if that would reverse the inter-package dependency)
 * Local (not installed) headers are "hard" to access when the user puts it in
   the their package: They are only accessible within a package but not across
   packages. They are furthermore not available in the `InstallArea` on cvmfs
   (though they are deployed elsewhere in cvmfs, because it's nice to have sources
   available). Since I don't want to add the complexity of packages, I don't
   make them accessible.

Clearly that is not ideal for all users but keeps the complexity of the phase
space of configurations managable (one per project instead of one per package).

And if a user wants additional local headers, they can provide absolute paths
to cvmfs.

In a first version, I had just picked the includes of one random translation
unit and did some replacements in my text editor to reach the Compiler Explorer
syntax. But that wouldn't scale when the nightlies would switch to a new ROOT
or LCG version and has limited copy-and-paste capabilities (different ROOT
version in different nightly slot configurations).

So I figured I need an automatic parser of the `compile_commands.json`, ideally
written in a language that has a json parser, and ended up writing
[**github.com/pseyfert/compilecommands\_to\_compilerexplorer**](https://github.com/pseyfert/compilecommands_to_compilerexplorer)
in go. The decision for go was mostly because I want to learn it and practise
probably helps.

The go application basically goes through all translation units specified in
the json database, gets the command, filters `-I` and `-isystem` out and builds
the union of all include paths of all translation units. Then applies the
filters I want to apply (replace paths in `InstallArea`s with their guessed
target on cvmfs, add the project itself, remove source dirs), according to
heuristics derived from how I read these paths.

I then run a loop over several of these jsons to get a set of libraries and
versions and write out the compiler explorer configuration.

The only slightly tricky bit was specify a library once, list all versions
once, and then list the path *for each* version of a library.

The application is now in a cron job, replacing the Compiler Explorer
configuration once per hour based on what it finds installed on cvmfs.
(Nightly slots are often unavailable because the deployment isn't finished or
the build failed, or the version/slot combination just isn't selected for
deployment - my loop is pretty stupid - so the program is tolerant on not
finding projects and only aborts if not a single project is found. That would
result in an invalid configuration at the moment.)

One last aspect I realised when browsing through twitter is, that the config
file writing should better be atomic (such that Compiler Explorer only ever
sees complete files). I use [renameio-write](https://github.com/google/renameio) for
that purpose.

### Compilers

On my laptop I hardly faced problems with using Compiler Explorer as-is. The
system compiler was decent enough, although we got a few new build errors with
gcc 8, that I didn't yet see in the nightlies back when they were gcc 7 as
latest. Setting up compilers for our stack is a bit of an unexpected nightmare
(setting up `LD_LIBRARY_PATH` and other toolchain parts or the undocumented
`COMPILER_PATH` for clang). The only way to compile stuff is the `lcg-`
wrappers such as

```sh
/cvmfs/lhcb.cern.ch/lib/lhcb/LBSCRIPTS/LBSCRIPTS_v9r2p4/InstallArea/scripts/lcg-g++-7.3.0
```

These are shell scripts that forward all arguments to a compiler that's called
with the right environment setup.

We're working on an updated setup (actually, the latest version of the old
wrappers fails on my VM because the host OS isn't detected correctly). The new
config for my VM is:

```
compilers=&lcgg:&lcgclang
defaultCompiler=lcgg730

group.lcgg.compilers=lcgg493:lcgg720:lcgg820:lcgg620:lcgg730:lcgg710:lcgg810
group.lcgg.groupName=lcg-g++

compiler.lcgg493.name=lcg-g++-4.9.3
compiler.lcgg720.name=lcg-g++-7.2.0
compiler.lcgg820.name=lcg-g++-8.2.0
compiler.lcgg620.name=lcg-g++-6.2.0
compiler.lcgg730.name=lcg-g++-7.3.0
compiler.lcgg710.name=lcg-g++-7.1.0
compiler.lcgg810.name=lcg-g++-8.1.0
compiler.lcgg493.exe=/cvmfs/lhcb.cern.ch/lib/bin/x86_64-centos7/lcg-g++-4.9.3
compiler.lcgg720.exe=/cvmfs/lhcb.cern.ch/lib/bin/x86_64-centos7/lcg-g++-7.2.0
compiler.lcgg820.exe=/cvmfs/lhcb.cern.ch/lib/bin/x86_64-centos7/lcg-g++-8.2.0
compiler.lcgg620.exe=/cvmfs/lhcb.cern.ch/lib/bin/x86_64-centos7/lcg-g++-6.2.0
compiler.lcgg730.exe=/cvmfs/lhcb.cern.ch/lib/bin/x86_64-centos7/lcg-g++-7.3.0
compiler.lcgg710.exe=/cvmfs/lhcb.cern.ch/lib/bin/x86_64-centos7/lcg-g++-7.1.0
compiler.lcgg810.exe=/cvmfs/lhcb.cern.ch/lib/bin/x86_64-centos7/lcg-g++-8.1.0

group.lcgclang.compilers=lcgclang600:hack
group.lcgclang.groupName=lcg-clang++

compiler.lcgclang600.name=lcg-clang++-6.0.0
compiler.lcgclang600.exe=/cvmfs/lhcb.cern.ch/lib/bin/x86_64-centos7/lcg-clang++-6.0.0

compiler.hack.name=hacked clang 6.0.0 with gcc 7.3.0 binutils
compiler.hack.exe=/home/pseyfert/hack
```

The released clang 6 uses gcc 6.2 binutils / headers, so fails with c++17. My
hack uses a newer version that uses gcc7 headers but isn't fully deployed yet,
so I won't go into details. The setup isn't automised at the moment, and the
defaultCompiler is just what I see in the imho most relevant nightly slot.

### Compiler options

Our stack is c++17 and for sure compilations w/o c++11 will fail immediately.
We're also interested in vectorization, which depends on compiler flags, so I
want to provide users with the nightly compiler options by default, so they're
looking at "the right" thing. For `etc/config/c++.defaults.properties` this
means:

```
group.lcgg.options=-DBOOST_FILESYSTEM_VERSION=3 -DBOOST_SPIRIT_USE_PHOENIX_V3 -DGAUDI_LINKER_LIBRARY -DGAUDI_V20_COMPAT -DPACKAGE_NAME="CompilerExplorer" -DPACKAGE_VERSION="v0r0" -D_GNU_SOURCE -Df2cFortran -Dlinux -Dunix -mavx2 -mfma -fmessage-length=0 -pipe -Wall -Wextra -Werror=return-type -pthread -pedantic -Wwrite-strings -Wpointer-arith -Woverloaded-virtual -Wsuggest-override -std=c++17 -fdiagnostics-color -O3 -DNDEBUG -fPIC 
compiler.lcgg493.options=-DBOOST_FILESYSTEM_VERSION=3 -DBOOST_SPIRIT_USE_PHOENIX_V3 -DGAUDI_LINKER_LIBRARY -DGAUDI_V20_COMPAT -DPACKAGE_NAME="CompilerExplorer" -DPACKAGE_VERSION="v0r0" -D_GNU_SOURCE -Df2cFortran -Dlinux -Dunix -mavx2 -mfma -fmessage-length=0 -pipe -Wall -Wextra -Werror=return-type -pthread -pedantic -Wwrite-strings -Wpointer-arith -Woverloaded-virtual -fdiagnostics-color -O3 -DNDEBUG -fPIC 
compiler.hack.options=--gcc-toolchain=/cvmfs/sft.cern.ch/lcg/releases/clang/6.0.0-6647e/x86_64-centos7-gcc62-opt/../../../gcc/7.3.0-90605/x86_64-centos7-gcc62-opt -DBOOST_FILESYSTEM_VERSION=3 -DBOOST_SPIRIT_USE_PHOENIX_V3 -DGAUDI_LINKER_LIBRARY -DGAUDI_V20_COMPAT -DPACKAGE_NAME="CompilerExplorer" -DPACKAGE_VERSION="v0r0" -D_GNU_SOURCE -Df2cFortran -Dlinux -Dunix -mavx2 -mfma -fmessage-length=0 -pipe -Wall -Wextra -Werror=return-type -pthread -pedantic -Wwrite-strings -Wpointer-arith -Woverloaded-virtual -std=c++17 -fdiagnostics-color -O3 -DNDEBUG -fPIC 
```

This is copy and paste from the compiler wrapper and the
`compile_commands.json`. Automisation is yet to come. These settings are the
same for all compilers (except where they don't make sense - clang needs
`-Wsuggest-override` removed because it doesn't exist. The old gcc doesn't know
c++17). An issue with this is, compiler options depend on the nightly slot
(what I use as library version) and not on compilers. So I intend to provide
the user with several different `g++ 7`, pointing to the same binary, but
configured with different options. The options shall then come from the
`compile_commands.json`. The draft is in the parser on github, but it's a bit
more tricky because these options change from translation unit to translation
unit: dedicated settings by a package developer and we differ between libraries
and modules through `-DGAUDI_LINKER_LIBRARY`, so I doubt a union of all options
is a good idea.

Another trouble is character escaping. As you can see above, the
`compile_commands.json` syntax for `-D` defined strings is

```json
-DPACKAGE_NAME=\\\"BrunelCache\\\"
```

and Compiler Explorer needs it as

```
-DPACKAGE_NAME="CompilerExplorer"
```

my draft for option parsing
[here](https://github.com/pseyfert/compilecommands_to_compilerexplorer/blob/d6651a8f9e2bd65c28ca8e09e9e0dd1d38928b32/cc2ce/in.go#L87)
removes these three-backslash-single-double-quote and replaces them by a single
double-quote (twice per definition). Although the need to escape these in the
go source doesn't lead to pretty code.

### Formatting

We have our own in house clang-format style. By default Compiler Explorer only
offers the clang-format builtin styles, but the drop down menu can be extended
in `static/settings.js` to

```js
    var formats = ["file", "Google", "LLVM", "Mozilla", "Chromium", "WebKit"];
```

The only question is, how does clang-format then find our format style? I put
it in the `compiler-explorer` main directory and that seems to work.

The formatting result also depends on the clang-format version, so I use the
one we deployed for the gitlab continous integration, in `etc/config/compiler-explorer.defaults.properties`

```
formatter.clangformat.exe=/cvmfs/lhcb.cern.ch/lib/bin/x86_64-centos7/lcg-clang-format-3.9
```

And then I disabled other styles in `etc/config/compiler-explorer.defaults.properties` with

```
formatter.clangformat.styles=file
```

Users will get an error message if they try other styles (I could also remove
them from the menu above …).

### other tweaks

Building with range v3, my compilations often timed out, so the timeout in
`etc/config/compiler-explorer.defaults.properties` needs to be increased to

```
compileTimeoutMs=90000
```

### server setup

#### network and security

As Matt Godbolt said in his talk at CppCon, compilers are a big security
nightmare. Basically my Compiler Explorer instance's website grants remote code
execution to everybody. Including writing to the file system. Since I didn't
even want to hang that behind the CERN firewall into the network for everybody
to use, or even allow individual users to access my CERN account's home dir, I
set up a VM that does not use the CERN user database but has "normal" local
users. I also set up iptables to only accept connections from localhost, so
only users to whom I gave a login can access Compiler Explorer. At this point
the website hardly grants them more rights than they have already (they can
just compile on the command line).

```
# /etc/iptables/iptables.rules
...
-A INPUT -p tcp -m tcp --dport 10240 -j DROP
```

The machine is a CernVM (strangely running Scientific Linux 7 and not CentOS 7,
but so be it) where iptables is not enabled by default, so

```sh
systemctl enable iptables
systemctl start iptables
```

Users can now use the website by opening an ssh tunnel (due to firewall
settings only from within the CERN network):

```sh
ssh -L <localport>:127.0.0.1:10240 <username>@pseyfert-ce
```

and then access `localhost:<localport>` in their local webbrowser.

#### Compiler Explorer process

On my local laptop I just start Compiler Explorer with `make` in a shell. For a
running server that's not cool (when I log out, the site goes down). A first
bodge was to run `make` in a screen session, but that still requires manual
restarting after a reboot. I was told systemd is what to cool kids are doing.
Originally I wanted to run Compiler Explorer as a user (not as root), so I
looked up how users can run systemd. The problem is, that feature is not
available on the redhat version SL7 is based on, but I found that root can
specify in system systemd files as which user a process should be running:

```
[Unit]
Description=Compiler Explorer
# dependencies are a bit guessed at this point, based on units I saw on the system
After=autofs.service network-online.target network.target cernvm.service
Wants=autofs.service network-online.target cernvm.service

[Service]
WorkingDirectory=/home/pseyfert/compiler-explorer
ExecStart=/usr/bin/make
# run as `pseyfert`
User=pseyfert

[Install]
WantedBy=default.target
```

The configuration generation remains a cron job of the user for now.

### other annoyances

Scientific Linux 7 comes with a too-old node.js version. Had to follow the
node.js instructions how to install a more recent version. Fingers crossed the
next auto-update won't destroy it.
