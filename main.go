package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic"
)

// Log json format
type Log struct {
	Level        int    `json:"level"`
	Message      string `json:"message"`
	Timestamp    int64  `json:"timestamp"`
	ContainerUID string `json:"containerUid"`
	GroupUID     string `json:"groupUid"`
	ProjectID    string `json:"projectId"`
	ServiceID    string `json:"serviceId"`
	ProjectUID   string `json:"projectUid"`
	ServiceUID   string `json:"serviceUid"`
	UID          string `json:"uid"`
}

var client elastic.Client

func main() {
	ctx := context.Background()

	client, err := elastic.NewClient(elastic.SetSniff(false))
	if err != nil {
		panic(err)
	}
	defer client.Stop()

	// Ping the Elasticsearch server to get e.g. the version number
	info, code, err := client.Ping("http://127.0.0.1:9200").Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)

	r := gin.Default()

	r.GET("/logs", getLogs)

	r.Run()
}

func getLogs(c *gin.Context) {
	ctx := context.Background()
	// termQuery := elastic.NewBoolQuery()

	client, err := elastic.NewClient(elastic.SetSniff(false))
	if err != nil {
		panic(err)
	}
	defer client.Stop()

	searchResult, err := client.Search().
		Index("logs"). // search in index "twitter"
		// Query(termQuery). // specify the query
		Do(ctx) // execute
	if err != nil {
		// Handle error
		panic(err)
	}

	// Here's how you iterate through results with full control over each step.
	if searchResult.Hits.TotalHits > 0 {
		fmt.Printf("Found a total of %d logs\n", searchResult.Hits.TotalHits)

		// Iterate through results
		for _, hit := range searchResult.Hits.Hits {
			// hit.Index contains the name of the index

			// Deserialize hit.Source into a Tweet (could also be just a map[string]interface{}).
			var t Log
			err := json.Unmarshal(*hit.Source, &t)
			if err != nil {
				// Deserialization failed
			}

			// Work with tweet
			fmt.Printf("Log by %d: %s\n", t.Timestamp, t.Message)
		}
	} else {
		// No hits
		fmt.Print("Found no logs\n")
	}

	c.JSON(200, gin.H{
		"message": "bar",
	})
}
