package mfd

import (
	"encoding/xml"
)

const (
	TypeHtmlInput    = "HTML_INPUT"
	TypeHtmlText     = "HTML_TEXT"
	TypeHtmlEditor   = "HTML_EDITOR"
	TypeHtmlCheckbox = "HTML_CHECKBOX"
	TypeHtmlDateTime = "HTML_DATETIME"
	TypeHtmlDate     = "HTML_DATE"
	TypeHtmlTime     = "HTML_TIME"
	TypeHtmlSelect   = "HTML_SELECT"
	TypeHtmlFile     = "HTML_FILE"

	String   = "TYPE_STRING"
	Integer  = "TYPE_INTEGER"
	Float    = "TYPE_FLOAT"
	Boolean  = "TYPE_BOOLEAN"
	DateTime = "TYPE_DATETIME"
	Date     = "TYPE_DATE"
	Time     = "TYPE_TIME"
	Array    = "TYPE_ARRAY"

	CanPages = "CanPages"
	CanCache = "CanCache"

	NullableYes   = "Yes"
	NullableNo    = "No"
	NullableEmpty = "CheckEmpty"

	SearchPage     = "page"
	SearchPageSize = "pageSize"

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

type Project struct {
	XMLName         xml.Name `xml:"Project" json:"-"`
	Name            string
	TemplateVersion string
	PackageNames    []string `xml:"PackageNames>string" json:"-"`
	Packages        Packages
	Filename        string `xml:"-"`
}

type Package struct {
	XMLName  xml.Name `xml:"Package" json:"-"`
	XMLxsi   string   `xml:"xmlns:xsi,attr"`
	XMLxsd   string   `xml:"xmlns:xsd,attr"`
	Name     string
	Entities []*Entity `xml:"Entities>Entity"`
}

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

type EazeEntities struct {
	XMLName  xml.Name     `xml:"ArrayOfEazeEntity" json:"-"`
	XMLxsi   string       `xml:"xmlns:xsi,attr"`
	XMLxsd   string       `xml:"xmlns:xsd,attr"`
	Entities []EazeEntity `xml:"EazeEntity"`
}

type EazeEntity struct {
	XMLName           xml.Name `xml:"EazeEntity" json:"-"`
	Name              string   `xml:"Name,attr"`
	TerminalPath      string
	CacheTime         int
	CacheDependencies string
	Search            []TemplateAttribute `xml:"Search>TemplateAttribute"`
	Attributes        []TemplateAttribute `xml:"Attributes>TemplateAttribute"`
}

//<TemplateAttribute Name="formId" TitleKey="vt.form.formId" Type="HTML_INPUT" CanEdit="false" CanShow="false" CanSearch="false">
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
type Parameter struct {
	XMLName xml.Name `xml:"AttributeParameter" json:"-"`
	Key     string   `xml:"Key,attr"`
	Value   string   `xml:"Value,attr"`
}

type Packages map[string]*Package
type Attributes []*Attribute
type Parameters []Parameter

// Actions.xml Model

type ActionsXML struct {
	XMLName xml.Name       `xml:"actions"`
	Actions Actions        `xml:"action"`
	Groups  []*ActionGroup `xml:"group"`
}

type ActionGroup struct {
	XMLName xml.Name       `xml:"group"`
	Name    string         `xml:"name,attr"`
	Groups  []*ActionGroup `xml:"group"`
	Actions Actions        `xml:"action"`
}

type Action struct {
	XMLName   xml.Name         `xml:"action"`
	Name      string           `xml:"name,attr"`
	Path      string           `xml:"path,omitempty"`
	Request   *ActionParams    `xml:"parameters>request>param,omitempty"`
	Response  *ActionParams    `xml:"parameters>response>param,omitempty"`
	Session   *ActionParams    `xml:"parameters>session>param,omitempty"`
	Redirects *ActionRedirects `xml:"redirects>redirect,omitempty"`
}

type ActionRedirect struct {
	XMLName xml.Name `xml:"redirect,omitempty"`
	Name    string   `xml:"name,attr"`
	Path    string   `xml:"path,attr"`
	Host    string   `xml:"host,attr,omitempty"`
}

type ActionParam struct {
	XMLName xml.Name `xml:"param,omitempty"`
	Name    string   `xml:"name,attr"`
	Value   string   `xml:",innerxml"`
}
type Actions []*Action
type ActionParams []*ActionParam
type ActionRedirects []*ActionRedirect

// Pages.xml Model
type PagesXML struct {
	XMLName xml.Name    `xml:"sites"`
	Sites   []*PageSite `xml:"site"`
}

type PageSite struct {
	XMLName xml.Name   `xml:"site"`
	Name    string     `xml:"name,attr"`
	Names   string     `xml:"names,attr,omitempty"`
	Hosts   PageHosts  `xml:"hosts>host,omitempty"`
	Groups  PageGroups `xml:"pages>pageGroup,omitempty"`
	Pages   Pages      `xml:"pages>page,omitempty"`
}

type Page struct {
	XMLName  xml.Name `xml:"page"`
	Uri      string   `xml:"uri,attr"`
	Actions  string   `xml:"actions,omitempty"`
	Template string   `xml:"template,omitempty"`
	PageBootShutdown
	PageSitemap
}

type PageGroup struct {
	XMLName    xml.Name   `xml:"pageGroup"`
	Name       string     `xml:"name,attr"`
	PageGroups PageGroups `xml:"pageGroup,omitempty"`
	Pages      Pages      `xml:"page,omitempty"`
	PageBootShutdown
	PageSitemap
}

type PageHost struct {
	XMLName xml.Name    `xml:"host"`
	Name    string      `xml:"name,attr"`
	Actions PageActions `xml:"action"`
}

type PageAction struct {
	XMLName xml.Name `xml:"action"`
	Name    string   `xml:"name,attr"`
	Value   string   `xml:",innerxml"`
}

type PageSitemap struct {
	SitemapEnable bool `xml:"sitemap-enable,attr,omitempty"`
}

type PageBootShutdown struct {
	Boot     *string `xml:"boot,attr"`
	Shutdown *string `xml:"shutdown,attr"`
}

type PageHosts []*PageHost
type Pages []*Page
type PageGroups []*PageGroup
type PageActions []*PageAction

// ru.xml
type LanguageXML struct {
	XMLName  xml.Name      `xml:"language"`
	Language string        `xml:"name,attr"`
	Nodes    LanguageNodes `xml:",any"`
}

type LanguageNode struct {
	XMLName  xml.Name
	Value    string        `xml:",chardata"`
	RawValue string        `xml:",innerxml"`
	Nodes    LanguageNodes `xml:",any"`
}

type LanguageNodes []*LanguageNode
