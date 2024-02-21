package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stonear/go-template/client/pokemon"
	"github.com/stonear/go-template/response"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
)

type PokemonService interface {
	Index(c *gin.Context)
}

func NewPokemonService(
	log *otelzap.Logger,
	pokemonClient pokemon.Pokemon,
) PokemonService {
	return &pokemonService{
		log:           log,
		pokemonClient: pokemonClient,
	}
}

type pokemonService struct {
	log           *otelzap.Logger
	pokemonClient pokemon.Pokemon
}

func (s *pokemonService) Index(c *gin.Context) {
	ctx := c.Request.Context()
	pokemon, err := s.pokemonClient.GetPokemon(ctx)
	if err != nil {
		s.log.Ctx(ctx).Error(
			"failed to call pokemonClient.GetPokemon",
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, response.New(
			response.CodeGeneralError,
			nil,
		))
		return
	}
	c.JSON(http.StatusOK, response.New(
		response.CodeSuccess,
		pokemon,
	))
}
