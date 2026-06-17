//ff:func feature=gate type=helper control=sequence level=error
//ff:what JSONL export sink을 생성한다. path의 부모 디렉터리가 있으면 미리 만들어 둔다(append 대상 파일은 첫 Emit 때 생성). reins cli.newJSONLSink 동등 구현.

package humaquest

import (
	"os"
	"path/filepath"
)

// newJSONLSink returns a sink writing to path, creating the parent directory. It
// mirrors reins' unexported cli.newJSONLSink so the cover command can build the
// export sink without reaching into reins internals.
func newJSONLSink(path string) (*jsonlSink, error) {
	if dir := filepath.Dir(path); dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, err
		}
	}
	return &jsonlSink{path: path}, nil
}
