/*
 *   Copyright (c) 2020
 *   All rights reserved.
 */
// Copyright (c) 2020 Contaim, LLC
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package spec_test

import (
	"bytes"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/responserms/spec"
	"github.com/stretchr/testify/assert"
)

func TestDiagnostics(tt *testing.T) {
	diags := &spec.Diagnostics{
		Spec: spec.NewSubset(),
		Diags: hcl.Diagnostics{
			{
				Severity: hcl.DiagError,
				Summary:  "This is a fake error.",
				Detail:   "This was a fake error.",
			},
		},
	}

	tt.Run("Raw() returns the raw hcl.Diagnostics", func(t *testing.T) {
		assert.IsType(t, hcl.Diagnostics{}, diags.Raw())
	})

	tt.Run("HasErrors() returns true when errors are in diagnostics", func(t *testing.T) {
		assert.True(t, diags.HasErrors())
	})

	tt.Run("Error() returns a string with the error message", func(t *testing.T) {
		assert.Contains(t, diags.Error(), "This is a fake error.")
	})

	tt.Run("Errs() returns a slice of error instances", func(t *testing.T) {
		assert.Len(t, diags.Errs(), 1)
	})

	// we don't thoroughly test this since it is actually passed directly into
	// the HCL language handler which does have more test coverage.
	tt.Run("WriteText() writes diagnostics to the `to` io.Writer", func(t *testing.T) {
		b := new(bytes.Buffer)
		if err := diags.WriteText(b, 0, false); err != nil {
			t.Fatalf("diags.WriteText() returned an error: %s", err)
		}

		assert.Contains(t, b.String(), "This is a fake error.")
	})
}
