package controller

import (
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Hello struct {
	// to do add properties
}

func (h *Hello) Default(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	io.WriteString(w, "Hello world!\n")
}

func (h *Hello) Name(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	io.WriteString(w, "Hello "+p.ByName("name")+"!\n")
}
