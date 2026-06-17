//ff:type feature=gate type=model
//ff:what cover 명령용 export sink. quest.Sink을 구현해 아이템 하나당 JSON 한 줄을 파일에 append(JSONL)한다. reins의 jsonlSink가 unexported라 동등 구현을 둔다 — quest.Export가 종단·미방출 아이템만 방출하므로 증분·멱등이다.

package humaquest

import "github.com/park-jun-woo/reins/pkg/quest"

// jsonlSink is the cover command's quest.Sink: it appends one JSON line per item to a
// file (JSONL). reins ships an equivalent sink in package cli, but it is unexported, so
// cover provides its own to replicate the PASS-lock export tail with exported
// primitives only (quest.Export + quest.Session.Save). Export is incremental and
// idempotent because quest.Export emits only terminal, not-yet-emitted items.
type jsonlSink struct {
	path string
}

var _ quest.Sink = (*jsonlSink)(nil)
