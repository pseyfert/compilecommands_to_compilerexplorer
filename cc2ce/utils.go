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
	"bytes"
	"strings"
)

// Convert a quasi-set of strings (a map[string]bool) into a colon separated string.
// Only the map keys are considered, values are ignored.
// This is equivalent to strings.Join(stringset, ":").
func ColonSeparateArray(stringset []string) string {
	return strings.Join(stringset, ":")
}

// Convert a quasi-set of strings (a map[string]bool) into a colon separated string.
// Only the map keys are considered, values are ignored.
func ColonSeparateMap(stringset map[string]bool) string {
	var b bytes.Buffer
	addseparator := false
	for k, _ := range stringset {
		if addseparator {
			b.WriteString(":")
		} else {
			addseparator = true
		}
		b.WriteString(k)
	}
	return b.String()
}
