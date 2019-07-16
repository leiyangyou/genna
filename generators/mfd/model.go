package mfd

import (
	"encoding/xml"
)

// input types
const (
	TypeHTMLInput    = "HTML_INPUT"
	TypeHTMLText     = "HTML_TEXT"
	TypeHTMLEditor   = "HTML_EDITOR"
	TypeHTMLCheckbox = "HTML_CHECKBOX"
	TypeHTMLDateTime = "HTML_DATETIME"
	TypeHTMLDate     = "HTML_DATE"
	TypeHTMLTime     = "HTML_TIME"
	TypeHTMLSelect   = "HTML_SELECT"
	TypeHTMLFile     = "HTML_FILE"
)

// types
const (
	String   = "TYPE_STRING"
	Integer  = "TYPE_INTEGER"
	Float    = "TYPE_FLOAT"
	Boolean  = "TYPE_BOOLEAN"
	DateTime = "TYPE_DATETIME"
	Date     = "TYPE_DATE"
	Time     = "TYPE_TIME"
	Array    = "TYPE_ARRAY"
)

// model options
const (
	CanPages = "CanPages"
	CanCache = "CanCache"
)

// nullable options
const (
	NullableYes   = "Yes"
	NullableNo    = "No"
	NullableEmpty = "CheckEmpty"
)

// special searches
const (
	SearchPage     = "page"
	SearchPageSize = "pageSize"
)

// search types
const (
	SearchEquals     = "SEARCHTYPE_EQUALS"
	SearchNotEquals  = "SEARCHTYPE_NOT_EQUALS"
	SearchNull       = "SEARCHTYPE_NULL"
	SearchNotNull    = "SEARCHTYPE_NOT_NULL"
	SearchGE         = "SEARCHTYPE_GE"
	SearchLE         = "SEARCHTYPE_LE"
	SearchG          = "SEARCHTYPE_G"
	SearchL          = "SEARCHTYPE_L"
	SearchLeftLike   = "SEARCHTYPE_LEFT_LIKE"
	SearchLeftILike  = "SEARCHTYPE_LEFT_ILIKE"
	SearchRightLike  = "SEARCHTYPE_RIGHT_LIKE"
	SearchRightILike = "SEARCHTYPE_RIGHT_ILIKE"
	SearchLike       = "SEARCHTYPE_LIKE"
	SearchILike      = "SEARCHTYPE_ILIKE"
	SearchArray      = "SEARCHTYPE_ARRAY"
	SearchNotArray   = "SEARCHTYPE_NOT_INARRAY"
)

// Project is xml element
type Project struct {
	XMLName         xml.Name `xml:"Project" json:"-"`
	Name            string
	TemplateVersion string
	PackageNames    []string `xml:"PackageNames>string" json:"-"`
	Packages        Packages
	Filename        string `xml:"-"`
}

// Package is xml element
type Package struct {
	XMLName  xml.Name `xml:"Package" json:"-"`
	XMLxsi   string   `xml:"xmlns:xsi,attr"`
	XMLxsd   string   `xml:"xmlns:xsd,attr"`
	Name     string
	Entities []*Entity `xml:"Entities>Entity"`
}

// Entity is xml element
type Entity struct {
	XMLName           xml.Name    `xml:"Entity" json:"-"`
	Name              string      `xml:"Name,attr"`
	PackageName       string      `xml:"PackageName,attr"`
	Table             string      `xml:"Table,attr"`
	View              string      `xml:"View,attr"`
	DefaultConnection string      `xml:"DefaultConnection,attr"`
	Attributes        *Attributes `xml:"Attributes>Attribute,omitempty"`
	Search            *Attributes `xml:"Search>Attribute,omitempty"`
	Lists             *Attributes `xml:"Lists>Attribute,omitempty"`
	Flags             []string    `xml:"Flags>EFlag,omitempty"`
	EazeEntity        EazeEntity  `xml:"-"`
	Package           *Package    `xml:"-" json:"-"`
}

// Attribute is xml element
type Attribute struct {
	XMLName      xml.Name `xml:"Attribute" json:"-"`
	Name         string   `xml:"Name,attr"`
	Key          bool     `xml:"Key,attr"`
	Addable      bool     `xml:"Addable,attr"`
	Updatable    bool     `xml:"Updatable,attr"`
	MinValue     int      `xml:"MinValue,attr"`
	MaxValue     int      `xml:"MaxValue,attr"`
	DbName       string   `xml:"DbName,attr"`
	ComplexType  string   `xml:"ComplexType,attr,omitempty"`
	DefaultValue string   `xml:"DefaultValue,attr,omitempty"`
	ForeignKey   string   `xml:"ForeignKey,attr,omitempty"`
	FieldType    string   `xml:"FieldType,attr"`
	GoType       string   `xml:"GoType,attr"`
	Nullable     string   `xml:"Nullable,attr"`
	SearchType   string   `xml:"SearchType,attr"`
}

// EazeEntities is xml element
type EazeEntities struct {
	XMLName  xml.Name     `xml:"ArrayOfEazeEntity" json:"-"`
	XMLxsi   string       `xml:"xmlns:xsi,attr"`
	XMLxsd   string       `xml:"xmlns:xsd,attr"`
	Entities []EazeEntity `xml:"EazeEntity"`
}

// EazeEntity is xml element
type EazeEntity struct {
	XMLName           xml.Name `xml:"EazeEntity" json:"-"`
	Name              string   `xml:"Name,attr"`
	TerminalPath      string
	CacheTime         int
	CacheDependencies string
	Search            []TemplateAttribute `xml:"Search>TemplateAttribute"`
	Attributes        []TemplateAttribute `xml:"Attributes>TemplateAttribute"`
}

// TemplateAttribute is xml element
type TemplateAttribute struct {
	XMLName    xml.Name `xml:"TemplateAttribute" json:"-"`
	Name       string   `xml:"Name,attr"`
	TitleKey   string   `xml:"TitleKey,attr"`
	TitleValue string   `xml:"TitleValue,attr"`
	Type       string   `xml:"Type,attr"`
	CanEdit    bool     `xml:"CanEdit,attr"`
	CanShow    bool     `xml:"CanShow,attr"`
	CanSearch  bool     `xml:"CanSearch,attr"`
	Parameters `xml:"Parameters>AttributeParameter"`
}

// Parameter is xml element
type Parameter struct {
	XMLName xml.Name `xml:"AttributeParameter" json:"-"`
	Key     string   `xml:"Key,attr"`
	Value   string   `xml:"Value,attr"`
}

// Packages is xml element
type Packages map[string]*Package

// Attributes is xml element
type Attributes []*Attribute

// Parameters is xml element
type Parameters []Parameter
