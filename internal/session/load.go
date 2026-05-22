//ff:func feature=session type=reader control=sequence
//ff:what Loads session state from the .hurlfill/session.json file
package session

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func Load() (*Session, error) {
	data, err := os.ReadFile(filepath.Join(sessionDir, sessionFile))
	if err != nil {
		return nil, err
	}
	var s Session
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}
	return &s, nil
}
