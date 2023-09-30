package pokemon

import (
	"net/http"

	"github.com/stonear/go-template/library/httpclient"
	"go.uber.org/zap"
)

type Pokemon interface {
	GetPokemon() (*PokemonList, error)
}

func New(log *zap.Logger, client httpclient.Helper) Pokemon {
	return &pokemon{
		log:    log,
		client: client,
		// TODO: make this parameterizable
		baseUrl: "https://pokeapi.co/api/v2",
	}
}

type pokemon struct {
	log     *zap.Logger
	client  httpclient.Helper
	baseUrl string
}

func (c *pokemon) GetPokemon() (*PokemonList, error) {
	resp := new(PokemonList)
	fullPath := c.client.JoinPath(c.baseUrl, "pokemon")

	httpReq, err := http.NewRequest(http.MethodGet, fullPath, c.client.AnyToBuffer(nil))
	if err != nil {
		return nil, c.client.ErrorParseHTTPRequest(httpReq, err)
	}

	httpResp, err := c.client.Do(httpReq, &resp)
	if err != nil {
		return nil, c.client.ErrorParseHTTPResponse(httpResp, err)
	}

	return resp, nil
}
