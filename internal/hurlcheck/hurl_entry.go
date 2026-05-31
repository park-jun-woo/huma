//ff:type feature=hurlcheck type=model
//ff:what HurlEntry is one request block in a .hurl file with its measured assertion depth (A-grade)
package hurlcheck

// HurlEntry is one request block in a .hurl file with its measured assertion
// depth (A-grade per §3.3).
type HurlEntry struct {
	Method  string
	URL     string
	Status  int  // 0 if no HTTP status line
	Skip    bool // [Options] skip: true
	Asserts int  // number of jsonpath/header/body assertions
	Grade   int  // 0..3 assertion depth
}
