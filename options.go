package enber

func WithDBConfig(db DatabaseConfig) extensionOption {
	return func(e *extension) {
		e.Config.DBConfig = &DBConfig{}
		switch db.Driver {
		case SQLite:
			e.Config.DBConfig.Driver = "sqlite3"
			e.Config.DBConfig.DriverImport = "github.com/mattn/go-sqlite3"
			if db.DBName != "" {
				e.Config.DBConfig.Dsn = db.DBName + ".sqlite?_fk=1&cache=shared"
			} else {
				e.Config.DBConfig.Dsn = "db.sqlite?_fk=1&cache=shared"
			}

		case MySQL:
			e.Config.DBConfig.Driver = "mysql"
			e.Config.DBConfig.DriverImport = "github.com/go-sql-driver/mysql"
			host := "127.0.0.1"
			port := "5432"
			if db.Host != "" {
				host = db.Host
			}
			if db.Port != "" {
				port = db.Port
			}
			e.Config.DBConfig.Dsn = "host=" + host + " port=" + port + " user=" + db.User + " password=" + db.Password + " dbname=" + db.DBName

		case PostgreSQL:
			e.Config.DBConfig.Driver = "postgres"
			e.Config.DBConfig.DriverImport = "github.com/lib/pq"
			host := "127.0.0.1"
			port := "5432"
			if db.Host != "" {
				host = db.Host
			}

			if db.Port != "" {
				port = db.Port
			}

			e.Config.DBConfig.Dsn = "host=" + host + " port=" + port + " user=" + db.User + " password=" + db.Password + " dbname=" + db.DBName
		}
	}
}

func WithAuth(method authentication, methods ...authentication) extensionOption {
	return func(e *extension) {
		e.Config.Server.Jwt = JWT == method || vin(JWT, methods)
	}
}

func WithServerConfig(config *ServerConfig) extensionOption {
	return func(e *extension) {
		if config.Pkg == "" {
			config.Pkg = "main"
		}
		if config.FuncName == "" {
			config.FuncName = "main"
		}
		if config.Port == "" {
			config.Port = "5000"
		}
		if config.Filename == "" {
			config.Filename = "main"
		}
		if config.Prefix == "" {
			config.Prefix = "api"
		}
		e.Config.Server = config
	}
}

func WithSwagger(b bool) extensionOption {
	return func(e *extension) {
		e.Config.Server.Swagger = b
	}
}

func WithAppConfig(c AppConfig) extensionOption {
	return func(e *extension) {
		if c.Pkg != "" {
			e.Config.App.Pkg = c.Pkg
		}
		if c.RootPath != "" {
			e.Config.App.RootPath = c.RootPath
		}
	}
}

func WithTypeScript(path string) extensionOption {
	return func(e *extension) {
		e.Config.TsConfig = &TsConfig{
			Path: path,
		}
	}
}

func WithPrivacy(b bool) extensionOption {
	return func(e *extension) {
		e.Config.Privacy = b
	}
}

func WithDebug(b bool) extensionOption {
	return func(e *extension) {
		e.Config.Debug = b
	}
}
