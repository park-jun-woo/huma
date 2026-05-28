//ff:func feature=adapter type=helper control=iteration dimension=1
//ff:what Collects covered and total line maps from a JaCoCo package for a matching handler file
package adapter

import "path/filepath"

// collectPackageLines iterates source files in a JaCoCo package and populates
// covered/total maps for lines matching the handler file and line range.
func collectPackageLines(pkg jacocoPackage, handlerFile string, startLine, endLine int, covered, total map[int]bool) {
	baseName := filepath.Base(handlerFile)
	for _, sf := range pkg.SourceFiles {
		if !matchSourceFile(sf.Name, baseName, handlerFile, pkg.Name) {
			continue
		}
		collectSourceFileLines(sf, startLine, endLine, covered, total)
	}
}
