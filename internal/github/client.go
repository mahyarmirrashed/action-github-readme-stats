package github

import (
	"net/http"

	"github.com/Khan/genqlient/graphql"
)

const githubGraphQLEndpoint = "https://api.github.com/graphql"

func NewClient(token string) graphql.Client {
	httpClient := &http.Client{
		Transport: &tokenTransport{
			underlying: http.DefaultTransport,
			token:      token,
		},
	}
	return graphql.NewClient(githubGraphQLEndpoint, httpClient)
}

type tokenTransport struct {
	underlying http.RoundTripper
	token      string
}

func (t *tokenTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+t.token)
	return t.underlying.RoundTrip(req)
}
