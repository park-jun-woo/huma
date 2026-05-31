//ff:func feature=hurlcheck type=helper control=sequence
//ff:what Flushes the current hurl entry and starts a new one from a method line
package hurlcheck

// startEntry flushes the current entry and begins a new one.
func (a *entryAccumulator) startEntry(line string) {
	a.flush()
	m := entryMethodRe.FindStringSubmatch(line)
	a.cur = &HurlEntry{Method: m[1], URL: m[2]}
}
