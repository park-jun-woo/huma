//ff:type feature=hurlcheck type=model
//ff:what entryAccumulator incrementally builds graded HurlEntry records while scanning a hurl file line by line
package hurlcheck

// entryAccumulator incrementally builds HurlEntry records as lines are scanned.
type entryAccumulator struct {
	entries []HurlEntry
	cur     *HurlEntry
}
