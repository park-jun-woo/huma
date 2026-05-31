//ff:func feature=session type=helper control=sequence
//ff:what Reports whether a status (TODO/IMPROVE/UNVERIFIED) still requires work
package session

// needsWork reports whether a status represents an endpoint still needing work.
func needsWork(st Status) bool {
	return st == StatusTodo || st == StatusImprove || st == StatusUnverified
}
