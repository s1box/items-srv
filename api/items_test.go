package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/s1box/dc-go-srv/db"
	"github.com/stretchr/testify/assert"
)

type fakeDbClient struct {
	items  []db.Item
	nextId int64
}

func NewFakeDbClient() *fakeDbClient {
	return &fakeDbClient{nextId: 1}
}

func (c *fakeDbClient) SelectAllItems() ([]db.Item, error) {
	return c.items, nil
}
func (c *fakeDbClient) GetItemById(id string) (*db.Item, error) {
	for _, item := range c.items {
		if fmt.Sprint(item.Id) == id {
			return &item, nil
		}
	}
	return nil, nil
}
func (c *fakeDbClient) GetRandomItem() (*db.Item, error) {
	if len(c.items) > 0 {
		return &c.items[0], nil
	}
	return nil, nil
}
func (c *fakeDbClient) InsertItem(item db.Item) (string, error) {
	item.Id = c.nextId
	c.nextId++
	c.items = append(c.items, item)
	return fmt.Sprint(item.Id), nil
}
func (c *fakeDbClient) DeleteItemById(id string) (bool, error) {
	for i, item := range c.items {
		if fmt.Sprint(item.Id) == id {
			c.items = append(c.items[:i], c.items[i+1:]...)
			return true, nil
		}
	}
	return false, nil
}
func (c *fakeDbClient) Reset() {
	c.items = nil
	c.nextId = 1
}

var _ db.ItemsDbClient = (*fakeDbClient)(nil)

func TestGetItems(t *testing.T) {
	fakeDbClient := NewFakeDbClient()

	fakeDbClient.InsertItem(db.Item{Name: "test1", Num: 1.1})
	fakeDbClient.InsertItem(db.Item{Name: "test2", Num: 2.2})

	itemsApi := &itemsApiImpl{dbClient: fakeDbClient}

	testRouter := gin.Default()
	testRouter.GET("/items", itemsApi.GetItems)
	req, _ := http.NewRequest("GET", "/items", nil)
	reqRecoder := httptest.NewRecorder()
	testRouter.ServeHTTP(reqRecoder, req)

	assert.Equal(t, http.StatusOK, reqRecoder.Code)
	assert.Equal(t, "*", reqRecoder.Header()["Access-Control-Allow-Origin"][0])

	var items []db.Item
	json.Unmarshal(reqRecoder.Body.Bytes(), &items)
	assert.Len(t, items, 2)
}

func TestGetItem(t *testing.T) {
	fakeDbClient := NewFakeDbClient()

	fakeDbClient.InsertItem(db.Item{Name: "test1", Num: 1.1})
	fakeDbClient.InsertItem(db.Item{Name: "test2", Num: 2.2})

	itemsApi := &itemsApiImpl{dbClient: fakeDbClient}

	testRouter := gin.Default()
	testRouter.GET("/item/:id", itemsApi.GetItemById)
	req, _ := http.NewRequest("GET", "/item/2", nil)
	reqRecoder := httptest.NewRecorder()
	testRouter.ServeHTTP(reqRecoder, req)

	assert.Equal(t, http.StatusOK, reqRecoder.Code)
	assert.Equal(t, "*", reqRecoder.Header()["Access-Control-Allow-Origin"][0])

	var item db.Item
	json.Unmarshal(reqRecoder.Body.Bytes(), &item)
	assert.Equal(t, int64(2), item.Id)
	assert.Equal(t, "test2", item.Name)
	assert.Equal(t, 2.2, item.Num)
}

func TestGetItemNotFound(t *testing.T) {
	fakeDbClient := NewFakeDbClient()

	fakeDbClient.InsertItem(db.Item{Name: "test1", Num: 1.1})
	fakeDbClient.InsertItem(db.Item{Name: "test2", Num: 2.2})

	itemsApi := &itemsApiImpl{dbClient: fakeDbClient}

	testRouter := gin.Default()
	testRouter.GET("/item/:id", itemsApi.GetItemById)
	req, _ := http.NewRequest("GET", "/item/5", nil)
	reqRecoder := httptest.NewRecorder()
	testRouter.ServeHTTP(reqRecoder, req)

	assert.Equal(t, http.StatusNotFound, reqRecoder.Code)
	assert.Equal(t, "*", reqRecoder.Header()["Access-Control-Allow-Origin"][0])
}

func TestPostItem(t *testing.T) {
	fakeDbClient := NewFakeDbClient()

	fakeDbClient.InsertItem(db.Item{Name: "test1", Num: 1.1})
	fakeDbClient.InsertItem(db.Item{Name: "test2", Num: 2.2})

	itemsApi := &itemsApiImpl{dbClient: fakeDbClient}

	testRouter := gin.Default()
	testRouter.POST("/items", itemsApi.PostItem)
	req, _ := http.NewRequest("POST", "/items", strings.NewReader(`{"Name": "test3", "Num": 3.3}`))
	reqRecoder := httptest.NewRecorder()
	testRouter.ServeHTTP(reqRecoder, req)

	assert.Equal(t, http.StatusCreated, reqRecoder.Code)
	assert.Equal(t, "*", reqRecoder.Header()["Access-Control-Allow-Origin"][0])

	repJson := struct {
		Id string `json:"id"`
	}{}
	json.Unmarshal(reqRecoder.Body.Bytes(), &repJson)
	assert.Equal(t, "3", repJson.Id)
}

func TestDeleteItemById(t *testing.T) {
	fakeDbClient := NewFakeDbClient()

	fakeDbClient.InsertItem(db.Item{Name: "test1", Num: 1.1})
	fakeDbClient.InsertItem(db.Item{Name: "test2", Num: 2.2})

	itemsApi := &itemsApiImpl{dbClient: fakeDbClient}

	testRouter := gin.Default()
	testRouter.DELETE("/items/:id", itemsApi.DeleteItemById)
	req, _ := http.NewRequest("DELETE", "/items/2", nil)
	reqRecoder := httptest.NewRecorder()
	testRouter.ServeHTTP(reqRecoder, req)

	assert.Equal(t, http.StatusNoContent, reqRecoder.Code)
	assert.Equal(t, "*", reqRecoder.Header()["Access-Control-Allow-Origin"][0])
}

func TestDeleteItemByIdNotFound(t *testing.T) {
	fakeDbClient := NewFakeDbClient()

	fakeDbClient.InsertItem(db.Item{Name: "test1", Num: 1.1})
	fakeDbClient.InsertItem(db.Item{Name: "test2", Num: 2.2})

	itemsApi := &itemsApiImpl{dbClient: fakeDbClient}

	testRouter := gin.Default()
	testRouter.DELETE("/items/:id", itemsApi.DeleteItemById)
	req, _ := http.NewRequest("DELETE", "/items/5", nil)
	reqRecoder := httptest.NewRecorder()
	testRouter.ServeHTTP(reqRecoder, req)

	assert.Equal(t, http.StatusNotFound, reqRecoder.Code)
	assert.Equal(t, "*", reqRecoder.Header()["Access-Control-Allow-Origin"][0])
}
