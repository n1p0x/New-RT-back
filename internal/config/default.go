package config

var defaultConfig = map[string]interface{}{
	"mode": "debug",

	"server.host":             "localhost",
	"server.port":             8000,
	"service.writeTimeout":    "15s",
	"server.gracefulShutdown": "30s",

	"db.host":           "localhost",
	"db.user":           "postgres",
	"db.password":       "postgres",
	"db.name":           "roulette",
	"db.port":           5432,
	"db.migrate.enable": true,
	"db.migrate.dir":    "",

	"tg.botToken":      "",
	"tg.adminBotToken": "",
	"tg.clientID":      "",
	"tg.clientHash":    "",
	"tg.clientPhone":   "",

	"ton.isTestnet":              true,
	"ton.tonCenterApiKey":        "",
	"ton.tonCenterApiKeyTestnet": "",
	"ton.adminWallet":            "",
	"ton.mnemonic":               "",
}
