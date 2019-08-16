package settings

import (
	"os"
)

var RPDisplayName = getenv("RP_NAME", "c't demo")
var RPID = getenv("RP_ID", "localhost")

var CookieName = getenv("COOKIE_NAME", "webauthdemo")

var ApiPort = getenv("SERVER_PORT", "80")
var ApiDebug = getenv("SERVER_DEBUG", "false")
var ApiName = getenv("SERVER_NAME", "webauthn")

var JWTKey = getenv("JWT_KEY", "secret_key")
var JWTTime = getenv("JWT_Time", "60")
var JWTIssuer = getenv("JWT_ISSUER", "localhost")

var Timezone = getenv("SERVER_TIMEZONE", "Europe/Berlin")

var SqliteStorage = getenv("SQLITE_STORAGE", "/db/gorm.db")

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
