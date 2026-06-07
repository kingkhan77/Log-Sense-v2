package engine

import (
	"encoding/json"
	"strconv"

	"github.com/kingkhan77/log-sense/internal/models"
)

type RuleQuery struct {
	Level            string            `json:"level"`
	MessageContains  string            `json:"message_contains"`
	Fields           map[string]string `json:"fields"`
}

func buildCountQuery(rule models.AlertRule) (map[string]interface{}, error) {
	must := []interface{}{
		map[string]interface{}{
			"term": map[string]interface{}{
				"tenant_id": rule.TenantID,
			},
		},
		map[string]interface{}{
			"term": map[string]interface{}{
				"service_id": rule.ServiceID,
			},
		},
	}

	if len(rule.Query) > 0 {
		var rq RuleQuery
		if err := json.Unmarshal(rule.Query, &rq); err == nil {
			if rq.Level != "" {
				must = append(must, map[string]interface{}{
					"term": map[string]interface{}{
						"level": rq.Level,
					},
				})
			}
			if rq.MessageContains != "" {
				must = append(must, map[string]interface{}{
					"match": map[string]interface{}{
						"message": rq.MessageContains,
					},
				})
			}
			for k, v := range rq.Fields {
				must = append(must, map[string]interface{}{
					"term": map[string]interface{}{
						"metadata." + k: v,
					},
				})
			}
		}
	}

	return map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": must,
				"filter": []interface{}{
					map[string]interface{}{
						"range": map[string]interface{}{
							"timestamp": map[string]interface{}{
								"gte": "now-" + strconv.Itoa(rule.WindowMinutes) + "m",
							},
						},
					},
				},
			},
		},
	}, nil
}
