package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/kingkhan77/log-sense/internal/middleware"
	"github.com/gin-gonic/gin"
	opensearch "github.com/opensearch-project/opensearch-go"
)

type LogController struct {
	os *opensearch.Client
}

func NewLogController(os *opensearch.Client) *LogController {
	return &LogController{os: os}
}

func (c *LogController) Search(ctx *gin.Context) {
	tenantID  := middleware.TenantID(ctx)
	serviceID := ctx.Query("service_id")
	message   := ctx.Query("message")
	level     := ctx.Query("level")
	from      := ctx.DefaultQuery("from", "now-1h")
	to        := ctx.DefaultQuery("to", "now")

	must := []interface{}{
		map[string]interface{}{"term": map[string]interface{}{"tenant_id": tenantID}},
	}
	if serviceID != "" {
		must = append(must, map[string]interface{}{"term": map[string]interface{}{"service_id": serviceID}})
	}
	if level != "" {
		must = append(must, map[string]interface{}{"term": map[string]interface{}{"level": level}})
	}
	if message != "" {
		must = append(must, map[string]interface{}{"match": map[string]interface{}{"message": message}})
	}

	limit := 50
	if l, err := strconv.Atoi(ctx.DefaultQuery("limit", "")); err == nil && l > 0 && l <= 200 {
		limit = l
	}
	offset := 0
	if o, err := strconv.Atoi(ctx.DefaultQuery("offset", "")); err == nil && o >= 0 {
		offset = o
	}

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": must,
				"filter": []interface{}{
					map[string]interface{}{
						"range": map[string]interface{}{
							"timestamp": map[string]interface{}{"gte": from, "lte": to},
						},
					},
				},
			},
		},
		"sort": []interface{}{map[string]interface{}{"timestamp": map[string]interface{}{"order": "desc"}}},
		"size": limit,
		"from": offset,
	}

	body, _ := json.Marshal(query)

	res, err := c.os.Search(
		c.os.Search.WithIndex("logs"),
		c.os.Search.WithBody(bytes.NewReader(body)),
		c.os.Search.WithContext(ctx.Request.Context()),
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "opensearch unavailable"})
		return
	}
	defer res.Body.Close()

	if res.IsError() {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "search error: " + res.Status()})
		return
	}

	var result struct {
		Hits struct {
			Total struct {
				Value int `json:"value"`
			} `json:"total"`
			Hits []struct {
				Source map[string]interface{} `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "decode failed"})
		return
	}

	logs := make([]map[string]interface{}, 0, len(result.Hits.Hits))
	for _, h := range result.Hits.Hits {
		logs = append(logs, h.Source)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"total": result.Hits.Total.Value,
		"logs":  logs,
	})
}
