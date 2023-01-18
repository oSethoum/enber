package enber

import (
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
)

type (
	database        uint
	authentication  string
	extensionOption func(*extension)
	DebugInfo       struct {
		SchemaJson bool
		DebugOnly  bool
	}
	extension struct {
		entc.DefaultExtension
		hooks        []gen.Hook
		Config       config
		TemplateData *templateData
		DebugInfo    *DebugInfo
	}
	config struct {
		Jwt          bool
		Swagger      bool
		Pkg          string
		Prefix       string
		DriverImport string
		Driver       string
		Dsn          string
		Server       ServerConfig
		App          AppConfig
		Debug        bool
	}

	DatabaseConfig struct {
		Driver   database
		DBName   string
		Host     string
		Port     string
		User     string
		Password string
	}

	AppConfig struct {
		Pkg      string
		RootPath string
	}

	ServerConfig struct {
		Pkg      string
		FuncName string
		Port     string
		Filename string
	}

	templateData struct {
		Config       config
		InputNodes   []*inputNode
		QueryNodes   []*queryNode
		QueryImports []string
		TypesImports []string
		InputImports []string
	}

	inputNode struct {
		Name         string
		ShouldType   bool
		ShouldInput  bool
		CreateFields []*inputField
		UpdateFields []*inputField
		CreateEdges  []*inputEdge
		UpdateEdges  []*inputEdge
	}

	queryNode struct {
		Name   string
		Fields []*queryField
		Edges  []*queryEdge
	}

	queryField struct {
		Name            string
		TypeString      string
		Optional        bool
		Boolean         bool
		Enum            bool
		Comparable      bool
		String          bool
		WithComment     bool
		EdgeFieldOrEnum bool
	}

	queryEdge struct {
		Name string
		Node string
	}

	inputField struct {
		Name     string
		Type     string
		Set      string
		SetParam string
		Slice    bool
		Check    bool
		Clear    bool
	}

	inputEdge struct {
		Name     string
		Type     string
		Set      string
		SetParam string
		Slice    bool
		JTag     string
		Check    bool
		Clear    bool
	}

	file struct {
		Path   string
		Buffer string
	}
)

type (
	jnode struct {
		Name   string
		Fields []*jfield
		Edges  []*jedge
	}
	jfield struct {
		Name          string
		StructField   string
		Default       bool
		Enums         []gen.Enum
		Unique        bool
		Nillable      bool
		Optional      bool
		Type          string
		TypePkgPath   string
		TypePkgName   string
		TypeConstName string
		TypeIdent     string
	}
	jedge struct {
		Name     string
		Optional bool
		Inverse  string
		Unique   bool
	}
	jgraph struct {
		Nodes []*jnode
	}
)

const (
	SQLite database = iota
	MySQL
	PostgreSQL
)

const (
	JWT authentication = "JWT"
)
