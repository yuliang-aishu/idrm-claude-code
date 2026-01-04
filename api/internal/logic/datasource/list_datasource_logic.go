package datasource

import (
	"context"

	datasourceModel "github.com/yuliang-aishu/idrm-claude-code/spec-cc-0104/model/datasource/datasource"
	"github.com/yuliang-aishu/idrm-claude-code/spec-cc-0104/api/internal/svc"
	"github.com/yuliang-aishu/idrm-claude-code/spec-cc-0104/api/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type ListDataSourceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListDataSourceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListDataSourceLogic {
	return &ListDataSourceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ListDataSource 实现列表查询逻辑
func (l *ListDataSourceLogic) ListDataSource(req *types.ListDataSourceReq) (resp *types.ListDataSourceResp, err error) {
	// 参数校验
	if req.Offset <= 0 {
		req.Offset = 1
	}
	if req.Limit <= 0 || req.Limit > 2000 {
		req.Limit = 10
	}

	// 转换查询参数
	query := &datasourceModel.DataSourceQuery{
		Offset:    req.Offset,
		Limit:     req.Limit,
		Keyword:   req.Keyword,
		Status:    req.Status,
		Sort:      req.Sort,
		Direction: req.Direction,
	}

	// 调用 Model 层查询
	dataSources, totalCount, err := l.svcCtx.DataSourceModel.FindList(l.ctx, query)
	if err != nil {
		l.Errorf("查询数据源列表失败: %v", err)
		return nil, err
	}

	// 转换为响应格式
	entries := make([]types.DataSourceResp, 0, len(dataSources))
	for _, ds := range dataSources {
		entries = append(entries, l.toResp(ds))
	}

	return &types.ListDataSourceResp{
		Entries:    entries,
		TotalCount: totalCount,
	}, nil
}

// toResp 转换为响应结构
func (l *ListDataSourceLogic) toResp(data *datasourceModel.DataSource) types.DataSourceResp {
	return types.DataSourceResp{
		Id:          data.Id,
		Name:        data.Name,
		Type:        data.Type,
		Host:        data.Host,
		Port:        data.Port,
		Database:    data.Database,
		Description: data.Description,
		Status:      data.Status,
		SortOrder:   data.SortOrder,
		CreatedAt:   data.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   data.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
