module github.com/park-jun-woo/huma

go 1.23.0

require (
	github.com/park-jun-woo/reins v0.0.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/cobra v1.10.2 // indirect
	github.com/spf13/pflag v1.0.9 // indirect
)

// Local development wiring: resolve the reins framework and its toulmin
// dependency from their sibling checkouts instead of published modules.
replace github.com/park-jun-woo/reins => ../reins

replace github.com/park-jun-woo/toulmin => ../toulmin
