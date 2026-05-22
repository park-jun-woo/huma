//ff:func feature=session type=engine control=sequence
//ff:what Persists session state to .hurlfill/session.json
package session

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func (s *Session) Save() error {
	if err := os.MkdirAll(sessionDir, 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(sessionDir, sessionFile), data, 0o644)
}
