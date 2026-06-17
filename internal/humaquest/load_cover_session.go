//ff:func feature=gate type=helper control=sequence level=error
//ff:what cover용 세션 로드. quest.Load로 sessionPath를 읽되, 파일 부재는 명확한 에러로 환원한다 — cover는 scan으로 시드된 세션을 전제로 한다(submit/loop와 달리 빈 세션을 새로 만들지 않는다).

package humaquest

import (
	"fmt"
	"os"

	"github.com/park-jun-woo/reins/pkg/quest"
)

// loadCoverSession loads the seeded session at path. Unlike reins' loadSession (which
// creates a fresh empty session for the first scan), cover requires an existing
// seeded session, so a missing file is a clear, actionable error rather than a silent
// empty loop.
func loadCoverSession(path string) (*quest.Session, error) {
	s, err := quest.Load(path)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("no session at %s — run `huma scan` first", path)
	}
	if err != nil {
		return nil, err
	}
	return s, nil
}
