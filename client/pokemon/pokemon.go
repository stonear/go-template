package pokemon

import (
	"context"
	"net/http"

	"github.com/stonear/go-template/library/httpclient"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
)

type Pokemon interface {
	GetPokemon(ctx context.Context) (*PokemonList, error)
}

func New(log *otelzap.Logger, client httpclient.Helper) Pokemon {
	return &pokemon{
		log:    log,
		client: client,
		// TODO: make this parameterizable
		baseUrl: "https://pokeapi.co/api/v2",
	}
}

type pokemon struct {
	log     *otelzap.Logger
	client  httpclient.Helper
	baseUrl string
}

func (c *pokemon) GetPokemon(ctx context.Context) (*PokemonList, error) {
	resp := new(PokemonList)
	fullPath := c.client.JoinPath(c.baseUrl, "pokemon")

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, fullPath, c.client.AnyToBuffer(nil))
	if err != nil {
		return nil, c.client.ErrorParseHTTPRequest(httpReq, err)
	}

	httpResp, err := c.client.Do(httpReq, &resp)
	if err != nil {
		return nil, c.client.ErrorParseHTTPResponse(httpResp, err)
	}

	return resp, nil
}
