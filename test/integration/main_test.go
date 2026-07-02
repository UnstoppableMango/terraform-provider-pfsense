package integration_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/unstoppablemango/terraform-provider-pfsense/test/mock"
)

var mockServerURL string

func TestMain(m *testing.M) {
	srv := mock.NewServer()
	defer srv.Close()
	mockServerURL = srv.URL

	binaryPath := os.Getenv("PFSENSE_PROVIDER_BINARY")
	if binaryPath == "" {
		binaryPath = "result/bin/terraform-provider-pfsense"
	}
	rcPath := writeTerraformRC(filepath.Dir(binaryPath))
	defer os.Remove(rcPath)
	os.Setenv("TF_CLI_CONFIG_FILE", rcPath)

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
	f.Close()
	return f.Name()
}
