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

func colonSeparateArray(stringset []string) string {
	var retval string
	for i, s := range stringset {
		if i != 0 {
			retval += ":"
		}
		retval += s
	}
	return retval
}

func colonSeparateMap(stringset map[string]bool) string {
	var retval string
	addseparator := false
	for k, _ := range stringset {
		if addseparator {
			retval += ":"
		} else {
			addseparator = true
		}
		retval += k
	}
	return retval
}
