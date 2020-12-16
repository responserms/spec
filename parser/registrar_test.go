// Copyright (c) 2020 Contaim, LLC
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parser_test

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/responserms/spec/parser"
	"github.com/stretchr/testify/assert"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
)

type testSpecBlockDef struct{}

func (t *testSpecBlockDef) Spec() hcldec.Spec {
	return &hcldec.AttrSpec{
		Name:     "test",
		Type:     cty.String,
		Required: true,
	}
}

func TestNewRegistrar(t *testing.T) {
	t.Run("properly sets defaults", func(t *testing.T) {
		reg := parser.NewRegistrar(2)

		assert.Equal(t, 0, reg.NextOrder)
		assert.Equal(t, 2, reg.IncreaseNextOrderBy)
	})
}

func TestRegistrar(t *testing.T) {
	t.Run("registrar allows registering raw blocks", func(t *testing.T) {
		registrar := parser.NewRegistrar(1)
		registration := &parser.Registration{
			BlockName:  "test",
			Order:      100,
			Definition: &testSpecBlockDef{},
		}

		registrar.AddRegistration(registration)

		assert.IsType(t, []*parser.Registration{}, registrar.Registrations())
		assert.Len(t, registrar.Registrations(), 1)
	})

	t.Run("RegisterBlock adds proper order", func(t *testing.T) {
		reg := parser.NewRegistrar(1)

		// change the order
		reg.IncreaseNextOrderBy = 100

		// register the block
		reg.RegisterBlock("test", &testSpecBlockDef{})
		reg.RegisterBlock("another", &testSpecBlockDef{})

		assert.Equal(t, 100, reg.Registrations()[1].Order)
	})
}

func TestBuild(t *testing.T) {
	t.Run("ensure Build() adds all attributes in hcldec.ObjectSpec before returning", func(t *testing.T) {
		reg := parser.NewRegistrar(1)
		def := &testSpecBlockDef{}
		reg.RegisterBlock("test", def)

		spec := reg.Build()

		assert.IsType(t, hcldec.ObjectSpec{}, spec)
		assert.IsType(t, def.Spec(), spec.(hcldec.ObjectSpec)["test"])
	})
}

var execOrder = 1
var _ parser.BlockDefinition = (*orderTrackingDefSpec)(nil)

type orderTrackingDefSpec struct {
	order int
}

func (s *orderTrackingDefSpec) Spec() hcldec.Spec {
	s.order = execOrder
	execOrder++

	return hcldec.ObjectSpec{}
}
func (s *orderTrackingDefSpec) Order() int {
	return s.order
}

var _ parser.FunctionInjector = (*funcDefSpec)(nil)

type funcDefSpec struct{}

func (s *funcDefSpec) Spec() hcldec.Spec {
	return &hcldec.ObjectSpec{}
}

func (s *funcDefSpec) Functions(v cty.Value) parser.InjectableFunctions {
	return map[string]function.Function{
		"noop": function.New(&function.Spec{
			Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
				return cty.StringVal("noop"), nil
			},
		}),
	}
}

var _ parser.VariableInjector = (*varDefSpec)(nil)

type varDefSpec struct{}

func (s *varDefSpec) Spec() hcldec.Spec {
	return &hcldec.ObjectSpec{}
}

func (s *varDefSpec) Variables(v cty.Value) parser.InjectableVariables {
	return map[string]cty.Value{
		"one": cty.NumberIntVal(1),
	}
}

func TestParse(t *testing.T) {
	t.Run("ensure Parse() initializes nil maps", func(t *testing.T) {
		reg := parser.NewRegistrar(1)
		ctx := &hcl.EvalContext{
			Functions: nil,
			Variables: nil,
		}

		reg.Parse(nil, ctx)

		assert.NotNil(t, ctx.Functions)
		assert.NotNil(t, ctx.Variables)
	})

	t.Run("ensure blocks are ordered properly", func(t *testing.T) {
		reg := parser.NewRegistrar(1)

		last := &orderTrackingDefSpec{3}
		reg.AddRegistration(&parser.Registration{
			BlockName:  "last",
			Order:      100,
			Definition: last,
		})

		second := &orderTrackingDefSpec{2}
		reg.AddRegistration(&parser.Registration{
			BlockName:  "second",
			Order:      50,
			Definition: second,
		})

		first := &orderTrackingDefSpec{1}
		reg.AddRegistration(&parser.Registration{
			BlockName:  "first",
			Order:      1,
			Definition: first,
		})

		reg.Parse(hcl.EmptyBody(), &hcl.EvalContext{})

		assert.Equal(t, 1, first.Order())
		assert.Equal(t, 2, second.Order())
		assert.Equal(t, 3, last.Order())
	})

	t.Run("ensure functions are added to EvalContext after each Spec()", func(t *testing.T) {
		reg := parser.NewRegistrar(1)
		reg.RegisterBlock("funcs", &funcDefSpec{})
		reg.RegisterBlock("vars", &varDefSpec{})

		ctx := &hcl.EvalContext{}
		reg.Parse(hcl.EmptyBody(), ctx)

		assert.IsType(t, function.Function{}, ctx.Functions["noop"])
		assert.IsType(t, cty.Value{}, ctx.Variables["one"])
	})
}
