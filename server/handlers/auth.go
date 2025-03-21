package handlers

import (
	"fmt"
	"net/http"

	"github.com/khalidibnwalid/Luma/core"
)

func (ctx *HandlerContext) AuthLogin(w http.ResponseWriter, r *http.Request) {
	token, err := core.GenerateJwtToken(ctx.JwtSecret, "000000000000000000000000")
	if err != nil {
		fmt.Println(err)
	}

	w.Header().Add("Authorization", "Bearer "+token)
}

func (ctx *HandlerContext) AuthGET(w http.ResponseWriter, r *http.Request) {

}
