//ff:func feature=gate type=helper control=sequence level=error
//ff:what jsonlSink.Emit. 아이템 하나를 JSON으로 직렬화해 개행을 붙여 sink 파일에 append한다(파일이 없으면 생성). reins cli.jsonlSink.Emit 동등 구현.

package humaquest

import (
	"encoding/json"
	"os"

	"github.com/park-jun-woo/reins/pkg/quest"
)

// Emit appends one JSON-encoded item, newline-terminated, to the sink's file.
func (s *jsonlSink) Emit(it *quest.Item) error {
	f, err := os.OpenFile(s.path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	b, err := json.Marshal(it)
	if err != nil {
		return err
	}
	b = append(b, '\n')
	_, err = f.Write(b)
	return err
}
