// Copyright (c) 2020 Contaim, LLC
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package spec

import (
	"path"
	"path/filepath"

	"github.com/hashicorp/hcl/v2"
)

// diagnostic messages
const (
	DiagCannotDetermineFileType       = "Cannot determine file type based on extension, only .json and .hcl files are supported"
	DiagCannotDetermineFileTypeDetail = "You must provide a file with either a .json or .hcl extension, only json and hcl files are supported"

	DiagGlobError       = "There was a problem parsing the file pattern"
	DiagGlobErrorDetail = "The file pattern was not able to be parsed. This might be an implementation problem."
)

// Files accepts many file paths and processes each. All files provided will be processed
// against the same Spec so all files should be of the same type. If not, the diagnostics
// will return errors for things not expected by the current Spec.
func (s *Spec) Files(filenames ...string) *Diagnostics {
	diags := hcl.Diagnostics{}

	for _, filename := range filenames {
		switch ext := path.Ext(filename); ext {
		case ".json":
			diags = diags.Extend(s.parseJSONFile(filename))
		case ".hcl":
			diags = diags.Extend(s.parseHCLFile(filename))
		default:
			diags = diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  DiagCannotDetermineFileType,
				Detail:   DiagCannotDetermineFileTypeDetail,
			})
		}
	}

	return newDiagnostics(s, diags)
}

// ParsedFiles returns all of the filenames that we've parsed through various parsing
// methods.
func (s *Spec) ParsedFiles() []string {
	files := []string{}

	for filename := range s.files {
		files = append(files, filename)
	}

	return files
}

// FileGlob works the same as Files but instead builds a list of filenames to process
// using the provided glob pattern.
func (s *Spec) FileGlob(pattern string) *Diagnostics {
	filenames, _ := filepath.Glob(pattern)

	// TODO: revisit error handling here
	// A note on not handling the error:
	//   I'm not happy about this (or maybe I am?), but I could not generate
	//   a test to break this with error checking, including an old GitHub
	//   issue that caused this to recurse before (https://github.com/golang/go/issues/15879)
	//   due to the OS itself. I am going to leave this without but if we notice
	//   issues we can bring in error handling.

	return s.Files(filenames...)
}
