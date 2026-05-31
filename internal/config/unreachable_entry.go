//ff:type feature=config type=model
//ff:what UnreachableEntry is a verifiable reachability-exemption artifact (reason + source evidence) for a branch
package config

// UnreachableEntry is a verifiable reachability-exemption artifact (§3.7).
// A missing/unbindable branch is only exempted (toward DONE, or toward
// branch-line binding) when the user records a reason backed by source
// evidence — silence is never an exemption.
type UnreachableEntry struct {
	Endpoint string `yaml:"endpoint"`
	Status   int    `yaml:"status"`
	Reason   string `yaml:"reason"`
	Evidence string `yaml:"evidence"`
}
