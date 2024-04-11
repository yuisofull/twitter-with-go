package tweetstorage

import "github.com/olivere/elastic/v7"

type esStore struct {
	client *elastic.Client
}

func NewESStore(client *elastic.Client) *esStore {
	return &esStore{client: client}
}
