package integration_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/unstoppablemango/terraform-provider-pfsense/mock"
)

func TestMain(m *testing.M) {
	srv := mock.NewServer()
	defer srv.Close()

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
