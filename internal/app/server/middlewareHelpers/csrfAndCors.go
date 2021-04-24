package middlewareHelpers

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"net/http"
	"net/url"
)

func SetupCSRFAndCORS(e *echo.Echo, allowedOrigin string, debug bool) {
	if len(allowedOrigin) > 0 {
		url, err := url.Parse(allowedOrigin)
		if err != nil {
			log.Fatal(err)
		}
		csrfCookieDomain := url.Hostname()
		if len(csrfCookieDomain) == 0 {
			log.Fatal("Invalid domain specified in allowedOrigin")
		}

		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			Skipper: func(c echo.Context) bool {
				return debug
			},
			AllowOrigins:     []string{allowedOrigin},
			AllowCredentials: true,
		}))

		e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
			Skipper: func(c echo.Context) bool {
				return debug
			},
			CookieSameSite: http.SameSiteStrictMode,
			CookieDomain:   csrfCookieDomain,
			CookiePath:     "/",
		}))
		log.Printf("INFO: %s added to CORS and CSRF protection enabled for %s\n", allowedOrigin, csrfCookieDomain)
	} else {
		log.Println("WARN: allowedOrigin not set in config => CORS and CSRF middlewares are not enabled!")
	}
}
