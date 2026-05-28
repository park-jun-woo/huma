//ff:func feature=adapter type=engine control=sequence
//ff:what Launches the Java server process with JaCoCo agent injected via JAVA_TOOL_OPTIONS
package adapter

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Start launches the Java server process with JaCoCo agent injected
// via the JAVA_TOOL_OPTIONS environment variable.
func (a *JavaAdapter) Start() error {
	parts := strings.Fields(a.cfg.Start)
	if len(parts) == 0 {
		return fmt.Errorf("empty start command")
	}

	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr

	env := os.Environ()

	absDir, err := filepath.Abs(a.jacocoDir)
	if err != nil {
		return fmt.Errorf("abs jacoco dir: %w", err)
	}

	agentPath := resolveJacocoAgent(a.cfg.Env)
	if agentPath != "" {
		env = append(env, fmt.Sprintf(
			"JAVA_TOOL_OPTIONS=-javaagent:%s=destfile=%s/jacoco.exec,output=file",
			agentPath, absDir,
		))
	}

	for k, v := range a.cfg.Env {
		if k == "JACOCO_AGENT" {
			continue
		}
		env = append(env, k+"="+v)
	}
	cmd.Env = env

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start server: %w", err)
	}
	a.proc = cmd
	return nil
}
