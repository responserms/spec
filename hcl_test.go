// Copyright (c) 2020 Contaim, LLC
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package spec_test

import (
	"testing"

	"github.com/responserms/spec"
	"github.com/stretchr/testify/assert"
)

func TestParseHCL(tt *testing.T) {
	var rawHCL = []byte(`vars {}`)

	tt.Run("test parsing", func(t *testing.T) {
		subset := spec.NewSubset()
		diags := subset.ParseHCL(rawHCL, "somefile.hcl")

		assert.False(tt, diags.HasErrors())
	})
}

func TestParseHCLFile(tt *testing.T) {
	tt.Run("file parsing returns no errors when the file exists", func(t *testing.T) {
		subset := spec.NewSubset()
		diags := subset.ParseHCLFile("./testdata/test.hcl")

		assert.False(t, diags.HasErrors())
	})

	tt.Run("file parsing returns diagnostics when the file is missing", func(t *testing.T) {
		subset := spec.NewSubset()
		diags := subset.ParseHCLFile("./testdata/does_not_exist.hcl")

		assert.True(t, diags.HasErrors())
	})
}
