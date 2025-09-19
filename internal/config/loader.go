package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func Load[T any](data string) (T, error) {
	var opts T

	if data == "" {
		return opts, nil
	}

	if _, err := os.Stat(data); err == nil {
		content, err := os.ReadFile(data)
		if err != nil {
			return opts, fmt.Errorf("failed to read config file: %w", err)
		}
		if err := json.Unmarshal(content, &opts); err != nil {
			return opts, fmt.Errorf("failed to parse config file: %w", err)
		}
		return opts, nil
	}

	if strings.HasPrefix(strings.TrimSpace(data), "{") {
		if err := json.Unmarshal([]byte(data), &opts); err != nil {
			return opts, fmt.Errorf("failed to parse config JSON string: %w", err)
		}
		return opts, nil
	}

	return opts, fmt.Errorf("invalid config input: not a file or JSON string")
}
