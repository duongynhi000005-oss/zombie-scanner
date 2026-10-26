package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAllRegionsOverridesConfigAndEnvRegion(t *testing.T) {
	cfgPath := filepath.Join(t.TempDir(), ".zombie-scanner.yaml")
	if err := os.WriteFile(cfgPath, []byte("region: ap-south-1\n"), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	t.Setenv("ZOMBIE_SCANNER_REGION", "eu-west-1")

	cmd := NewRootCommand("test", "abc1234")
	cmd.SetArgs([]string{"version", "--config", cfgPath, "--all-regions"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("expected --all-regions to succeed when config/env sets region, got error: %v", err)
	}

	gotRegion, err := cmd.PersistentFlags().GetString("region")
	if err != nil {
		t.Fatalf("get region flag: %v", err)
	}
	if gotRegion != "" {
		t.Fatalf("expected region flag to remain empty when --all-regions is set, got %q", gotRegion)
	}
}

func TestRegionAndAllRegionsFlagsRemainMutuallyExclusive(t *testing.T) {
	cmd := NewRootCommand("test", "abc1234")
	cmd.SetArgs([]string{"version", "--no-config", "--region", "us-east-1", "--all-regions"})
	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error when both --region and --all-regions are passed on the CLI")
	}
	if !strings.Contains(err.Error(), "if any flags in the group [region all-regions] are set none of the others can be") {
		t.Fatalf("unexpected error: %v", err)
	}
}
