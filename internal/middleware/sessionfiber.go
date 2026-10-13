package middleware

import (
	"log/slog"
	"sigolang/config"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/extractors"
	"github.com/gofiber/fiber/v3/middleware/session"

	// "github.com/gofiber/storage/redis/v3"
	"github.com/gofiber/utils/v2"
)

const SessionKeyPrefix string = "session_id"

func NewSessionFiberMiddleware() fiber.Handler {
	c := config.Get()

	slog.Info("Session Cookie is: " + c.SessionCookie)

	// storage := redis.New(redis.Config{
	// 	URL: c.Cache.CacheUri,
	// })

	sessCfg := session.Config{
		// Storage: storage,
		KeyGenerator: func() string {
			token := utils.SecureToken()
			return SessionKeyPrefix + ":" + token
		},
		CookieSecure:    true,  // HTTPS only
		CookieHTTPOnly:  true,  // Prevent XSS
		CookieSameSite:  "Lax", // CSRF protection
		CookiePath:      "/",
		IdleTimeout:     24 * 2 * time.Hour, // Session timeout
		AbsoluteTimeout: 24 * 7 * time.Hour, // Maximum session life, zero to follow idle timeout
		Extractor:       extractors.FromCookie("__Host-" + c.SessionCookie),
	}

	sessionMiddleware := session.New(sessCfg)
	return sessionMiddleware
}
