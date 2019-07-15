package mfd

import (
	"fmt"
	"strings"

	"github.com/dizzyfool/genna/model"
	"github.com/dizzyfool/genna/util"
)

func NewEazeEntities(entities []model.Entity) EazeEntities {
	eazeEntities := make([]EazeEntity, len(entities))
	for i, ent := range entities {
		eazeEntities[i] = NewEazeEntity(ent)
	}

	return EazeEntities{
		Entities: eazeEntities,

		XMLxsi: "http://www.w3.org/2001/XMLSchema-instance",
		XMLxsd: "http://www.w3.org/2001/XMLSchema",
	}
}

func NewEazeEntity(entity model.Entity) EazeEntity {
	cacheDeps := make([]string, len(entity.Relations))
	for i, relation := range entity.Relations {
		cacheDeps[i] = relation.TargetPGFullName
	}

	search := make([]TemplateAttribute, len(entity.Columns))
	attributes := make([]TemplateAttribute, len(entity.Columns))
	for i, column := range entity.Columns {
		search[i] = NewTemplateAttribute(entity, column, false)
		attributes[i] = NewTemplateAttribute(entity, column, true)
	}

	return EazeEntity{
		Name:              entity.GoName,
		TerminalPath:      util.Dash(entity.PGName),
		CacheTime:         0,
		CacheDependencies: strings.Join(cacheDeps, ","),
		Search:            search,
		Attributes:        attributes,
	}
}

func NewTemplateAttribute(entity model.Entity, column model.Column, isSearch bool) TemplateAttribute {
	inputType := InputType(column, true)

	return TemplateAttribute{
		Name:       column.PGName,
		TitleKey:   LangPath(entity, column),
		Type:       inputType,
		CanEdit:    CanEdit(column, isSearch),
		CanShow:    CanShow(column, isSearch),
		CanSearch:  CanSearch(column, isSearch),
		Parameters: InputParams(inputType, column, isSearch),
	}
}

func CanEdit(column model.Column, isSearch bool) bool {
	if isSearch {
		return false
	}

	if column.PGName == "createdAt" || column.PGName == "modifiedAt" {
		return false
	}
	return true
}

func CanShow(column model.Column, isSearch bool) bool {
	if isSearch {
		return false
	}

	if column.PGType == model.TypePGText {
		return false
	}

	return true
}

func CanSearch(column model.Column, isSearch bool) bool {
	if !isSearch {
		return false
	}

	return column.GoType != model.TypeTime
}

func InputType(column model.Column, isSearch bool) string {
	if column.PGType == model.TypePGText && column.PGName == "description" && !isSearch {
		return TypeHtmlEditor
	}

	if column.IsFK {
		return TypeHtmlSelect
	}

	switch column.PGType {
	case model.TypePGText:
		return TypeHtmlText
	case model.TypeTime, model.TypePGTimetz:
		return TypeHtmlTime
	case model.TypePGDate:
		return TypeHtmlDate
	case model.TypePGTimestamp, model.TypePGTimestamptz:
		return TypeHtmlDateTime
	case model.TypePGBool:
		return TypeHtmlCheckbox
	}

	return TypeHtmlInput
}

func InputParams(inputType string, column model.Column, isSearch bool) Parameters {
	switch inputType {
	case TypeHtmlInput:
		return HtmlInputParams(column)
	case TypeHtmlText:
		return HtmlTextParams(column)
	case TypeHtmlEditor:
		return HtmlEditorParams(column)
	case TypeHtmlTime, TypeHtmlDate, TypeHtmlDateTime:
		return HtmlDateTimeParams(column)
	case TypeHtmlCheckbox:
		return HtmlCheckboxParams(column)
	case TypeHtmlSelect:
		return HtmlSelectParams(column, isSearch || column.Nullable)
	}

	return nil
}

func HtmlInputParams(column model.Column) Parameters {
	return Parameters{
		Parameter{Key: "name", Value: column.GoName},
		Parameter{Key: "value", Value: column.GoName},
		Parameter{Key: "size", Value: "80"},
		Parameter{Key: "controlId", Value: util.LowerFirst(column.GoName)},
	}
}

func HtmlCheckboxParams(column model.Column) Parameters {
	return Parameters{
		Parameter{Key: "name", Value: column.GoName},
		Parameter{Key: "value", Value: column.GoName},
		Parameter{Key: "controlId", Value: util.LowerFirst(column.GoName)},
	}
}

func HtmlTextParams(column model.Column) Parameters {
	return Parameters{
		Parameter{Key: "name", Value: column.GoName},
		Parameter{Key: "value", Value: column.GoName},
		Parameter{Key: "size", Value: "80"},
		Parameter{Key: "rows", Value: "5"},
		Parameter{Key: "controlId", Value: util.LowerFirst(column.GoName)},
	}
}

func HtmlEditorParams(column model.Column) Parameters {
	return HtmlTextParams(column)
}

func HtmlDateTimeParams(column model.Column) Parameters {
	return HtmlTextParams(column)
}

func HtmlSelectParams(column model.Column, nullable bool) Parameters {
	n := "false"
	if nullable {
		n = "true"
	}

	return Parameters{
		Parameter{Key: "name", Value: column.GoName},
		Parameter{Key: "values", Value: ""},
		Parameter{Key: "dataKey", Value: ""},
		Parameter{Key: "dataValue", Value: ""},
		Parameter{Key: "value", Value: column.GoName},
		Parameter{Key: "controlId", Value: util.LowerFirst(column.GoName)},
		Parameter{Key: "nullValue", Value: n},
	}
}

func LangPath(entity model.Entity, column model.Column) string {
	return fmt.Sprintf("vt.%s.%s", entity.GoName, column.GoName)
}
