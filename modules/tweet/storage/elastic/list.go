package tweetstorage

import (
	"context"
	"encoding/json"
	"github.com/olivere/elastic/v7"
	"log"
	"strings"
	"twitter/common"
	tweetmodel "twitter/modules/tweet/model"
)

func (s *esStore) ListTweetWithCondition(
	ctx context.Context,
	filter *tweetmodel.Filter,
	paging *common.Paging,
	moreKeys ...string,
) ([]tweetmodel.Tweet, error) {
	var results []tweetmodel.Tweet

	var empty []tweetmodel.Tweet

	q := elastic.NewBoolQuery()
	if f := filter; f != nil {
		if f.UserID > 0 {
			q = q.Must(elastic.NewMatchQuery("user_id", f.UserID))
		}

		search := strings.TrimSpace(f.Search)
		if len(search) > 0 {
			q = q.Must(elastic.NewFuzzyQuery("text_content", search))
		}
	}

	// q = q.Must(elastic.NewRangeQuery("shipping_fee_per_km").From(2).To(3))
	src, err := q.Source()
	if err != nil {
		log.Fatal(err)
	}

	data, err := json.Marshal(src)
	if err != nil {
		log.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	log.Println("got", got)

	offset := (paging.Page - 1) * paging.Limit

	searchResult, err := s.client.Search().
		Index(tweetmodel.Index). // search in index "twitter"
		Query(q).
		// specify the query
		Sort("id", false).                      // sort by "user" field, ascending
		From(offset).Size(paging.Limit).        // take documents 0-9
		Pretty(true).RestTotalHitsAsInt(false). // pretty print request and response JSON
		Do(ctx)                                 // execute
	if err != nil {
		// Handle error
		log.Println("Error in search: ", err)
		return empty, err
	}

	// TotalHits is another convenience function that works even when something goes wrong.

	// Here's how you iterate through results with full control over each step.
	if searchResult.TotalHits() > 0 {

		// Iterate through results
		for _, hit := range searchResult.Hits.Hits {
			// hit.Index contains the name of the index

			// Deserialize hit.Source into a Tweet (could also be just a map[string]interface{}).
			var tw tweetmodel.TweetES
			log.Println("hit.Source: ", string(hit.Source))
			err := json.Unmarshal(hit.Source, &tw)

			if err != nil {
				// Deserialization failed
				log.Println("Deserialization failed")
			}
			results = append(results, tw.ToTweet())

			// Work with tweet
			log.Printf("Tweet by %d: %s\n\n", tw.Id, tw.Text)
		}
	} else {
		// No hits
		log.Printf("Found no restaurant\n")
	}

	paging.Total = searchResult.TotalHits()
	return results, nil
}
