//ff:func feature=config type=reader control=sequence
//ff:what Reads and validates .huma/unreachable.yaml, dropping entries without a reason and evidence
package config

import (
	"errors"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// LoadUnreachable reads .huma/unreachable.yaml. A missing file is not an error
// (returns nil). Entries lacking a reason or evidence are dropped — they do not
// constitute a valid exemption.
func LoadUnreachable() ([]UnreachableEntry, error) {
	path := filepath.Join(".huma", "unreachable.yaml")
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var raw []UnreachableEntry
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, err
	}
	return validUnreachable(raw), nil
}
