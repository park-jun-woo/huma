//ff:type feature=gate type=model
//ff:what generate 모드에서 기존(사람이 쓴) .hurl을 보호하기 위한 stash. 생성물을 정규 경로에 덮어쓰기 전 직전 내용을 담아두고, 비-PASS면 복원·PASS면 폐기(생성물 확정)한다. existed=false면 직전에 파일이 없었다는 뜻이라 복원 시 삭제로 되돌린다.

package humaquest

// hurlStash captures the on-disk .hurl content (if any) that existed before a
// generate-mode overwrite, so a non-PASS attempt can restore the user's prior
// asset and a PASS keeps the generated file. When existed is false there was no
// prior file, so restore removes the generated one.
type hurlStash struct {
	path    string
	content []byte
	existed bool
}
