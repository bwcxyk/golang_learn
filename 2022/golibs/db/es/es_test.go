/*
   Copyright (c) [2021] IT.SOS
   golibs is licensed under Mulan PSL v2.
   You can use this software according to the terms and conditions of the Mulan PSL v2.
   You may obtain a copy of Mulan PSL v2 at:
            http://license.coscl.org.cn/MulanPSL2
   THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
   See the Mulan PSL v2 for more details.
*/

package es

import (
	"bytes"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"golang.org/x/net/context"
	"log"
	"strconv"
	"strings"
	"sync"
	"testing"
)

func TestNewEs(t *testing.T) {
	var buf bytes.Buffer
	query := map[string]interface{}{
		"highlight": map[string]interface{}{
			"pre_tags":  []string{"<font-size='red'>"},
			"post_tags": []string{"</font>"},
			"fields": map[string]interface{}{
				"title": map[string]interface{}{},
				"intro": map[string]interface{}{},
				"data":  map[string]interface{}{},
			},
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		panic("Error encoding query:" + err.Error())
	}
	// Perform the search request.
	es := NewEs()
	res, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex("canal"),
		//es.Search.WithDocumentType("study_notes"),
		es.Search.WithFrom(0),
		es.Search.WithSize(10),
		//es.Search.WithSeqNoPrimaryTerm(true),
		es.Search.WithSourceExcludes("is_del"),
		es.Search.WithAnalyzer("ik_smart"),
		//es.Search.WithDocvalueFields(""),
		// https://lucene.apache.org/core/2_9_4/queryparsersyntax.html
		es.Search.WithQuery("第一天 AND is_del:1 AND is_state:2"),
		es.Search.WithDefaultOperator("and"),
		//es.Search.WithExplain(true),
		//es.Search.WithSuggestText()
		es.Search.WithBody(&buf),
		es.Search.WithTrackTotalHits(true),
		es.Search.WithErrorTrace(),
		es.Search.WithPretty(),
	)

	if err != nil {
		panic("Error getting response: " + err.Error())
	}
	defer res.Body.Close()
	print(res.String())
}

func TestNewEs2(t *testing.T) {
	es := NewEs()
	t.Run("查看es详细信息", func(t *testing.T) {
		res, err := es.Info()
		if err != nil {
			log.Fatalf("Error getting response: %s", err)
		}

		defer res.Body.Close()
		log.Println(res)
	})

	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	var (
		r  map[string]interface{}
		wg sync.WaitGroup
	)

	log.SetFlags(0)
	t.Run("存数据", func(t *testing.T) {
		// Initialize a client with the default settings.
		//
		// An `ELASTICSEARCH_URL` environment variable will be used when exported.
		//

		// 1. Get cluster info
		//
		res, err := es.Info()
		if err != nil {
			log.Fatalf("Error getting response: %s", err)
		}
		defer res.Body.Close()
		// Check response status
		if res.IsError() {
			log.Fatalf("Error: %s", res.String())
		}
		// Deserialize the response into a map.
		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			log.Fatalf("Error parsing the response body: %s", err)
		}
		// Print client and server version numbers.
		log.Printf("Client: %s", elasticsearch.Version)
		log.Printf("Server: %s", r["version"].(map[string]interface{})["number"])
		log.Println(strings.Repeat("~", 37))

		//2. Index documents concurrently

		for i, title := range []string{"中国人pk美国人", "中国人pk韩国人"} {
			wg.Add(1)

			go func(i int, title string) {
				defer wg.Done()

				// Build the request body.
				var b strings.Builder
				b.WriteString(`{"title" : "`)
				b.WriteString(title)
				b.WriteString(`"}`)

				// Set up the request object.
				req := esapi.IndexRequest{
					Index:      "canal",
					DocumentID: strconv.Itoa(i + 1),
					//OpType: "_doc",
					Body:    strings.NewReader(b.String()),
					Refresh: "true",
				}

				// Perform the request with the client.
				res, err := req.Do(context.Background(), es)
				if err != nil {
					log.Fatalf("Error getting response: %s", err)
				}
				defer res.Body.Close()

				if res.IsError() {
					log.Printf("[%s] Error indexing document ID=%d", res.Status(), i+1)
				} else {
					// Deserialize the response into a map.
					var r map[string]interface{}
					if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
						log.Printf("Error parsing the response body: %s", err)
					} else {
						// Print the response status and indexed document version.
						log.Printf("[%s] %s; version=%d", res.Status(), r["result"], int(r["_version"].(float64)))
					}
				}
			}(i, title)
		}
		wg.Wait()

		log.Println(strings.Repeat("-", 37))
	})
	t.Run("查数据", func(t *testing.T) {
		// 3. Search for the indexed documents
		//
		// Build the request body.
		var buf bytes.Buffer
		query := map[string]interface{}{
			"query": map[string]interface{}{
				"match": map[string]interface{}{
					"title": "中国",
				},
			},
		}
		if err := json.NewEncoder(&buf).Encode(query); err != nil {
			log.Fatalf("Error encoding query: %s", err)
		}

		// Perform the search request.
		res, err := es.Search(
			es.Search.WithContext(context.Background()),
			es.Search.WithIndex("canal"),
			es.Search.WithDocumentType("_doc"),
			es.Search.WithBody(&buf),
			es.Search.WithTrackTotalHits(true),
			es.Search.WithPretty(),
		)
		if err != nil {
			log.Fatalf("Error getting response: %s", err)
		}
		defer res.Body.Close()

		if res.IsError() {
			var e map[string]interface{}
			if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
				log.Fatalf("Error parsing the response body: %s", err)
			} else {
				// Print the response status and error information.
				log.Fatalf("[%s] %s: %s",
					res.Status(),
					e["error"].(map[string]interface{})["type"],
					e["error"].(map[string]interface{})["reason"],
				)
			}
		}

		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			log.Fatalf("Error parsing the response body: %s", err)
		}
		// Print the response status, number of results, and request duration.
		log.Printf(
			"[%s] %d hits; took: %dms",
			res.Status(),
			int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)),
			int(r["took"].(float64)),
		)
		// Print the ID and document source for each hit.
		for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
			log.Printf(" * ID=%s, %s", hit.(map[string]interface{})["_id"], hit.(map[string]interface{})["_source"])
		}

		log.Println(strings.Repeat("=", 37))
	})

}
