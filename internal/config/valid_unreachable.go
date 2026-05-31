//ff:func feature=config type=helper control=iteration dimension=1
//ff:what Filters reachability-exemption entries to those with both a reason and evidence
package config

// validUnreachable keeps only entries that carry both a reason and source
// evidence — the rest are unverifiable and therefore not valid exemptions.
func validUnreachable(raw []UnreachableEntry) []UnreachableEntry {
	out := make([]UnreachableEntry, 0, len(raw))
	for _, e := range raw {
		if e.Reason == "" || e.Evidence == "" {
			continue
		}
		out = append(out, e)
	}
	return out
}
