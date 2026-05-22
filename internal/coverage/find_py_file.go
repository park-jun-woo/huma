//ff:func feature=coverage type=helper control=iteration dimension=1
//ff:what Finds the coverage.py file entry matching the given handler file path
package coverage

// findPyFile searches the coverage report for a file matching the given handler file path.
func findPyFile(report coveragePyReport, handlerFile string) *coveragePyFile {
	for path, fd := range report.Files {
		if matchPyFile(path, handlerFile) {
			fd := fd
			return &fd
		}
	}
	return nil
}
