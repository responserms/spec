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
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

// ParseHCL parses the raw src as HCL.
func (s *Spec) ParseHCL(src []byte, filename string) *Diagnostics {
	return newDiagnostics(s, s.parseHCL(src, filename))
}

func (s *Spec) parseHCL(src []byte, filename string) hcl.Diagnostics {
	file, diags := hclsyntax.ParseConfig(src, filename, hcl.Pos{Byte: 0, Line: 1, Column: 1})
	s.files[filename] = file

	return diags
}

// ParseHCLFile parses a single HCL file by reading it from the filesystem.
func (s *Spec) ParseHCLFile(filename string) *Diagnostics {
	return newDiagnostics(s, s.parseHCLFile(filename))
}

func (s *Spec) parseHCLFile(filename string) hcl.Diagnostics {
	src, err := ioutil.ReadFile(filename)
	if err != nil {
		return hcl.Diagnostics{
			{
				Severity: hcl.DiagError,
				Summary:  "Failed to read file",
				Detail:   fmt.Sprintf("The HCL file %q could not be read.", filename),
			},
		}
	}

	return s.parseHCL(src, filename)
}
