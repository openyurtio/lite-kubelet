package manifest

import "testing"

func TestGetPodsManifestPath(t *testing.T) {
	p := GetPodsManifestPath()
	t.Logf("path, %v\n", p)

	p = GetPodsManifestPath()
	t.Logf("path, %v\n", p)
}
