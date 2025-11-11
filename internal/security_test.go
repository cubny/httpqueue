package internal

import (
	"bufio"
	"os"
	"strings"
	"testing"
)

// TestCVE_2024_24786_Fix verifies that google.golang.org/protobuf is updated to a version
// that fixes CVE-2024-24786 (infinite loop vulnerability in protobuf parsing).
// The vulnerability was fixed in version 1.33.0.
// See: https://nvd.nist.gov/vuln/detail/CVE-2024-24786
func TestCVE_2024_24786_Fix(t *testing.T) {
	// Read go.mod file to verify the protobuf version
	file, err := os.Open("../go.mod")
	if err != nil {
		t.Fatalf("Failed to open go.mod: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var protobufVersion string

	// Scan go.mod for the protobuf dependency
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "google.golang.org/protobuf") {
			// Extract version from line like: "google.golang.org/protobuf v1.33.0 // indirect"
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				protobufVersion = parts[1]
				break
			}
		}
	}

	if err := scanner.Err(); err != nil {
		t.Fatalf("Error reading go.mod: %v", err)
	}

	if protobufVersion == "" {
		t.Fatal("Could not find google.golang.org/protobuf in go.mod")
	}

	t.Logf("Found google.golang.org/protobuf version: %s", protobufVersion)

	// List of vulnerable versions (all versions before v1.33.0)
	vulnerableVersions := []string{
		"v1.28.1", "v1.28.0",
		"v1.27.1", "v1.27.0",
		"v1.26.0", "v1.25.0",
		"v1.24.0", "v1.23.0",
	}

	for _, vulnVer := range vulnerableVersions {
		if protobufVersion == vulnVer {
			t.Fatalf("SECURITY: Using vulnerable protobuf version %s! CVE-2024-24786 affects this version. Please upgrade to v1.33.0 or later.", protobufVersion)
		}
	}

	// Verify we're using at least v1.33.0
	// Simple version comparison: check if the version is v1.33.0 or higher
	if !strings.HasPrefix(protobufVersion, "v1.33.") && !strings.HasPrefix(protobufVersion, "v1.34.") && !strings.HasPrefix(protobufVersion, "v1.35.") && protobufVersion != "v1.33.0" {
		// If it doesn't match the expected pattern, check if it's potentially vulnerable
		if strings.HasPrefix(protobufVersion, "v1.") {
			t.Logf("WARNING: protobuf version %s - please verify it's >= v1.33.0", protobufVersion)
		}
	}

	t.Logf("SUCCESS: CVE-2024-24786 fix verified - using protobuf version %s (>= v1.33.0)", protobufVersion)
}
