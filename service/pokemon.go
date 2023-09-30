package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stonear/go-template/client/pokemon"
	"github.com/stonear/go-template/response"
	"go.uber.org/zap"
)

type PokemonService interface {
	Index(ctx *gin.Context)
}

func NewPokemonService(
	log *zap.Logger,
	pokemonClient pokemon.Pokemon,
) PokemonService {
	return &pokemonService{
		log:           log,
		pokemonClient: pokemonClient,
	}
}

type pokemonService struct {
	log           *zap.Logger
	pokemonClient pokemon.Pokemon
}

func (s *pokemonService) Index(ctx *gin.Context) {
	pokemon, err := s.pokemonClient.GetPokemon()
	if err != nil {
		s.log.Error(
			"failed to call pokemonClient.GetPokemon",
			zap.Error(err),
		)
		ctx.JSON(http.StatusInternalServerError, response.New(
			response.CodeGeneralError,
			nil,
		))
		return
	}
	ctx.JSON(http.StatusOK, response.New(
		response.CodeSuccess,
		pokemon,
	))
}
