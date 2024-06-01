package pokemon

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stonear/go-template/client/pokemon"
	"github.com/stonear/go-template/library/redis"
	"github.com/stonear/go-template/response"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
)

type Service interface {
	Index(c *gin.Context)
}

func New(
	log *otelzap.Logger,
	cache *redis.Client,
	pokemonClient pokemon.Pokemon,
) Service {
	return &service{
		log:           log,
		cache:         cache,
		pokemonClient: pokemonClient,
	}
}

type service struct {
	log           *otelzap.Logger
	cache         *redis.Client
	pokemonClient pokemon.Pokemon
}

// @Summary		Pokemon List
// @Description	get pokemon list
// @Tags		pokemon
// @Accept		json
// @Produce		json
// @Success		200			{object}	response.Response{data=[]pokemon.PokemonList}
// @Failure		500			{object}	response.Response{data=nil}
// @Router		/pokemon	[get]
func (s *service) Index(c *gin.Context) {
	ctx := c.Request.Context()
	pokemon, err := s.cache.Remember(ctx, "pokemon", 5*time.Minute,
		&pokemon.PokemonList{}, func() (any, error) {
			return s.pokemonClient.GetPokemon(ctx)
		},
	)
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
