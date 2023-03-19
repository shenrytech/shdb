package shdb

import (
	"context"
	"strings"
	"testing"

	"go.uber.org/zap"
)

func TestSearch(t *testing.T) {
	log, _ = zap.NewDevelopment()
	GenerateTestData(1)

	ctx := context.Background()

	var (
		nextPageToken string = ""
		res           *SearchResult
		err           error
	)
	allRes := &SearchResult{Hits: []*SearchHit{}}
	for {
		res, nextPageToken, err = Search(ctx, TObj, func(s string) bool {
			return strings.Contains(s, "Staffan Olsson")
		}, 10, "")

		if err != nil {
			t.Fail()
		}
		allRes.Hits = append(allRes.Hits, res.Hits...)
		if nextPageToken == "" {
			break
		}
	}

	if len(allRes.Hits) != 1 {
		t.Fail()
	}
}

func TestSearch2(t *testing.T) {
	log, _ = zap.NewDevelopment()
	count := 5000
	GenerateTestData(count)

	ctx := context.Background()

	var (
		nextPageToken string = ""
		res           *SearchResult
		err           error
	)
	allRes := &SearchResult{Hits: []*SearchHit{}}
	for {
		res, nextPageToken, err = Search(ctx, TObj, func(s string) bool {
			return strings.Contains(s, "Staffan Olsson")
		}, 100, nextPageToken)

		if err != nil {
			t.Fail()
		}
		allRes.Hits = append(allRes.Hits, res.Hits...)
		if nextPageToken == "" {
			break
		}
	}

	if len(allRes.Hits) != count {
		t.Fail()
	}
}
