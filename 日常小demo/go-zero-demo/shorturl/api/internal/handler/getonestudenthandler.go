package handler

import (
	"net/http"

	"shorturl/api/internal/logic"
	"shorturl/api/internal/svc"
	"shorturl/api/internal/types"

	"github.com/tal-tech/go-zero/rest/httpx"
)

func GetOneStudentHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetonestudentReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewGetOneStudentLogic(r.Context(), ctx)
		resp, err := l.GetOneStudent(req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
