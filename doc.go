// Copyright (c) 2020 Contaim, LLC
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package spec provides schema specification parsing for Response's HCL/JSON-based
// configuration and data files. This package is a convenience wrapper around the
// hcl/v2 package itself and adds helpers for dealing with our specifications
// directly. This is likely not useful if you do not need to parse one of our schema
// files or something similar.
package spec
