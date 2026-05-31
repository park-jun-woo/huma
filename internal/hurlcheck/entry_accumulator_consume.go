//ff:func feature=hurlcheck type=engine control=selection
//ff:what Consumes one hurl line, starting a new entry or updating the current entry's status/skip/asserts
package hurlcheck

import "strings"

// consume processes a single line, mutating the accumulator state.
func (a *entryAccumulator) consume(line string) {
	switch {
	case entryMethodRe.MatchString(line):
		a.startEntry(line)
	case a.cur == nil:
		return
	case httpStatusRe.MatchString(line):
		a.cur.Status = parseStatusLine(line)
	case skipOptionRe.MatchString(strings.TrimSpace(line)):
		a.cur.Skip = true
	case bodyAssertRe.MatchString(line):
		a.cur.Asserts++
	}
}
