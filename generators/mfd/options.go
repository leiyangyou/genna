package mfd

import (
	"github.com/dizzyfool/genna/generators/base"
)

// Options stores generator options
type Options struct {
	base.Options

	// Package sets package name for model
	Package string
}
