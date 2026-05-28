//ff:func feature=adapter type=helper control=iteration dimension=1
//ff:what Resolves the JaCoCo agent JAR path from env, Maven repo, or project directories
package adapter

import (
	"os"
	"path/filepath"
	"sort"
)

// resolveJacocoAgent finds the JaCoCo agent JAR file.
// Search order:
// 1. cfg.Env["JACOCO_AGENT"] (user-specified)
// 2. Maven local repo: ~/.m2/repository/org/jacoco/org.jacoco.agent/*/org.jacoco.agent-*-runtime.jar
// 3. Project directories: build/jacoco/ or target/jacoco-agent/
func resolveJacocoAgent(env map[string]string) string {
	if p, ok := env["JACOCO_AGENT"]; ok && p != "" {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}

	home, err := os.UserHomeDir()
	if err == nil {
		pattern := filepath.Join(home, ".m2", "repository", "org", "jacoco",
			"org.jacoco.agent", "*", "org.jacoco.agent-*-runtime.jar")
		matches, _ := filepath.Glob(pattern)
		if len(matches) > 0 {
			sort.Strings(matches)
			return matches[len(matches)-1]
		}
	}

	projectDirs := []string{"build/jacoco", "target/jacoco-agent"}
	for _, dir := range projectDirs {
		pattern := filepath.Join(dir, "*.jar")
		matches, _ := filepath.Glob(pattern)
		for _, m := range matches {
			return m
		}
	}

	return ""
}
