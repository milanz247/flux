// Package config loads application configuration from the environment.
//
// Values are read from OS environment variables, with a `.env` file in the
// project root loaded first (existing OS variables always win). Only the
// standard library is used.
package config

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config holds every setting the application needs at runtime.
type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Auth     AuthConfig
	Mail     MailConfig
	Vite     ViteConfig
}

type AppConfig struct {
	Name  string
	Env   string // local | production | testing
	Key   string // application secret, used to sign JWTs and tokens
	URL   string
	Host  string
	Port  int
	Debug bool
}

type DatabaseConfig struct {
	Host     string
	Port     int
	Name     string
	User     string
	Password string
	Charset  string
}

// DSN returns the MySQL data source name for GORM.
func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		d.User, d.Password, d.Host, d.Port, d.Name, d.Charset)
}

type AuthConfig struct {
	// TokenTTLMinutes is the lifetime of issued JWT access tokens.
	TokenTTLMinutes int
	// SessionCookie is the name of the httpOnly cookie carrying the JWT.
	SessionCookie string
}

type MailConfig struct {
	// Driver is "log" (write mail to storage/logs/mail.log) or "smtp".
	Driver   string
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

type ViteConfig struct {
	// DevServer is the origin of the Vite dev server used when App.Env == "local".
	DevServer string
	// BuildDir is the directory (under public/) holding the production build.
	BuildDir string
}

// Load reads .env (if present) and the process environment into a Config.
func Load() (*Config, error) {
	loadDotEnv(".env")

	cfg := &Config{
		App: AppConfig{
			Name:  env("APP_NAME", "Flux"),
			Env:   env("APP_ENV", "local"),
			Key:   env("APP_KEY", ""),
			URL:   env("APP_URL", "http://localhost:8080"),
			Host:  env("APP_HOST", "0.0.0.0"),
			Port:  envInt("APP_PORT", 8080),
			Debug: envBool("APP_DEBUG", true),
		},
		Database: DatabaseConfig{
			Host:     env("DB_HOST", "127.0.0.1"),
			Port:     envInt("DB_PORT", 3306),
			Name:     env("DB_DATABASE", "flux"),
			User:     env("DB_USERNAME", "root"),
			Password: env("DB_PASSWORD", ""),
			Charset:  env("DB_CHARSET", "utf8mb4"),
		},
		Auth: AuthConfig{
			TokenTTLMinutes: envInt("AUTH_TOKEN_TTL", 60*24),
			SessionCookie:   env("AUTH_SESSION_COOKIE", "flux_session"),
		},
		Mail: MailConfig{
			Driver:   env("MAIL_DRIVER", "log"),
			Host:     env("MAIL_HOST", "127.0.0.1"),
			Port:     envInt("MAIL_PORT", 1025),
			Username: env("MAIL_USERNAME", ""),
			Password: env("MAIL_PASSWORD", ""),
			From:     env("MAIL_FROM", "hello@flux.local"),
		},
		Vite: ViteConfig{
			DevServer: env("VITE_DEV_SERVER", "http://localhost:5173"),
			BuildDir:  env("VITE_BUILD_DIR", "build"),
		},
	}

	if cfg.App.Key == "" {
		return nil, fmt.Errorf("APP_KEY is not set; add it to .env (any long random string)")
	}
	return cfg, nil
}

// loadDotEnv parses a .env file. Lines are KEY=VALUE; `#` starts a comment;
// values may be quoted. Variables already set in the OS environment win.
func loadDotEnv(path string) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, found := strings.Cut(line, "=")
		if !found {
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		value = strings.Trim(value, `"'`)
		if _, exists := os.LookupEnv(key); !exists {
			os.Setenv(key, value)
		}
	}
}

func env(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

func envInt(key string, fallback int) int {
	if v, ok := os.LookupEnv(key); ok {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}

func envBool(key string, fallback bool) bool {
	if v, ok := os.LookupEnv(key); ok {
		if b, err := strconv.ParseBool(v); err == nil {
			return b
		}
	}
	return fallback
}
