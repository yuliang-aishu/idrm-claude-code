// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package datasource

import (
	"net/http"

	"github.com/yuliang-aishu/idrm-claude-code/spec-cc-0104/api/internal/logic/datasource"
	"github.com/yuliang-aishu/idrm-claude-code/spec-cc-0104/api/internal/svc"
	"github.com/yuliang-aishu/idrm-claude-code/spec-cc-0104/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func CreateDataSourceHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateDataSourceReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := datasource.NewCreateDataSourceLogic(r.Context(), svcCtx)
		resp, err := l.CreateDataSource(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
