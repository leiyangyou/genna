package mfd

import (
	"encoding/xml"
	"io"
	"path"

	"github.com/dizzyfool/genna/generators/base"
	"github.com/dizzyfool/genna/lib"
	"github.com/dizzyfool/genna/util"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"golang.org/x/xerrors"
)

const (
	pkg = "pkg"
)

// CreateCommand creates generator command
func CreateCommand(logger *zap.Logger) *cobra.Command {
	return base.CreateCommand("mfd", "XML generator for mdf project", New(logger))
}

// Generator represents mfd generator
type Generator struct {
	logger  *zap.Logger
	options Options
}

// New creates basic generator
func New(logger *zap.Logger) *Generator {
	return &Generator{
		logger: logger,
	}
}

// Logger gets logger
func (g *Generator) Logger() *zap.Logger {
	return g.logger
}

// AddFlags adds flags to command
func (g *Generator) AddFlags(command *cobra.Command) {
	base.AddFlags(command)

	flags := command.Flags()
	flags.SortFlags = false

	flags.StringP(pkg, "p", util.DefaultPackage, "package for model files")
}

// ReadFlags read flags from command
func (g *Generator) ReadFlags(command *cobra.Command) error {
	var err error

	g.options.URL, g.options.Output, g.options.Tables, g.options.FollowFKs, err = base.ReadFlags(command)
	if err != nil {
		return err
	}

	flags := command.Flags()

	if g.options.Package, err = flags.GetString(pkg); err != nil {
		return err
	}

	// setting defaults
	g.options.Def()

	return nil
}

// Generate runs whole generation process
func (g *Generator) Generate() error {
	genna := genna.New(g.options.URL, g.logger)

	entities, err := genna.Read(g.options.Tables, g.options.FollowFKs, false)
	if err != nil {
		return xerrors.Errorf("read database error: %w", err)
	}

	dirname := path.Dir(g.options.Output)

	pack := NewPackage(entities, g.options)
	if err := saveXML(path.Join(dirname, g.options.Package+".xml"), pack); err != nil {
		return xerrors.Errorf("save package error: %w", err)
	}

	eazeEntities := NewEazeEntities(entities)
	if err := saveXML(path.Join(dirname, g.options.Package+".EazeEntity.xml"), eazeEntities); err != nil {
		return xerrors.Errorf("save package error: %w", err)
	}

	return nil
}

func saveXML(filename string, data interface{}) error {
	file, err := util.File(filename)
	if err != nil {
		return xerrors.Errorf("create file error: %w", err)
	}

	writer := io.Writer(file)
	if _, err := writer.Write([]byte("<?xml version=\"1.0\"?>\n")); err != nil {
		return xerrors.Errorf("write xml error: %w", err)
	}

	enc := xml.NewEncoder(writer)
	enc.Indent("", "    ")
	if err := enc.Encode(data); err != nil {
		return xerrors.Errorf("write xml error: %w", err)
	}

	return nil
}
