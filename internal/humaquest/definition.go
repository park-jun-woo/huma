//ff:type feature=gate type=model
//ff:what huma 퀘스트의 gate.Definition 구현체 humaDef 타입. Seed/Render/Prepare/Rules 4메서드는 각자 파일에 분리돼 있으며 Phase 001에서는 스텁 — Phase 002~005에서 본체가 채워진다. Def()로 cli.NewQuestCmd에 끼운다.

package humaquest

// humaDef is the huma quest's gate.Definition implementation. The four methods are
// stubs in Phase 001 (skeleton + library port); they are filled in later phases:
// Seed → Phase 002, Rules → Phase 002, Render → Phase 003, Prepare → Phase 004,
// and the CRI verdict Evaluator → Phase 005.
type humaDef struct{}
