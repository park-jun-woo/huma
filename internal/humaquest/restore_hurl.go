//ff:func feature=gate type=helper control=sequence level=error
//ff:what hurlStash를 디스크에 복원한다. 직전에 파일이 있었으면 그 내용을 되쓰고, 없었으면 생성물을 삭제해 무파일 상태로 되돌린다(사용자 자산 무손상). 비-PASS 종료/파싱 실패 경로에서 호출된다.

package humaquest

import "os"

// restoreHurl undoes a generate-mode overwrite from a hurlStash. If a prior file
// existed, its content is written back; otherwise the generated file is removed,
// returning the path to its no-file state. It is called on a non-PASS outcome or a
// parse failure so the user's prior .hurl is never lost.
func restoreHurl(st hurlStash) error {
	if !st.existed {
		if err := os.Remove(st.path); err != nil && !os.IsNotExist(err) {
			return err
		}
		return nil
	}
	return os.WriteFile(st.path, st.content, 0o644)
}
