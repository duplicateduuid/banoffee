package main

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func testNewResource(t *testing.T, repos *Repositories) *Resource {
	uuid := uuid.New().String()

	resource := NewResource("http://"+uuid+"/", uuid, nil, nil, nil)
	err := repos.resourceRepository.CreateResource(resource)

	if err != nil {
		t.Errorf("[ERROR] [testNewResource] failed to create resource: %s", err)
	}

	return resource
}

func TestCreateResource(t *testing.T) {
	t.Parallel()

	repos := newTestRepositories(t)

	uuid := uuid.New().String()
	resource := NewResource("http://"+uuid+"/", uuid, nil, nil, nil)
	auth, _ := testNewUser(t, &repos)

	router := newAuthTestRouter(t, repos, *auth)
	w := router.post("/resource", resource)

	assert.Equal(t, 204, w.Code, w.Body)
}

func TestResourceSearch(t *testing.T) {
	t.Parallel()

	repos := newTestRepositories(t)

	resource := testNewResource(t, &repos)
	user, _ := testNewUser(t, &repos)

	router := newAuthTestRouter(t, repos, *user)
	query := map[string]string{
		"name":   resource.Name,
		"limit":  "10",
		"offset": "0",
	}
	w := router.get("/resource/search", query)

	body := w.Body.String()

	var resp SearchResourceResponse
	json.Unmarshal([]byte(body), &resp)

	actual := resp.Resources[0]

	assert.Equal(t, 200, w.Code, resp)
	assert.Equal(t, resource.Url, actual.Url)
}
