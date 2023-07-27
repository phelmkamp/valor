// Copyright 2023 phelmkamp. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package unit provides an empty, unit type.
package unit

// Type is the unit type.
type Type = struct{}

// Unit is the value of the unit type.
//
// It's empty and occupies no space in memory.
var Unit = Type{}
