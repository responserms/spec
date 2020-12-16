// Copyright (c) 2020 Contaim, LLC
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package spec

import (
	"io"

	"github.com/hashicorp/hcl/v2"
)

// Diagnostics is used to represent a number of diagnostics returned from various places
// in Spec parsing and file reading. Diagnostics is general purpose in nature and is meant
// to be displayed to the end-user and read by machine.
type Diagnostics struct {
	Spec  *Spec
	Diags hcl.Diagnostics
}

// newDiagnostics creaes a new Diagnostics instance.
func newDiagnostics(spec *Spec, diags hcl.Diagnostics) *Diagnostics {
	return &Diagnostics{
		Spec:  spec,
		Diags: diags,
	}
}

// Raw returns the raw hcl.Diagnostics instance. This is useful if you need to interact
// with the raw implementation rather than the sugared version we provide. In most cases
// you won't need this.
func (d *Diagnostics) Raw() hcl.Diagnostics {
	return d.Diags
}

// HasErrors simply returns true if there are errors within this Diagnostics instance. If
// no errors this will return false.
func (d *Diagnostics) HasErrors() bool {
	return d.Diags.HasErrors()
}

// Error returns all of the diagnostics coerced into a string meant to be shown to the end-user.
// In general this should likely be avoided and you should use WriteText instead.
func (d *Diagnostics) Error() string {
	return d.Diags.Error()
}

// Errs returns all of the native Go error interfaces for the diagnostics returned.
func (d *Diagnostics) Errs() []error {
	return d.Diags.Errs()
}

// WriteText writes the output in a format easily understood by humans to the provided io.Writer. This
// is useful for showing the output of the diagnostics to end-users in a CLI environment but less useful
// when the caller is an application.
//
// Setting the width to 0 disables word wrapping. Setting the color to false will disable coloring of
// key information in the output. The output will contain relevant context such as line numbers and code
// snippets.
func (d *Diagnostics) WriteText(to io.Writer, width uint, color bool) error {
	wr := hcl.NewDiagnosticTextWriter(to, d.Spec.files, width, color)
	return wr.WriteDiagnostics(d.Diags)
}
