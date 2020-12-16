// Copyright (c) 2020 Contaim, LLC
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parser

import (
	"sort"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
)

// verify that our public APIs remain
var _ BlockRegistrar = (*Registrar)(nil)
var _ Builder = (*Registrar)(nil)
var _ Parser = (*Registrar)(nil)

// BlockRegistrar registers blocks by adding the given BlockDefinition for the root blockName.
type BlockRegistrar interface {
	RegisterBlock(blockName string, blockDef BlockDefinition)
	AddRegistration(reg *Registration)
}

// Builder builds and returns an hcldec.Spec used for further processing.
type Builder interface {
	Build() hcldec.Spec
}

// Parser accepts an hcl.Body and a pointer to an hcl.EvalContext, parses the hcl.Body against the
// hcl.EvalContext and returns any returned hcl.Diagnostics.
type Parser interface {
	Parse(body hcl.Body, ctx *hcl.EvalContext) hcl.Diagnostics
}

// Registration adds a BlockDefinition for the BlockName to the Registrar that will allow parsing
// the block in the schema according to the returned hcldec.Spec from BlockDefinition's Spec()
// method.
type Registration struct {
	BlockName  string
	Order      int
	Definition BlockDefinition
}

// Registrar allows registering many root hcldec.Spec's at a given block or argument which
// are all called independently allowing variables to be injected into the context after each
// run for variable interpolation that includes variables created from previous schema blocks
// or arguments.
//
// The registrar also allows injecting functions after specific hcldec.Spec's are processed,
// though in general this should be avoided.
type Registrar struct {
	NextOrder           int
	IncreaseNextOrderBy int
	registrations       []*Registration
}

// NewRegistrar creates a new Registrar instance, returning the pointer to be used for registering
// DefinitionBlock's.
func NewRegistrar(increaseNextOrderBy int) *Registrar {
	return &Registrar{
		NextOrder:           0,
		IncreaseNextOrderBy: increaseNextOrderBy,
		registrations:       make([]*Registration, 0),
	}
}

// Registrations returns pointers to all of the registered registrations. These may be modified
// as needed before calling Build() which will then convert them to a full hcl.Spec.
func (r *Registrar) Registrations() []*Registration {
	return r.registrations
}

// RegisterBlock registers the BlockDefinition with the given blockName, automatically ordering the
// schema to be processed at the NextOrder in the Registrar, then increasing the NextOrder by the
// Registrar's configured IncreaseNextOrderBy.
func (r *Registrar) RegisterBlock(blockName string, blockDef BlockDefinition) {
	reg := &Registration{
		BlockName:  blockName,
		Order:      r.NextOrder,
		Definition: blockDef,
	}

	// increase NextOrder by the value of IncreaseNextOrderBy
	r.NextOrder += r.IncreaseNextOrderBy
	r.AddRegistration(reg)
}

// AddRegistration adds a raw Registration to the slice of registered specs. When
// using this method directly you will be responsible for assigning the appropriate
// order to the registration. If the order is not placed before any blocks that might
// depend on variable scopes, variables will remain unknown and ultimately unvalidated.
//
// Unknown variables may lead to an invalid Result from a Schema.
func (r *Registrar) AddRegistration(reg *Registration) {
	r.registrations = append(r.registrations, reg)
}

// Build generates a full hcldec.Spec out of the registered BlockDefinition's allowing
// the full schema, with all blocks, to be used in a HCL parser.
func (r *Registrar) Build() hcldec.Spec {
	res := hcldec.ObjectSpec{}

	for _, reg := range r.registrations {
		res[reg.BlockName] = reg.Definition.Spec()
	}

	return res
}

// Parse handles parsing the hcl.Body against the dynamically generated hcldec.Spec from
// the registered BlockDefinition's.
func (r *Registrar) Parse(body hcl.Body, ctx *hcl.EvalContext) hcl.Diagnostics {
	ordered := r.registrations

	if ctx.Functions == nil {
		ctx.Functions = map[string]function.Function{}
	}

	if ctx.Variables == nil {
		ctx.Variables = map[string]cty.Value{}
	}

	sort.Slice(
		ordered,
		func(i, j int) bool {
			return ordered[i].Order < ordered[j].Order
		},
	)

	var lastBody = body
	var lastDiags = hcl.Diagnostics{}

	for _, reg := range ordered {
		val, body, diags := hcldec.PartialDecode(lastBody, reg.Definition.Spec(), ctx)
		lastDiags = lastDiags.Extend(diags)
		lastBody = body

		// if a FunctionInjector
		if inj, ok := reg.Definition.(FunctionInjector); ok {
			for k, v := range inj.Functions(val) {
				ctx.Functions[k] = v
			}
		}

		// if a VariableInjector
		if inj, ok := reg.Definition.(VariableInjector); ok {
			for k, v := range inj.Variables(val) {
				ctx.Variables[k] = v
			}
		}
	}

	return lastDiags
}
