# compilecommands 2 compilerexplorer 

[![Licence: GPL v3](https://img.shields.io/github/license/pseyfert/compilecommands_to_compilerexplorer.svg)](LICENSE)
[![travis Status](https://travis-ci.org/pseyfert/compilecommands_to_compilerexplorer.svg?branch=master)](https://travis-ci.org/pseyfert/compilecommands_to_compilerexplorer)

Helpers to create a [compiler explorer](https://godbolt.org/) configuration
from a `compile_commands.json` file.

Originally created to have the LHCb Experiment's nightly software builds
available in Compiler Explorer with transitive include dependencies.

Sorting out LHCb specific parts into dedicated packages is work in progress to
facilitate a more general usage.

Feedback (suggestions, wishes, reviews, bug reports, patches, pull requests,
improvements) is welcome, even if I don't find the time to follow them up.
Please consider that general purpose refactoring is not on my employer's
priorities.

I (very quickly) wrote down a [blog
post](https://pseyfert.web.cern.ch/pseyfert/blog/compiler-explorer-for-lhcb.html)
about my Compiler Explorer setup, its raw version can be found in the `doc`
directory. The blog post is under CC-BY-SA 4.0.
