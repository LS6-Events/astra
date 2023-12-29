package comparison

import (
	"github.com/google/go-cmp/cmp"
	"gopkg.in/yaml.v3"
	"os"
	"testing"
)

func readFile(t *testing.T, path string) []byte {
	t.Helper()

	contents, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file %s: %v", path, err)
	}

	return contents
}

func CompareYAML(t *testing.T, snapshotPath, outputPath string) {
	t.Helper()

	outputPathContents := readFile(t, outputPath)
	snapshotPathContents := readFile(t, snapshotPath)

	var snapshotMap map[any]any
	var outputMap map[any]any

	err := yaml.Unmarshal(snapshotPathContents, &snapshotMap)
	if err != nil {
		t.Fatalf("failed to unmarshal snapshot yaml: %v", err)
	}

	err = yaml.Unmarshal(outputPathContents, &outputMap)
	if err != nil {
		t.Fatalf("failed to unmarshal output yaml: %v", err)
	}

	if !cmp.Equal(outputMap, snapshotMap) {
		t.Errorf("snapshots do not match: %s", cmp.Diff(snapshotMap, outputMap))
	}
}
