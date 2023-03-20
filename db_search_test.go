// Copyright 2023 Shenry Tech AB
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package shdb

import (
	"context"
	"strings"
	"testing"
)

func TestSearch(t *testing.T) {
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
	count := 50
	pageSize := 10
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
		}, int32(count/pageSize), nextPageToken)

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
