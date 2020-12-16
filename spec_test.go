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
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/responserms/spec"
	"github.com/responserms/spec/parser"
	"github.com/stretchr/testify/assert"
	"github.com/zclconf/go-cty/cty"
)

type testSchema struct{}

func (t *testSchema) Name() string {
	return "test"
}

func (t *testSchema) Spec() hcldec.Spec {
	return &hcldec.BlockAttrsSpec{
		TypeName:    "test",
		ElementType: cty.String,
		Required:    true,
	}
}

var schemaList = parser.NamedBlockDefinitions{
	&testSchema{},
}

func TestNew(tt *testing.T) {
	tt.Run("spec.New registers all named block definitions", func(t *testing.T) {
		s := spec.New(schemaList)

		assert.Len(t, s.Registrations(), len(schemaList))
		assert.IsType(t, &spec.Spec{}, s)
	})
}

func TestNewSubset(tt *testing.T) {
	tt.Run("spec.NewSubset allows 0 definitions", func(t *testing.T) {
		s := spec.NewSubset()

		assert.IsType(t, &spec.Spec{}, s)
	})

	tt.Run("spec.NewSubset registers the blocks defined", func(t *testing.T) {
		s := spec.NewSubset(&testSchema{})

		assert.Len(t, s.Registrations(), 1)
		assert.Equal(t, "test", s.Registrations()[0].BlockName)
	})
}

func TestBuild(tt *testing.T) {
	tt.Run("Build() returns the raw hcldec.Spec", func(t *testing.T) {
		s := spec.NewSubset()

		assert.Implements(t, (*hcldec.Spec)(nil), s.Build())
	})
}

func TestBody(tt *testing.T) {
	tt.Run("Body() returns a body containing 0 files when no files are parsed", func(t *testing.T) {
		s := spec.NewSubset()
		fileDiags := s.FileGlob("./testdata/glob/*.hcl")
		assert.False(t, fileDiags.HasErrors())

		body := s.Body()
		attrs, _ := body.JustAttributes()

		// these are specified in individual files within ./testdata/glob/*.hcl
		assert.Equal(t, "one", attrs["one"].Name)
		assert.Equal(t, "two", attrs["two"].Name)
		assert.Equal(t, "three", attrs["three"].Name)
		assert.Equal(t, "four", attrs["four"].Name)
	})
}

func TestParse(tt *testing.T) {
	tt.Run("Parse() returns a custom Diagnostics object", func(t *testing.T) {
		s := spec.NewSubset()
		diags := s.Parse(&hcl.EvalContext{})

		assert.IsType(t, &spec.Diagnostics{}, diags)
	})
}

type decodeStruct struct {
	One   string `hcl:"one,attr"`
	Two   string `hcl:"two,attr"`
	Three string `hcl:"three,attr"`
	Four  string `hcl:"four,attr"`
}

func TestDecode(tt *testing.T) {
	tt.Run("Decode() properly decodes the body into a struct with hcl tags", func(t *testing.T) {
		s := spec.NewSubset()
		s.FileGlob("./testdata/glob/*.hcl")

		decode := &decodeStruct{}
		diags := s.Decode(&hcl.EvalContext{}, decode)

		assert.False(t, diags.HasErrors())
		assert.Equal(t, "one", decode.One)
		assert.Equal(t, "two", decode.Two)
		assert.Equal(t, "three", decode.Three)
		assert.Equal(t, "four", decode.Four)
	})
}
