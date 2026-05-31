//ff:func feature=hurlcheck type=helper control=sequence
//ff:what Flushes any pending entry and returns all accumulated graded entries
package hurlcheck

// finish flushes the final entry and returns all accumulated entries.
func (a *entryAccumulator) finish() []HurlEntry {
	a.flush()
	return a.entries
}
