//ff:func feature=adapter type=engine control=sequence
//ff:what Runs jacococli XML report generation and parses JaCoCo XML for handler coverage
package adapter

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Collect generates a JaCoCo XML report from the exec file using jacococli,
// then parses the XML to extract line coverage for the handler file.
func (a *JavaAdapter) Collect(handlerFile string, startLine, endLine int) (*CoverageResult, error) {
	absDir, err := filepath.Abs(a.jacocoDir)
	if err != nil {
		return nil, fmt.Errorf("abs jacoco dir: %w", err)
	}

	execFile := filepath.Join(absDir, "jacoco.exec")
	xmlFile := filepath.Join(absDir, "jacoco.xml")

	if _, err := os.Stat(execFile); os.IsNotExist(err) {
		return nil, nil
	}

	classesDir := "target/classes"
	sourceDir := "src/main/java"

	out, err := exec.Command("java", "-jar", findJacocoCLI(),
		"report", execFile,
		"--classfiles", classesDir,
		"--sourcefiles", sourceDir,
		"--xml", xmlFile,
	).CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("jacococli report: %w\n%s", err, string(out))
	}

	covered, total, err := parseJacocoXML(xmlFile, handlerFile, startLine, endLine)
	if err != nil {
		return nil, fmt.Errorf("parse jacoco xml: %w", err)
	}

	totalCount := len(total)
	coveredCount := len(covered)

	var pct float64
	if totalCount > 0 {
		pct = float64(coveredCount) / float64(totalCount) * 100
	}

	uncovered, err := readUncoveredLines(handlerFile, covered, total)
	if err != nil {
		uncovered = nil
	}

	return &CoverageResult{
		Covered:   coveredCount,
		Total:     totalCount,
		Percent:   pct,
		Uncovered: uncovered,
	}, nil
}

// findJacocoCLI locates the JaCoCo CLI JAR file.
func findJacocoCLI() string {
	home, err := os.UserHomeDir()
	if err == nil {
		pattern := filepath.Join(home, ".m2", "repository", "org", "jacoco",
			"org.jacoco.cli", "*", "org.jacoco.cli-*-nodeps.jar")
		matches, _ := filepath.Glob(pattern)
		if len(matches) > 0 {
			return matches[len(matches)-1]
		}
	}
	return "jacococli.jar"
}
