package mjd

import (
	"context"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/get"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/index"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	config2 "github.com/manjada/com/config"
	"strings"
)

var client *elasticsearch.TypedClient

type NoSqlInterface interface {
	Create(document interface{}) (*index.Response, error)
	GetByDocId(documentId string) (*get.Response, error)
	IsDocumentExist(documentId string) (bool, error)
	Search(map[string]interface{}) (*search.Response, error)
}

type Elastic struct {
	Index   string
	Refresh string
}

func NewElastic() (*Elastic, error) {
	var err error
	if client == nil {
		config := config2.GetConfig().NoSqlConfig
		var address []string
		address = strings.Split(config.Host, ",")
		client, err = elasticsearch.NewTypedClient(elasticsearch.Config{Addresses: address, Username: config.User, Password: config.Pass})
		if err != nil {
			return &Elastic{}, err
		}
		config2.Info("ES initialized...")
	}

	return &Elastic{}, err
}

func (e *Elastic) Create(document interface{}) (*index.Response, error) {
	// Build the request body.
	res, err := client.Index(e.Index).
		Request(document).
		Do(context.Background())
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (e *Elastic) GetByDocId(documentId string) (*get.Response, error) {
	// Build the request body.
	res, err := client.Get(e.Index, "121212").Do(context.Background())
	if err != nil {
		config2.Error(err)
		return nil, err
	}
	return res, err
}

func (e *Elastic) IsDocumentExist(documentId string) (bool, error) {
	fmt.Println("index ", e.Index, documentId)
	if exists, err := client.Exists(e.Index, documentId).IsSuccess(nil); !exists {
		// The document exists !
		return false, err
	}
	return false, nil
}

func (e *Elastic) Search(data map[string]interface{}) (*search.Response, error) {
	res, err := client.Search().
		Index(e.Index).
		Request(&search.Request{
			Query: &types.Query{
				Match: map[string]types.MatchQuery{
					"Nama": {Query: "ocky"},
				},
			},
		}).Do(context.Background())
	return res, err
}
