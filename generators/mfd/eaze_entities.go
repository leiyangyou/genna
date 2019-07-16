package mfd

import (
	"fmt"
	"strings"

	"github.com/dizzyfool/genna/model"
	"github.com/dizzyfool/genna/util"
)

// NewEazeEntities creates entities for xml
func NewEazeEntities(entities []model.Entity) EazeEntities {
	eazeEntities := make([]EazeEntity, len(entities))
	for i, ent := range entities {
		eazeEntities[i] = newEazeEntity(ent)
	}

	return EazeEntities{
		Entities: eazeEntities,

		XMLxsi: "http://www.w3.org/2001/XMLSchema-instance",
		XMLxsd: "http://www.w3.org/2001/XMLSchema",
	}
}

func newEazeEntity(entity model.Entity) EazeEntity {
	cacheDeps := make([]string, len(entity.Relations))
	for i, relation := range entity.Relations {
		cacheDeps[i] = relation.TargetPGFullName
	}

	search := make([]TemplateAttribute, len(entity.Columns))
	attributes := make([]TemplateAttribute, len(entity.Columns))
	for i, column := range entity.Columns {
		search[i] = newTemplateAttribute(entity, column, true)
		attributes[i] = newTemplateAttribute(entity, column, false)
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

func newTemplateAttribute(entity model.Entity, column model.Column, isSearch bool) TemplateAttribute {
	inputType := inputType(column, isSearch)

	return TemplateAttribute{
		Name:       column.PGName,
		TitleKey:   langPath(entity, column),
		Type:       inputType,
		CanEdit:    canEdit(column, isSearch),
		CanShow:    canShow(column, isSearch),
		CanSearch:  canSearch(column, isSearch),
		Parameters: inputParams(inputType, column, isSearch),
	}
}

func canEdit(column model.Column, isSearch bool) bool {
	if isSearch {
		return false
	}

	if column.PGName == "createdAt" || column.PGName == "modifiedAt" {
		return false
	}
	return true
}

func canShow(column model.Column, isSearch bool) bool {
	if isSearch {
		return false
	}

	if column.PGType == model.TypePGText {
		return false
	}

	return true
}

func canSearch(column model.Column, isSearch bool) bool {
	if !isSearch {
		return false
	}

	return column.GoType != model.TypeTime
}

func inputType(column model.Column, isSearch bool) string {
	if column.PGType == model.TypePGText && !isSearch {
		return TypeHTMLEditor
	}

	if column.IsFK {
		return TypeHTMLSelect
	}

	switch column.PGType {
	case model.TypePGText:
		return TypeHTMLText
	case model.TypeTime, model.TypePGTimetz:
		return TypeHTMLTime
	case model.TypePGDate:
		return TypeHTMLDate
	case model.TypePGTimestamp, model.TypePGTimestamptz:
		return TypeHTMLDateTime
	case model.TypePGBool:
		return TypeHTMLCheckbox
	}

	return TypeHTMLInput
}

func inputParams(inputType string, column model.Column, isSearch bool) Parameters {
	switch inputType {
	case TypeHTMLInput:
		return htmlInputParams(column)
	case TypeHTMLText:
		return htmlTextParams(column)
	case TypeHTMLEditor:
		return htmlEditorParams(column)
	case TypeHTMLTime, TypeHTMLDate, TypeHTMLDateTime:
		return htmlDateTimeParams(column)
	case TypeHTMLCheckbox:
		return htmlCheckboxParams(column)
	case TypeHTMLSelect:
		return htmlSelectParams(column, isSearch || column.Nullable)
	}

	return nil
}

func htmlInputParams(column model.Column) Parameters {
	return Parameters{
		Parameter{Key: "name", Value: column.GoName},
		Parameter{Key: "value", Value: column.GoName},
		Parameter{Key: "size", Value: "80"},
		Parameter{Key: "controlId", Value: util.LowerFirst(column.GoName)},
	}
}

func htmlCheckboxParams(column model.Column) Parameters {
	return Parameters{
		Parameter{Key: "name", Value: column.GoName},
		Parameter{Key: "value", Value: column.GoName},
		Parameter{Key: "controlId", Value: util.LowerFirst(column.GoName)},
	}
}

func htmlTextParams(column model.Column) Parameters {
	return Parameters{
		Parameter{Key: "name", Value: column.GoName},
		Parameter{Key: "value", Value: column.GoName},
		Parameter{Key: "size", Value: "80"},
		Parameter{Key: "rows", Value: "5"},
		Parameter{Key: "controlId", Value: util.LowerFirst(column.GoName)},
	}
}

func htmlEditorParams(column model.Column) Parameters {
	return htmlTextParams(column)
}

func htmlDateTimeParams(column model.Column) Parameters {
	return htmlTextParams(column)
}

func htmlSelectParams(column model.Column, nullable bool) Parameters {
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

func langPath(entity model.Entity, column model.Column) string {
	return fmt.Sprintf("vt.%s.%s", entity.GoName, column.GoName)
}
