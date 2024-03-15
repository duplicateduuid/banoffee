package main

import (
	"encoding/json"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func deleteResourceByUrl(url string) func(*sqlx.DB) {
	return func(db *sqlx.DB) {
		db.Exec(`DELETE FROM "resource" WHERE url = $1`, url)
	}
}

type searchResponse struct {
	Resource Resource `json:"resource"`
}

func TestResourceSearch(t *testing.T) {
	t.Parallel()

	resource := NewResource("http://example.com/", "Designing Data-Intensive Applications", nil, nil, nil)

	repos := newTestRepositories(deleteResourceByUrl(resource.Url), t)
	repos.resourceRepository.CreateResource(resource)

	router := newAuthTestRouter(repos, nil)
	query := map[string]string{
		"name": resource.Name,
	}
	w := router.get("/resource/search", query)

	resp := w.Body.String()

	var actual Resource
	json.Unmarshal([]byte(resp), &actual)

	assert.Equal(t, 200, w.Code, resp)
	assert.Equal(t, resource.Url, actual.Url)
}
