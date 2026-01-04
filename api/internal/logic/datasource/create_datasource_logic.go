package datasource

import (
	"context"

	"github.com/jinguoxing/idrm-go-base/errorx"
	"github.com/jinguoxing/idrm-go-base/validator"
	datasourceModel "github.com/yuliang-aishu/idrm-claude-code/spec-cc-0104/model/datasource/datasource"
	"github.com/yuliang-aishu/idrm-claude-code/spec-cc-0104/api/internal/svc"
	"github.com/yuliang-aishu/idrm-claude-code/spec-cc-0104/api/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateDataSourceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateDataSourceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateDataSourceLogic {
	return &CreateDataSourceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CreateDataSource 实现创建数据源逻辑
func (l *CreateDataSourceLogic) CreateDataSource(req *types.CreateDataSourceReq) (resp *types.CreateDataSourceResp, err error) {
	// 参数校验
	if err := validator.Validate(req); err != nil {
		l.Errorf("参数校验失败: %v", err)
		return nil, errorx.NewWithCode(datasourceModel.ErrCodeParamInvalid)
	}

	// 验证数据源类型
	if !l.isValidDataSourceType(req.Type) {
		l.Errorf("不支持的数据源类型: %s", req.Type)
		return nil, errorx.NewWithCode(datasourceModel.ErrCodeParamInvalid)
	}

	// 验证状态
	if req.Status == "" {
		req.Status = datasourceModel.DataSourceStatusEnabled
	} else if req.Status != datasourceModel.DataSourceStatusEnabled && req.Status != datasourceModel.DataSourceStatusDisabled {
		l.Errorf("无效的状态: %s", req.Status)
		return nil, errorx.NewWithCode(datasourceModel.ErrCodeParamInvalid)
	}

	// 检查名称是否重复
	exists, err := l.svcCtx.DataSourceModel.CheckNameExists(l.ctx, req.Name)
	if err != nil {
		l.Errorf("检查名称重复失败: %v", err)
		return nil, err
	}
	if exists {
		l.Errorf("数据源名称已存在: %s", req.Name)
		return nil, errorx.NewWithCode(datasourceModel.ErrCodeNameExists)
	}

	// 创建数据源实体
	dataSource := l.toDataSource(req)

	// 加密密码（TODO: 实现加密逻辑）
	// encryptedPassword, err := datasourceModel.EncryptPassword(req.Password)
	// if err != nil {
	// 	l.Errorf("密码加密失败: %v", err)
	// 	return nil, err
	// }
	// dataSource.Password = encryptedPassword

	// 生成 UUID v7 主键
	dataSource.Id = datasourceModel.GenerateUUIDv7()

	// 连接测试
	if err := l.svcCtx.DataSourceModel.TestConnection(l.ctx, dataSource); err != nil {
		l.Errorf("连接测试失败: %v", err)
		return nil, errorx.NewWithCode(datasourceModel.ErrCodeConnectionTest)
	}

	// 插入数据源
	createdDataSource, err := l.svcCtx.DataSourceModel.Insert(l.ctx, dataSource)
	if err != nil {
		l.Errorf("插入数据源失败: %v", err)
		return nil, err
	}

	l.Infof("创建数据源成功: %s", createdDataSource.Id)

	return &types.CreateDataSourceResp{
		Id:        createdDataSource.Id,
		CreatedAt: createdDataSource.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

// isValidDataSourceType 验证数据源类型
func (l *CreateDataSourceLogic) isValidDataSourceType(dataType string) bool {
	validTypes := map[string]bool{
		datasourceModel.DataSourceTypeMySQL:      true,
		datasourceModel.DataSourceTypePostgreSQL: true,
		datasourceModel.DataSourceTypeRedis:      true,
		datasourceModel.DataSourceTypeMongoDB:    true,
		datasourceModel.DataSourceTypeSQLServer:  true,
	}
	return validTypes[dataType]
}

// toDataSource 转换为 DataSource 实体
func (l *CreateDataSourceLogic) toDataSource(req *types.CreateDataSourceReq) *datasourceModel.DataSource {
	return &datasourceModel.DataSource{
		Name:        req.Name,
		Type:        req.Type,
		Host:        req.Host,
		Port:        req.Port,
		Database:    req.Database,
		Username:    req.Username,
		Password:    req.Password, // 将在上层加密
		Description: req.Description,
		Status:      req.Status,
		SortOrder:   req.SortOrder,
	}
}
