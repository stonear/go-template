package controller

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stonear/go-template/app/model"
)

type MhsController struct {
	// to do add properties
	Db *sql.DB
}

func (m *MhsController) Index(c *gin.Context) {
	rows, err := m.Db.Query("SELECT id, name, nrp FROM mhs")
	if err != nil {
		log.Fatalln(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "An error occured\n",
		})
	}
	defer rows.Close()

	result := model.Mhs{}
	results := []model.Mhs{}
	for rows.Next() {
		rows.Scan(&result.Id, &result.Name, &result.Nrp)
		results = append(results, result)
	}

	c.JSON(http.StatusOK, gin.H{
		"mhs": results,
	})
}
