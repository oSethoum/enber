package enber

func WithDBConfig(db DatabaseConfig) extensionOption {
	return func(e *extension) {
		switch db.Driver {
		case SQLite:
			e.Config.Driver = "sqlite3"
			e.Config.DriverImport = "github.com/mattn/go-sqlite3"
			if db.DBName != "" {
				e.Config.Dsn = db.DBName + ".sqlite?_fk=1&cahche=shared"
			} else {
				e.Config.Dsn = "db.sqlite?_fk=1&cache=shared"
			}

		case MySQL:
			e.Config.Driver = "mysql"
			e.Config.DriverImport = "github.com/go-sql-driver/mysql"
			host := "127.0.0.1"
			port := "5432"
			if db.Host != "" {
				host = db.Host
			}
			if db.Port != "" {
				port = db.Port
			}
			e.Config.Dsn = "host=" + host + " port=" + port + " user=" + db.User + " password=" + db.Password + " dbname=" + db.DBName

		case PostgreSQL:
			e.Config.Driver = "postgres"
			e.Config.DriverImport = "github.com/lib/pq"
			host := "127.0.0.1"
			port := "5432"
			if db.Host != "" {
				host = db.Host
			}

			if db.Port != "" {
				port = db.Port
			}

			e.Config.Dsn = "host=" + host + " port=" + port + " user=" + db.User + " password=" + db.Password + " dbname=" + db.DBName
		}
	}
}

func WithAuth(method authentication, methods ...authentication) extensionOption {
	return func(e *extension) {
		e.Config.Jwt = JWT == method || vin(JWT, methods)
	}
}

func WithServerConfig(config ServerConfig) extensionOption {
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
		e.Config.Server = config
	}
}

func WithApiPrefix(prefix string) extensionOption {
	return func(e *extension) {
		e.Config.Prefix = prefix
	}
}

func WithSwagger(b bool) extensionOption {
	return func(e *extension) {
		e.Config.Swagger = b
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
