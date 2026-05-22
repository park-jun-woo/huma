//ff:type feature=coverage type=model
//ff:what IstanbulPosition represents a line and column position in istanbul JSON format
package coverage

// istanbulPosition represents a line/column position in istanbul JSON.
type istanbulPosition struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}
