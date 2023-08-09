package block_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/vjerci/reverse-proxy/internal/block"
)

func TestGuardsCollection(t *testing.T) {
	guards := []block.Guard{
		&block.MethodGuard{
			Method: http.MethodDelete,
		},
		&block.HeaderGuard{
			Header: "header",
			Value:  "header",
		},
		&block.QueryParamGuard{
			QueryParam: "queryParam",
			Value:      "queryParam",
		},
	}

	coll := block.NewGuardsCollection(guards)

	deleteRequest, err := http.NewRequest(http.MethodDelete, "", strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}

	headerRequest, err := http.NewRequest(http.MethodGet, "", strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}
	headerRequest.Header.Set("header", "header")

	queryParamRequest, err := http.NewRequest(http.MethodGet, "http://localhost/?queryParam=queryParam", strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}

	passingRequest, err := http.NewRequest(http.MethodGet, "", strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		testName string
		req      *http.Request
		block    bool
	}{
		{
			testName: "method_guard",
			req:      deleteRequest,
			block:    true,
		},
		{
			testName: "header_guard",
			req:      headerRequest,
			block:    true,
		},
		{
			testName: "query_param_guard",
			req:      queryParamRequest,
			block:    true,
		},
		{
			testName: "passing_request",
			req:      passingRequest,
			block:    false,
		},
	}

	for _, test := range testCases {
		block := coll.ShouldBlock(test.req)
		if block != test.block {
			t.Fatalf("%s test case failed, expected outcome %t", test.testName, test.block)
		}
	}
}

func TestGuardsJoiner(t *testing.T) {
	guards := []block.Guard{
		&block.MethodGuard{
			Method: http.MethodDelete,
		},
		&block.PathGuard{
			Path: "/api",
		},
	}

	joiner := block.NewGuardsJoiner(guards)

	blockingRequest, err := http.NewRequest(http.MethodDelete, "http://domain.com/api", strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}

	passingRequest, err := http.NewRequest(http.MethodGet, "http://domain.com/api", strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		testName string
		req      *http.Request
		block    bool
	}{
		{
			testName: "all_guards_block",
			req:      blockingRequest,
			block:    true,
		},
		{
			testName: "only_1_guard_blocks",
			req:      passingRequest,
			block:    false,
		},
	}

	for _, test := range testCases {
		block := joiner.ShouldBlock(test.req)
		if block != test.block {
			t.Fatalf("%s test case failed, expected outcome %t", test.testName, test.block)
		}
	}
}
