package config

var defaultConfig = map[string]interface{}{
	"mode": "debug",

	"origin": "http://localhost:5173",

	"server.host":             "localhost",
	"server.port":             8000,
	"service.writeTimeout":    "15s",
	"server.gracefulShutdown": "30s",

	"db.host":     "localhost",
	"db.user":     "postgres",
	"db.password": "postgres",
	"db.port":     5432,

	"ton.isTestnet": true,
}
