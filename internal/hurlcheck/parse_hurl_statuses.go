//ff:func feature=hurlcheck type=parser control=iteration dimension=1
//ff:what Extracts all HTTP status code assertions from a hurl file
package hurlcheck

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

var httpStatusRe = regexp.MustCompile(`^HTTP\s+(\d{3})\s*$`)

// ParseHurlStatuses reads a .hurl file and extracts all HTTP <status> lines.
func ParseHurlStatuses(hurlPath string) ([]int, error) {
	f, err := os.Open(hurlPath)
	if err != nil {
		return nil, fmt.Errorf("open hurl file: %w", err)
	}
	defer f.Close()

	var statuses []int
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()
		m := httpStatusRe.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		code, err := strconv.Atoi(m[1])
		if err != nil {
			continue
		}
		statuses = append(statuses, code)
	}

	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("read hurl file: %w", err)
	}

	return statuses, nil
}
