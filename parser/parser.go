// Copyright (c) 2020 Contaim, LLC
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parser

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
)

// InjectableVariables is a map of variable names that correspond with a
// cty.Value of which will be injected into the running context after processing
// Spec()
type InjectableVariables map[string]cty.Value

// InjectableFunctions is a map of function names that correspond with a
// cty.Function of which will be injected into the context after processing
// Spec()
type InjectableFunctions map[string]function.Function

// NamedBlockDefinitions represents a slice of individual NamedBlockDefinition's pre-ordered
// in the order they should be processed.
type NamedBlockDefinitions []NamedBlockDefinition

// BlockDefinition describes the methods that must be implemented by all
// definition packages for each root-level key.
type BlockDefinition interface {

	// Spec must return a hcldec.Spec instance. This will be used to parse from
	// the root of the schema file.
	Spec() hcldec.Spec
}

// NamedBlockDefinition describes the methods that musy be implemented by all
// definition packages for each root-level key when the definition can be self-registered
// by it's block name.
type NamedBlockDefinition interface {
	BlockDefinition

	// Name must return the name of the block to be registered.
	Name() string
}

// VariableInjector allows injecting variables into a hcl.EvalContext after it has been
// created from the BlockDefinition's Spec() result.
type VariableInjector interface {
	BlockDefinition

	// Variables must return an instance of InjectableVariables. These variables
	// will be passed through all BlockDefinition's that follow it allowing values
	// created from this block to be available for all that follow.
	Variables(v cty.Value) InjectableVariables
}

// FunctionInjector allows injecting functions into a hcl.EvalContext after it has been
// created from the BlockDefinition's Spec() result.
type FunctionInjector interface {
	BlockDefinition

	// Functions must return an instance of InjectableFunctions. These functions
	// will be passed through all BlockDefinition's that follow it allowing functions
	// created from this block to be available for all that follow.
	Functions(v cty.Value) InjectableFunctions
}
