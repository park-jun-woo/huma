//ff:func feature=hurlcheck type=helper control=selection
//ff:what Grades and appends the current hurl entry, then clears it
package hurlcheck

// flush grades the current entry, appends it, and clears the cursor.
func (a *entryAccumulator) flush() {
	if a.cur == nil {
		return
	}
	a.cur.Grade = gradeEntry(*a.cur)
	a.entries = append(a.entries, *a.cur)
	a.cur = nil
}
