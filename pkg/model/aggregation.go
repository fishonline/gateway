package model

import (
	"encoding/json"
	"io"
	"regexp"

	"github.com/valyala/fasthttp"
)

// Node aggregation node struct
type Node struct {
	ClusterName string `json:"clusterName,omitempty"`
	URL         string `json:"url,omitempty"`
	Rewrite     string `json:"rewrite,omitempty"`
	AttrName    string `json:"attrName,omitempty"`
}

// Aggregation aggregation struct
// a aggregation container a url and some nodes
type Aggregation struct {
	URL     string         `json:"url"`
	Nodes   []*Node        `json:"nodes"`
	Pattern *regexp.Regexp `json:"-"`
}

// UnMarshalAggregation unmarshal
func UnMarshalAggregation(data []byte) *Aggregation {
	v := &Aggregation{}
	json.Unmarshal(data, v)
	// v.URL, _ = url.QueryUnescape(v.URL)
	return v
}

// UnMarshalAggregationFromReader unmarshal from reader
func UnMarshalAggregationFromReader(r io.Reader) (*Aggregation, error) {
	v := &Aggregation{}

	decoder := json.NewDecoder(r)
	err := decoder.Decode(v)

	return v, err
}

// NewAggregation create a Aggregation
func NewAggregation(url string, nodes []*Node) *Aggregation {
	return &Aggregation{
		URL:   url,
		Nodes: nodes,
	}
}

// Marshal marshal
func (a *Aggregation) Marshal() []byte {
	v, _ := json.Marshal(a)
	return v
}

func (a *Aggregation) getNodeURL(req *fasthttp.Request, node *Node) string {
	if node.Rewrite == "" {
		return node.URL
	}

	return a.Pattern.ReplaceAllString(string(req.URI().RequestURI()), node.Rewrite)
}

func (a *Aggregation) matches(req *fasthttp.Request) bool {
	return a.Pattern.Match(req.URI().RequestURI())
}

func (a *Aggregation) updateFrom(ang *Aggregation) {
	a.URL = ang.URL
	a.Nodes = ang.Nodes
}
