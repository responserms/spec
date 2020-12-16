/*
 *   Copyright (c) 2020
 *   All rights reserved.
 */
// Copyright (c) 2020 Contaim, LLC
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package spec

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/responserms/spec/parser"
)

type specFiles map[string]*hcl.File

// Spec represents a single type of HCL/JSON schema variant and provides helpers to easily
// parse raw bytes, files, and more against the schema. The Spec returns a custom Diagnostics
// rather than the hcl.Diagnostics allowing easy manipulation of our own errors.
type Spec struct {
	registrar *parser.Registrar
	files     specFiles
}

// New creates a new Spec instance with the pre-ordered slice of parser.NamedBlockDefiniion
// instances provided. This is the typical API where you will use the Schema variable from
// a particular schema.
func New(defs parser.NamedBlockDefinitions) *Spec {
	registrar := parser.NewRegistrar(1)

	for _, def := range defs {
		registrar.RegisterBlock(def.Name(), def)
	}

	return &Spec{
		registrar: registrar,
		files:     specFiles{},
	}
}

// NewSubset creates a new Spec instance with one or more parser.NamedBlockDefinition
// instances provided. This is primarily useful for validating only a specific block or
// for performing tests.
func NewSubset(defs ...parser.NamedBlockDefinition) *Spec {
	registrar := parser.NewRegistrar(1)

	for _, def := range defs {
		registrar.RegisterBlock(def.Name(), def)
	}

	return &Spec{
		registrar: registrar,
		files:     specFiles{},
	}
}

// Registrations returns all of the parser.Registration instances that are added to the
// registrar.
func (s *Spec) Registrations() []*parser.Registration {
	return s.registrar.Registrations()
}

// Build builds an hcldec.Spec instance from all registered BlockDefinition's.
func (s *Spec) Build() hcldec.Spec {
	return s.registrar.Build()
}

// Body returns an hcl.Body that merges all processed files into a single body for further
// processing.
func (s *Spec) Body() hcl.Body {
	files := []*hcl.File{}

	for _, file := range s.files {
		files = append(files, file)
	}

	return hcl.MergeFiles(files)
}

// Parse parses the provided hcl.Body, given the hcl.EvalContext against the generated
// hcldec.Spec and ordered according to the order that the BlockDefinition's were defined.
func (s *Spec) Parse(ctx *hcl.EvalContext) *Diagnostics {
	return newDiagnostics(s, s.registrar.Parse(s.Body(), ctx))
}

// Decode extracts the configuration within the given body into the given value. This value must
// be a non-nil pointer to either a struct or a map, where in the former case the configuration
// will be decoded using struct tags and in the latter case only attributes are allowed and their
// values are decoded into the map.
//
// The given hcl.EvalContext is used to resolve any variables or functions in expressions encountered while
// decoding. This may be nil to require only constant values, for simple applications that do not support
// variables or functions.
//
// The returned diagnostics should be inspected with its HasErrors method to determine if the populated
// value is valid and complete. If error diagnostics are returned then the given value may have been
// partially-populated but may still be accessed by a careful caller for static analysis and editor
// integration use-cases.
func (s *Spec) Decode(ctx *hcl.EvalContext, val interface{}) *Diagnostics {
	return newDiagnostics(s, gohcl.DecodeBody(s.Body(), ctx, val))
}
