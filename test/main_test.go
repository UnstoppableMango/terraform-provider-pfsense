package integration_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestMain(m *testing.M) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()
	if err := os.Setenv("PFSENSE_MOCK_URL", srv.URL); err != nil {
		panic(err)
	}

	binaryPath := os.Getenv("PFSENSE_PROVIDER_BINARY")
	if binaryPath == "" {
		binaryPath = "result/bin/terraform-provider-pfsense"
	}
	absBinaryPath, err := filepath.Abs(binaryPath)
	if err != nil {
		panic(err)
	}
	rcPath := writeTerraformRC(filepath.Dir(absBinaryPath))
	defer os.Remove(rcPath)
	if err := os.Setenv("TF_CLI_CONFIG_FILE", rcPath); err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

func writeTerraformRC(binDir string) string {
	content := fmt.Sprintf(`
provider_installation {
  dev_overrides {
    "registry.terraform.io/unstoppablemango/pfsense" = %q
  }
  direct {}
}
`, binDir)
	f, err := os.CreateTemp("", "terraformrc-*.tfrc")
	if err != nil {
		panic(err)
	}
	if _, err := f.WriteString(content); err != nil {
		panic(err)
	}
	if err := f.Close(); err != nil {
		panic(err)
	}
	return f.Name()
}
