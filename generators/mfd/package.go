package mfd

import (
	"fmt"
	"strings"

	"github.com/dizzyfool/genna/model"
	"github.com/dizzyfool/genna/util"
)

// NewPackage creates package for mfg
func NewPackage(entities []model.Entity, options Options) Package {
	xmlEntities := make([]*Entity, len(entities))
	for i, ent := range entities {
		xmlEntities[i] = newEntity(ent, options)
	}

	return Package{
		Name: options.Package,

		Entities: xmlEntities,

		XMLxsi: "http://www.w3.org/2001/XMLSchema-instance",
		XMLxsd: "http://www.w3.org/2001/XMLSchema",
	}
}

func newEntity(entity model.Entity, options Options) *Entity {
	searches := Attributes{
		newPageSearch(),
		newPageSizeSearch(),
	}

	attributes := Attributes{}
	for _, column := range entity.Columns {
		if search := newSuggestedSearch(column); search != nil {
			searches = append(searches, search)
		}
		attributes = append(attributes, newAttribute(column))
	}

	return &Entity{
		Name:              entity.GoName,
		PackageName:       options.Package,
		Table:             entity.PGFullName,
		View:              view(entity.PGSchema, entity.PGName),
		DefaultConnection: "",
		Attributes:        &attributes,
		Search:            &searches,
		Flags: []string{
			CanPages, CanCache,
		},
	}
}

func view(schema, table string) string {
	return util.JoinF(schema, fmt.Sprintf("get%s", strings.Title(table)))
}

func newAttribute(column model.Column) *Attribute {
	fkModel := ""
	if column.IsFK && column.Relation != nil {
		fkModel = column.Relation.GoType
	}

	return &Attribute{
		Name:        column.GoName,
		Key:         column.IsPK,
		Addable:     addable(column),
		Updatable:   updatable(column),
		MinValue:    0,
		MaxValue:    column.MaxLen,
		DbName:      column.PGName,
		FieldType:   fieldType(column),
		ForeignKey:  fkModel,
		ComplexType: complexType(column),
		Nullable:    nullable(column),
		GoType:      column.Type,
		SearchType:  SearchEquals,
	}
}

func nullable(column model.Column) string {
	switch {
	case column.IsPK || column.Nullable:
		return NullableYes
	case column.GoType == model.TypeString:
		return NullableEmpty
	default:
		return NullableNo
	}
}

func updatable(column model.Column) bool {
	if column.PGName == "createdAt" || column.PGName == "modifiedAt" {
		return false
	}

	return true
}

func addable(column model.Column) bool {
	if column.PGName == "createdAt" || column.PGName == "modifiedAt" {
		return false
	}

	return true
}

func fieldType(column model.Column) string {
	if column.IsArray {
		return Array
	}

	if column.PGType == model.TypePGDate {
		return Date
	}

	if column.PGType == model.TypePGTimetz || column.PGType == model.TypePGTime {
		return Time
	}

	switch column.GoType {
	case model.TypeString:
		return String
	case model.TypeInt, model.TypeInt32, model.TypeInt64:
		return Integer
	case model.TypeFloat32, model.TypeFloat64:
		return Float
	case model.TypeBool:
		return Boolean
	case model.TypeTime:
		return DateTime
	case model.TypeMapString, model.TypeMapInterface:
		return Array
	}

	return String
}

func phpType(column model.Column) string {
	switch column.GoType {
	case model.TypeString:
		return "string"
	case model.TypeInt, model.TypeInt32, model.TypeInt64:
		return "int"
	case model.TypeFloat32, model.TypeFloat64:
		return "float"
	case model.TypeBool:
		return "boolean"
	}

	return "string"
}

func complexType(column model.Column) string {
	if column.IsArray {
		phpType := phpType(column)
		for i := 0; i < column.Dimensions; i++ {
			phpType = fmt.Sprintf("[]%s", phpType)
		}
		return phpType
	}

	if column.PGType == model.TypePGJSON || column.PGType == model.TypePGJSONB {
		return "json"
	}

	return ""
}

func newSuggestedSearch(column model.Column) *Attribute {
	// TODO add suggested searches
	if !column.IsArray && column.PGType == model.TypePGText {
		return newSearch(
			fmt.Sprintf("%sILike", column.GoName),
			column.PGName,
			"",
			String,
			column.GoType,
			SearchILike,
		)
	}

	return nil
}

func newPageSearch() *Attribute {
	return newSearch(SearchPage, SearchPage, "0", Integer, Integer, SearchEquals)
}

func newPageSizeSearch() *Attribute {
	return newSearch(SearchPageSize, SearchPageSize, "25", Integer, Integer, SearchEquals)
}

func newSearch(name, dbName, defaultValue, fieldType, goType, searchType string) *Attribute {
	return &Attribute{
		Name:         name,
		Key:          false,
		Addable:      false,
		Updatable:    false,
		MinValue:     0,
		MaxValue:     0,
		DbName:       dbName,
		DefaultValue: defaultValue,
		FieldType:    fieldType,
		GoType:       goType,
		Nullable:     "Yes",
		SearchType:   searchType,
	}
}
