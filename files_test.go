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

func TestFiles(tt *testing.T) {
	subset := spec.NewSubset()

	tt.Run("a valid hcl file returns no errors", func(t *testing.T) {
		hclFile := subset.Files("./testdata/test.hcl")
		assert.False(t, hclFile.HasErrors())
	})

	tt.Run("a valid json file returns no errors", func(t *testing.T) {
		jsonFile := subset.Files("./testdata/test.json")
		assert.False(t, jsonFile.HasErrors())
	})

	tt.Run("an invalid filetype results in a diagnostic error", func(t *testing.T) {
		invalidFile := subset.Files("./testdata/thisdoesnotexist.responsermsfile")

		assert.True(t, invalidFile.HasErrors())
		assert.Contains(t, invalidFile.Error(), spec.DiagCannotDetermineFileType)
	})
}

func TestFileGlob(tt *testing.T) {
	subset := spec.NewSubset()

	tt.Run("finds and parses all files we expect", func(t *testing.T) {
		diags := subset.FileGlob("./testdata/glob/*.hcl")
		assert.False(t, diags.HasErrors())
		assert.Len(t, subset.ParsedFiles(), 4)
	})
}
