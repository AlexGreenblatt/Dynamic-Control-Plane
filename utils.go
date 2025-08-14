package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func readRouteConfigJSON() ([]RouteConfig, error) {
	file, err := os.ReadFile("route_config.json")
	if err != nil {
		return nil, fmt.Errorf("error reading config: %w", err)
	}

	var configs []RouteConfig
	if err := json.Unmarshal(file, &configs); err != nil {
		return nil, fmt.Errorf("error parsing JSON: %w", err)
	}
	return configs, nil
}

func readRegoPoliciesFromFiles(fileNames []string) ([]string, error) {
	var policies []string
	for _, file := range fileNames {
		fileData, err := os.ReadFile(file)
		if err != nil {
			return nil, fmt.Errorf("error reading policy file %s: %w", file, err)
		}
		policies = append(policies, string(fileData))
	}
	return policies, nil
}
