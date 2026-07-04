package goswagger

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
	"github.com/swaggo/swag"
)

// Configure applies Swagger host, base path, and schemes from environment variables.
//
// Environment-driven behaviour (aligned across all backend services):
//   - prod: callers should skip registering /swagger (see Enabled).
//   - dev/staging: HOST + SWAGGER_BASE_PATH (gateway); defaults to defaultBasePath when unset.
//   - local/other: 127.0.0.1:APP_PORT, empty base path unless SWAGGER_BASE_PATH is set.
//
// Helm injects HOST, ENV, SWAGGER_BASE_PATH (from ingress.path), and SWAGGER_SCHEMES.
// defaultBasePath is the service gateway prefix (e.g. "/api/payment") used when
// SWAGGER_BASE_PATH is not set.
func Configure(spec *swag.Spec, appPort, defaultBasePath string) {
	if spec == nil {
		return
	}

	env := strings.ToLower(strings.TrimSpace(viper.GetString("ENV")))
	port := strings.TrimSpace(appPort)
	if port == "" {
		port = strings.TrimSpace(viper.GetString("APP_PORT"))
	}
	if port == "" {
		port = "8080"
	}

	defaultBasePath = strings.TrimRight(strings.TrimSpace(defaultBasePath), "/")

	switch env {
	case "dev", "staging":
		applyGateway(spec, defaultBasePath)
	default:
		applyLocal(spec, port, defaultBasePath)
	}
}

// Enabled reports whether Swagger UI should be registered (false in production).
func Enabled() bool {
	return strings.ToLower(strings.TrimSpace(viper.GetString("ENV"))) != "prod"
}

func applyGateway(spec *swag.Spec, defaultBasePath string) {
	host := strings.TrimSpace(viper.GetString("HOST"))
	if host != "" {
		spec.Host = host
	}

	bp := swaggerBasePath(defaultBasePath)
	spec.BasePath = bp

	if sch := parseSchemes(viper.GetString("SWAGGER_SCHEMES")); len(sch) > 0 {
		spec.Schemes = sch
	} else if bp != "" {
		spec.Schemes = []string{"https"}
	}
}

func applyLocal(spec *swag.Spec, port, defaultBasePath string) {
	spec.Host = fmt.Sprintf("127.0.0.1:%s", port)

	if bp := strings.TrimRight(strings.TrimSpace(viper.GetString("SWAGGER_BASE_PATH")), "/"); bp != "" {
		spec.BasePath = bp
	} else {
		spec.BasePath = ""
	}

	if sch := parseSchemes(viper.GetString("SWAGGER_SCHEMES")); len(sch) > 0 {
		spec.Schemes = sch
	} else {
		spec.Schemes = []string{"http"}
	}

	_ = defaultBasePath
}

func swaggerBasePath(defaultBasePath string) string {
	if bp := strings.TrimRight(strings.TrimSpace(viper.GetString("SWAGGER_BASE_PATH")), "/"); bp != "" {
		return bp
	}
	return defaultBasePath
}

func parseSchemes(raw string) []string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}
	var out []string
	for _, p := range strings.Split(raw, ",") {
		p = strings.TrimSpace(strings.ToLower(p))
		if p == "http" || p == "https" {
			out = append(out, p)
		}
	}
	return out
}
