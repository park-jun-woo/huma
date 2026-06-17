//ff:func feature=gate type=helper control=sequence
//ff:what huma 퀘스트 정의(humaDef)를 gate.Definition으로 반환해 cli.NewQuestCmd에 끼우는 Def() 생성자.

package humaquest

import (
	"github.com/park-jun-woo/reins/pkg/gate"
)

// Def returns the huma quest definition to wire into cli.NewQuestCmd.
func Def() gate.Definition { return humaDef{} }
