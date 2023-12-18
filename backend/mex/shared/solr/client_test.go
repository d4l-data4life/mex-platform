package solr

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

func TestClient_DoJsonQuery(t *testing.T) {
	var foundVal uint32 = 1
	var startVal uint32 = 0
	scoreVal := 3.14
	docsVal := []GenericObject{
		{
			"title":    "Some doc",
			"category": "book"},
		{
			"title":    "Top speeds of sloths",
			"category": "study",
		}}
	solrResponseObj := QueryResponse{
		Response: QueryResult{
			NumFound:      foundVal,
			NumFoundExact: true,
			Start:         startVal,
			MaxScore:      scoreVal,
			Docs:          docsVal,
		},
		Facets: nil,
		// Highlighting: nil,
	}
	solrResponseBodyBytes, _ := json.Marshal(solrResponseObj)
	solrResponseBody := string(solrResponseBodyBytes)
	type solrResp struct {
		statusCode int
		body       string
	}
	type solrQueryParams struct {
		urlParams url.Values
		body      *QueryBody
	}
	getAllQuery := "*:*"
	testQueryParams := url.Values{}
	testQueryParams.Add("hl", "true")
	testQueryBody := &QueryBody{
		Query: getAllQuery,
	}
	tests := []struct {
		name            string
		solrResponse    solrResp
		solrQueryParams solrQueryParams
		wantResponse    *QueryResponse
		wantStatusCode  int
		wantErr         bool
	}{
		{
			name: "Returns if the body can be parsed, the parsed response and the olr response code is returned",
			solrQueryParams: solrQueryParams{
				urlParams: url.Values{},
				body:      testQueryBody,
			},
			solrResponse: solrResp{
				987,
				solrResponseBody,
			},
			wantResponse:   &solrResponseObj,
			wantStatusCode: 987,
		},
		{
			name: "Returns error if Solr responds with empty JSON body",
			solrQueryParams: solrQueryParams{
				urlParams: url.Values{},
				body:      testQueryBody,
			},
			solrResponse: solrResp{
				http.StatusOK,
				"",
			},
			wantErr: true,
		},
		{
			name: "Returns error if Solr responds with un-" +
				"parseable JSON body",
			solrQueryParams: solrQueryParams{
				urlParams: url.Values{},
				body:      testQueryBody,
			},
			solrResponse: solrResp{
				http.StatusOK,
				"{\"notThere\": \"Never seen this...\"",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange - set up mock server to return the required value
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.solrResponse.statusCode)
				if tt.solrResponse.body != "" {
					fmt.Fprintln(w, tt.solrResponse.body)
				}
			}))
			defer ts.Close()
			// Act
			solrApi := NewClient(ts.URL, "dummy")
			gotResp, gotStatusCode, err := solrApi.DoJSONQuery(context.TODO(), tt.solrQueryParams.urlParams, tt.solrQueryParams.body)
			// Assert
			if (err != nil) != tt.wantErr {
				t.Errorf("DoJsonQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotStatusCode != tt.wantStatusCode {
				t.Errorf("DoJsonQuery(): received staus code %d, wanted %d", gotStatusCode, tt.wantStatusCode)
			}
			if !reflect.DeepEqual(gotResp, tt.wantResponse) {
				t.Errorf("DoJsonQuery(): received response object = %v, wanted %v", gotResp, tt.wantResponse)
			}
		})
	}
}
