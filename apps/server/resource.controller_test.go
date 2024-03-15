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
	search := SearchResourceRequest{
		Name: "Designing Data-Intensive Applications",
	}

	json_user, _ := json.Marshal(search)
	router := newAuthTestRouter(repos, nil)
	w := router.get("/resource/search", json_user)

	var actual Resource
	_ = json.Unmarshal(w.Body.Bytes(), &actual)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, resource.Url, actual.Url)
}
