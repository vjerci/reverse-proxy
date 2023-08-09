// I would use open policy agent if this was real world example, but as far as i understand the goal is to use something simpler like json rules
package block

import (
	"net/http"
	"strings"
)

type Guard interface {
	ShouldBlock(req *http.Request) bool
}

type DecodedGuard interface {
	Guard
	IsValid() bool
}

type HeaderGuard struct {
	Header string `mapstructure:"header"`
	Value  string `mapstructure:"value"`
}

func (guard *HeaderGuard) ShouldBlock(req *http.Request) bool {
	return req.Header.Get(guard.Header) == guard.Value
}

func (guard *HeaderGuard) IsValid() bool {
	return guard.Header != "" && guard.Value != ""
}

type QueryParamGuard struct {
	QueryParam string `mapstructure:"query_param"`
	Value      string `mapstructure:"value"`
}

func (guard *QueryParamGuard) ShouldBlock(req *http.Request) bool {
	return req.URL.Query().Get(guard.QueryParam) == guard.QueryParam
}

func (guard *QueryParamGuard) IsValid() bool {
	return guard.QueryParam != "" && guard.Value != ""
}

type MethodGuard struct {
	Method string `mapstructure:"method"`
}

func (guard *MethodGuard) ShouldBlock(req *http.Request) bool {
	return req.Method == guard.Method
}

func (guard *MethodGuard) IsValid() bool {
	return guard.Method != ""
}

type PathGuard struct {
	Path string `mapstructure:"path"`
}

func (guard *PathGuard) ShouldBlock(req *http.Request) bool {
	return strings.Index(req.URL.Path, guard.Path) == 0
}

func (guard *PathGuard) IsValid() bool {
	return guard.Path != ""
}

// for guardscollecation if any guard blocks it results into blocking request
type GuardsCollection struct {
	guards []Guard
}

func NewGuardsCollection(guards []Guard) Guard {
	return &GuardsCollection{
		guards: guards,
	}
}

func (collection *GuardsCollection) ShouldBlock(req *http.Request) bool {
	for _, guard := range collection.guards {
		if guard.ShouldBlock(req) {
			return true
		}
	}

	return false
}

// for guards joiner each Guard must block to result into block request
type GuardsJoiner struct {
	guards []Guard
}

func NewGuardsJoiner(guards []Guard) Guard {
	return &GuardsJoiner{
		guards: guards,
	}
}

func (collection *GuardsJoiner) ShouldBlock(req *http.Request) bool {
	for _, guard := range collection.guards {
		if !guard.ShouldBlock(req) {
			return false
		}
	}

	return true
}
