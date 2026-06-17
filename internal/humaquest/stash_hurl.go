//ff:func feature=gate type=helper control=sequence level=error
//ff:what 정규 경로의 기존 .hurl 내용을 hurlStash로 읽어둔다(generate 덮어쓰기 전 보호). 파일이 없으면 existed=false인 빈 stash를 돌려준다(에러 아님). 진짜 읽기 에러만 전파.

package humaquest

import "os"

// stashHurl reads the current .hurl at path into a hurlStash so a generate-mode
// overwrite can be undone. A missing file is not an error: it yields an
// existed=false stash (restore will then remove the generated file). Only a
// genuine read failure propagates.
func stashHurl(path string) (hurlStash, error) {
	b, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return hurlStash{path: path, existed: false}, nil
	}
	if err != nil {
		return hurlStash{}, err
	}
	return hurlStash{path: path, content: b, existed: true}, nil
}
