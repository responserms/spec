// Copyright (c) 2020 Contaim, LLC
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package spec

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/stretchr/testify/assert"
)

func TestInternalDiagnostics(tt *testing.T) {
	tt.Run("newDiagnostics creates a new intance of Diagnostics", func(t *testing.T) {
		spec := &Spec{}
		diag := newDiagnostics(spec, hcl.Diagnostics{})

		assert.IsType(t, &Diagnostics{}, diag)
	})
}
