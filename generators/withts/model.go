package withts

import (
	"fmt"
	"html/template"
	"strings"

	"github.com/dizzyfool/genna/model"
	"github.com/dizzyfool/genna/util"
)

// TemplatePackage stores package info
type TemplatePackage struct {
	Package string

	HasImports bool
	Imports    []string

	Entities []TemplateEntity
}

// NewTemplatePackage creates a package for template
func NewTemplatePackage(entities []model.Entity, options Options) TemplatePackage {
	imports := util.NewSet()

	models := make([]TemplateEntity, len(entities))
	for i, entity := range entities {
		models[i] = NewTemplateEntity(entity, options)

		for _, imp := range models[i].Imports {
			imports.Add(imp)
		}
	}

	return TemplatePackage{
		Package: options.Package,

		HasImports: imports.Len() > 0,
		Imports:    imports.Elements(),

		Entities: models,
	}
}

// TemplateEntity stores struct info
type TemplateEntity struct {
	model.Entity

	Tag template.HTML

	NoAlias bool
	Alias   string

	Columns []TemplateColumn

	HasRelations bool
	Relations    []TemplateRelation
}

// NewTemplateEntity creates an entity for template
func NewTemplateEntity(entity model.Entity, options Options) TemplateEntity {
	if entity.HasMultiplePKs() {
		options.KeepPK = true
	}

	columns := make([]TemplateColumn, len(entity.Columns))
	for i, column := range entity.Columns {
		columns[i] = NewTemplateColumn(&entity, column, options)
	}

	relations := make([]TemplateRelation, len(entity.Relations))
	for i, relation := range entity.Relations {
		relations[i] = NewTemplateRelation(relation)
	}

	tags := util.NewAnnotation()
	tags.AddTag("sql", util.Quoted(entity.PGFullName, true))
	if !options.NoAlias {
		tags.AddTag("sql", fmt.Sprintf("alias:%s", util.DefaultAlias))
	}

	if !options.NoDiscard {
		// leading comma is required
		tags.AddTag("pg", ",discard_unknown_columns")
	}

	return TemplateEntity{
		Entity: entity,
		Tag:    template.HTML(fmt.Sprintf("`%s`", tags.String())),

		NoAlias: options.NoAlias,
		Alias:   util.DefaultAlias,

		Columns: columns,

		HasRelations: len(relations) > 0,
		Relations:    relations,
	}
}

// TemplateColumn stores column info
type TemplateColumn struct {
	model.Column

	Tag     template.HTML
	Comment template.HTML
	IsCreatedAt bool
	IsUpdatedAt bool
}

// NewTemplateColumn creates a column for template
func NewTemplateColumn(entity *model.Entity, column model.Column, options Options) TemplateColumn {
	if !options.KeepPK && column.IsPK {
		column.GoName = util.ID
	}

	comment := ""
	tags := util.NewAnnotation()
	tags.AddTag("sql", column.PGName)

	if column.IsPK {
		tags.AddTag("sql", "pk")
	}

	if column.PGType == model.TypePGHstore {
		tags.AddTag("sql", "hstore")
	} else if column.IsArray {
		tags.AddTag("sql", "array")
	}

	if !column.Nullable && !column.IsPK {
		tags.AddTag("sql", "notnull")
	}

	if options.SoftDelete == column.PGName && column.Nullable && column.GoType == model.TypeTime && !column.IsArray {
		tags.AddTag("pg", "soft_delete")
	}

	// As default is not exposed in genna, assuming default is now

	isCreatedAt := false
	if options.CreatedAt == column.PGName && column.GoType == model.TypeTime && !column.IsArray {
		tags.AddTag("sql", "default:now()")
		isCreatedAt = true
	}

	isUpdatedAt := false
	if options.UpdatedAt == column.PGName && column.GoType == model.TypeTime && !column.IsArray {
		tags.AddTag("sql", "default:now()")
		entity.Imports = append(entity.Imports,  "context", "github.com/go-pg/pg/orm")
		isUpdatedAt = true
	}

	if column.GoType == model.TypeInterface {
		comment = "// unsupported"
		tags = util.NewAnnotation().AddTag("sql", "-")
	}

	return TemplateColumn{
		Column: column,

		Tag:     template.HTML(fmt.Sprintf("`%s`", tags.String())),
		Comment: template.HTML(comment),
		IsCreatedAt: isCreatedAt,
		IsUpdatedAt: isUpdatedAt,
	}
}

// TemplateRelation stores relation info
type TemplateRelation struct {
	model.Relation

	Tag     template.HTML
	Comment template.HTML
}

// NewTemplateRelation creates relation for template
func NewTemplateRelation(relation model.Relation) TemplateRelation {
	comment := ""
	tags := util.NewAnnotation().AddTag("pg", "fk:"+strings.Join(relation.FKFields, ","))
	if len(relation.FKFields) > 1 {
		comment = "// unsupported"
		tags.AddTag("sql", "-")
	}

	return TemplateRelation{
		Relation: relation,

		Tag:     template.HTML(fmt.Sprintf("`%s`", tags.String())),
		Comment: template.HTML(comment),
	}
}
