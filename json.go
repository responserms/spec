// Copyright (c) 2020 Contaim, LLC
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package spec

import (
	"fmt"
	"io/ioutil"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/json"
)

// ParseJSON parses the raw src as JSON.
func (s *Spec) ParseJSON(src []byte, filename string) *Diagnostics {
	return newDiagnostics(s, s.parseJSON(src, filename))
}

func (s *Spec) parseJSON(src []byte, filename string) hcl.Diagnostics {
	file, diags := json.Parse(src, filename)
	s.files[filename] = file

	return diags
}

// ParseJSONFile parses a single JSON file by reading it from the filesystem.
func (s *Spec) ParseJSONFile(filename string) *Diagnostics {
	return newDiagnostics(s, s.parseJSONFile(filename))
}

func (s *Spec) parseJSONFile(filename string) hcl.Diagnostics {
	src, err := ioutil.ReadFile(filename)
	if err != nil {
		return hcl.Diagnostics{
			{
				Severity: hcl.DiagError,
				Summary:  "Failed to read file",
				Detail:   fmt.Sprintf("The JSON file %q could not be read.", filename),
			},
		}
	}

	return s.parseJSON(src, filename)
}
