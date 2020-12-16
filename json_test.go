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

func TestParseJSON(tt *testing.T) {
	var rawJSON = []byte(`{}`)

	tt.Run("test parsing", func(t *testing.T) {
		subset := spec.NewSubset()
		diags := subset.ParseJSON(rawJSON, "somefile.json")

		assert.False(tt, diags.HasErrors())
	})
}

func TestParseJSONFile(tt *testing.T) {
	tt.Run("file parsing returns no errors when the file exists", func(t *testing.T) {
		subset := spec.NewSubset()
		diags := subset.ParseJSONFile("./testdata/test.json")

		assert.False(t, diags.HasErrors())
	})

	tt.Run("file parsing returns diagnostics when the file is missing", func(t *testing.T) {
		subset := spec.NewSubset()
		diags := subset.ParseJSONFile("./testdata/does_not_exist.json")

		assert.True(t, diags.HasErrors())
	})
}
