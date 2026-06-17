//ff:func feature=gate type=helper control=iteration dimension=1
//ff:what Definition.Seed. scanEndpoints로 OpenAPI/엣지함수에서 엔드포인트를 추출한 뒤, 엔드포인트마다 Key=ep.ID·State=TODO인 quest.Item을 만들고 SetPayload(ep)로 Endpoint를 Payload에 실어 반환한다. 세션 영속화는 reins가 한다.

package humaquest

import (
	"github.com/park-jun-woo/reins/pkg/quest"
)

// Seed scans endpoints from the OpenAPI/edge-function source given in args and
// creates one TODO quest.Item per endpoint, carrying the scanner.Endpoint as the
// item Payload (via SetPayload, never the field directly). reins persists the
// returned items; Seed does no session I/O of its own.
func (humaDef) Seed(args []string) ([]*quest.Item, error) {
	eps, err := scanEndpoints(args...)
	if err != nil {
		return nil, err
	}

	items := make([]*quest.Item, 0, len(eps))
	for i := range eps {
		it := &quest.Item{Key: eps[i].ID, State: quest.TODO}
		if err := it.SetPayload(eps[i]); err != nil {
			return nil, err
		}
		items = append(items, it)
	}
	return items, nil
}
