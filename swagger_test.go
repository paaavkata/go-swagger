package goswagger

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/swaggo/swag"
)

func TestConfigure_devGateway(t *testing.T) {
	viper.Reset()
	viper.Set("ENV", "dev")
	viper.Set("HOST", "dev.file-convert.online")
	viper.Set("SWAGGER_BASE_PATH", "/api/payment")

	spec := &swag.Spec{}
	Configure(spec, "8080", "/api/fallback")

	if spec.Host != "dev.file-convert.online" {
		t.Fatalf("host: got %q", spec.Host)
	}
	if spec.BasePath != "/api/payment" {
		t.Fatalf("basePath: got %q", spec.BasePath)
	}
	if len(spec.Schemes) != 1 || spec.Schemes[0] != "https" {
		t.Fatalf("schemes: got %v", spec.Schemes)
	}
}

func TestConfigure_devDefaultBasePath(t *testing.T) {
	viper.Reset()
	viper.Set("ENV", "dev")
	viper.Set("HOST", "dev.file-convert.online")

	spec := &swag.Spec{}
	Configure(spec, "8080", "/api/conversion")

	if spec.BasePath != "/api/conversion" {
		t.Fatalf("basePath: got %q", spec.BasePath)
	}
}

func TestConfigure_localDirect(t *testing.T) {
	viper.Reset()
	viper.Set("ENV", "local")

	spec := &swag.Spec{}
	Configure(spec, "9090", "/api/payment")

	if spec.Host != "127.0.0.1:9090" {
		t.Fatalf("host: got %q", spec.Host)
	}
	if spec.BasePath != "" {
		t.Fatalf("basePath: got %q, want empty for local", spec.BasePath)
	}
	if len(spec.Schemes) != 1 || spec.Schemes[0] != "http" {
		t.Fatalf("schemes: got %v", spec.Schemes)
	}
}

func TestEnabled(t *testing.T) {
	viper.Reset()
	viper.Set("ENV", "prod")
	if Enabled() {
		t.Fatal("expected disabled in prod")
	}
	viper.Set("ENV", "dev")
	if !Enabled() {
		t.Fatal("expected enabled in dev")
	}
}
