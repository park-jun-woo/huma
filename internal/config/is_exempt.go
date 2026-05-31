//ff:func feature=config type=helper control=iteration dimension=1
//ff:what Reports whether a given endpoint+status branch has a valid unreachable.yaml exemption
package config

// IsExempt reports whether a (method+path, status) branch has a valid
// reachability exemption. The endpoint key is matched against "<METHOD> <PATH>".
func IsExempt(entries []UnreachableEntry, endpoint string, status int) bool {
	for _, e := range entries {
		if e.Endpoint == endpoint && e.Status == status {
			return true
		}
	}
	return false
}
