package manifest

import (
	"os"
	"path/filepath"

	"k8s.io/klog/v2"
)

const MQTTManifestRootPath = "/etc/kubernetes/mqttManifests/"

func GetPodsManifestPath() string {
	return filepath.Join(MQTTManifestRootPath, "pods")
}

func GetNodesManifestPath() string {
	return filepath.Join(MQTTManifestRootPath, "nodes")
}

func GetLeasesManifestPath() string {
	return filepath.Join(MQTTManifestRootPath, "leases")
}

func GetEventsManifestPath() string {
	return filepath.Join(MQTTManifestRootPath, "events")
}

func MkdirAllSubManifestPath() {
	allManifestPath := make([]string, 0, 5)
	allManifestPath = append(allManifestPath,
		GetPodsManifestPath(),
		GetNodesManifestPath(),
	)

	for _, p := range allManifestPath {
		if err := os.MkdirAll(p, os.ModePerm); err != nil {
			klog.Fatalf("Failed to create mqtt manifest source directory[%s]: %v", p, err)
		}
	}
}
