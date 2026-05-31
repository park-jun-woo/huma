//ff:func feature=hurlcheck type=parser control=iteration dimension=1
//ff:what Reads a hurl file and splits it into graded request entries
package hurlcheck

import (
	"bufio"
	"fmt"
	"os"
)

// ParseHurlEntries reads a .hurl file and splits it into request entries,
// measuring each entry's assertion depth.
func ParseHurlEntries(hurlPath string) ([]HurlEntry, error) {
	f, err := os.Open(hurlPath)
	if err != nil {
		return nil, fmt.Errorf("open hurl file: %w", err)
	}
	defer f.Close()

	acc := &entryAccumulator{}
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for sc.Scan() {
		acc.consume(sc.Text())
	}
	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("read hurl file: %w", err)
	}
	return acc.finish(), nil
}
