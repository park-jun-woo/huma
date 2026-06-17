//ff:func feature=gate type=helper control=sequence level=error
//ff:what reins applyVerdict+exportAndSave 꼬리를 export 프리미티브로 재현한다: quest.Apply(UTC RFC3339)→Save→quest.Export→Save(Export 실패여도 Emitted 래칫 보존). verdict를 래칫에 적용·영속화하는 단일 지점. PASS 잠금 권한은 게이트뿐 — 여기선 주어진 verdict만 적용한다.

package humaquest

import (
	"time"

	"github.com/park-jun-woo/reins/pkg/quest"
)

// commitVerdict ratchets a verdict onto an item and persists it, replicating reins'
// unexported applyVerdict → exportAndSave tail with exported primitives only:
// quest.Apply (UTC RFC3339) → Session.Save → quest.Export → Session.Save (so the
// Emitted ratchet survives an Export failure). It does not originate PASS — the
// caller supplies the gate verdict; this is purely the persistence/export tail.
func commitVerdict(s *quest.Session, it *quest.Item, v quest.Verdict, sink quest.Sink, sessionPath string) error {
	now := time.Now().UTC().Format(time.RFC3339)
	quest.Apply(it, v, now)
	if err := s.Save(sessionPath); err != nil {
		return err
	}
	_, exportErr := quest.Export(s, sink)
	saveErr := s.Save(sessionPath)
	if exportErr != nil {
		return exportErr
	}
	return saveErr
}
