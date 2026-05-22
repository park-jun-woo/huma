//ff:type feature=coverage type=model
//ff:what IstanbulRange represents a start/end position range in istanbul JSON format
package coverage

// istanbulRange represents a start/end range in istanbul JSON.
type istanbulRange struct {
	Start istanbulPosition `json:"start"`
	End   istanbulPosition `json:"end"`
}
