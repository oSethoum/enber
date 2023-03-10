package enber

import (
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
)

type (
	database        uint
	authentication  string
	extensionOption func(*extension)
	extension       struct {
		entc.DefaultExtension
		hooks        []gen.Hook
		Config       config
		TemplateData *templateData
	}
	config struct {
		DBConfig *DBConfig
		Server   *ServerConfig
		TsConfig *TsConfig
		App      AppConfig
		Privacy  bool
		Debug    bool
	}

	TsConfig struct {
		Path string
	}

	DBConfig struct {
		DriverImport string
		Driver       string
		Dsn          string
		UniqueID     bool
	}

	AppConfig struct {
		Pkg      string
		RootPath string
	}

	ServerConfig struct {
		Pkg      string
		FuncName string
		Port     string
		FileName string
		Prefix   string
		Swagger  bool
		Jwt      bool
	}

	DatabaseConfig struct {
		Driver   database
		UniqueID bool
		DBName   string
		Host     string
		Port     string
		User     string
		Password string
	}

	templateData struct {
		Config       config
		InputNode    *inputNode
		QueryNode    *queryNode
		InputNodes   []*inputNode
		QueryNodes   []*queryNode
		QueryImports []string
		TypesImports []string
		InputImports []string
	}

	inputNode struct {
		Name         string
		IDType       string
		ShouldType   bool
		ShouldInput  bool
		CreateFields []*inputField
		UpdateFields []*inputField
		CreateEdges  []*inputEdge
		UpdateEdges  []*inputEdge
	}

	queryNode struct {
		Name   string
		IDType string
		Fields []*queryField
		Edges  []*queryEdge
	}

	queryField struct {
		Name            string
		TypeString      string
		TsType          string
		DepthType       string
		Optional        bool
		Boolean         bool
		Enum            bool
		Order           bool
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
		TsType   string
		Set      string
		SetParam string
		Slice    bool
		Check    bool
		TsCheck  bool
		Clear    bool
	}

	inputEdge struct {
		Name     string
		Type     string
		TsType   string
		Set      string
		SetParam string
		Slice    bool
		JTag     string
		Check    bool
		TsCheck  bool
		Clear    bool
	}

	file struct {
		Path   string
		Buffer string
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
