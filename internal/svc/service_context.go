package svc

import (
	"context"
	"fmt"
	"time"

	"gozero-ruoyi-vue-plus/internal/config"
	"gozero-ruoyi-vue-plus/internal/model/sys"
	"gozero-ruoyi-vue-plus/internal/oss"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config    config.Config
	RedisConn *redis.Redis
	DB        sqlx.SqlConn

	// 系统模型 - 按字母顺序
	SysClientModel        sys.SysClientModel
	SysConfigModel        sys.SysConfigModel
	SysDeptModel          sys.SysDeptModel
	SysDictDataModel      sys.SysDictDataModel
	SysDictTypeModel      sys.SysDictTypeModel
	SysLogininforModel    sys.SysLogininforModel
	SysMenuModel          sys.SysMenuModel
	SysNoticeModel        sys.SysNoticeModel
	SysOperLogModel       sys.SysOperLogModel
	SysOssModel           sys.SysOssModel
	SysOssConfigModel     sys.SysOssConfigModel
	SysPostModel          sys.SysPostModel
	SysRoleModel          sys.SysRoleModel
	SysRoleDeptModel      sys.SysRoleDeptModel
	SysRoleMenuModel      sys.SysRoleMenuModel
	SysSocialModel        sys.SysSocialModel
	SysTenantModel        sys.SysTenantModel
	SysTenantPackageModel sys.SysTenantPackageModel
	SysUserModel          sys.SysUserModel
	SysUserPostModel      sys.SysUserPostModel
	SysUserRoleModel      sys.SysUserRoleModel

	// OSS客户端管理器
	OssManager *oss.OssManager
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化MySQL连接并测试
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	if err := checkMySQLConnection(conn); err != nil {
		logx.Errorf("MySQL连接失败: %v", err)
		panic(fmt.Sprintf("MySQL连接失败，应用启动终止: %v", err))
	}
	logx.Infof("MySQL连接成功")

	// 初始化Redis连接并测试
	rds := redis.MustNewRedis(c.Redis)
	if err := checkRedisConnection(rds); err != nil {
		logx.Errorf("Redis连接失败: %v", err)
		panic(fmt.Sprintf("Redis连接失败，应用启动终止: %v", err))
	}
	logx.Infof("Redis连接成功")

	return &ServiceContext{
		Config:    c,
		RedisConn: rds,
		DB:        conn,

		// 初始化所有系统模型 - 不使用缓存
		SysClientModel:        sys.NewSysClientModel(conn),
		SysConfigModel:        sys.NewSysConfigModel(conn),
		SysDeptModel:          sys.NewSysDeptModel(conn),
		SysDictDataModel:      sys.NewSysDictDataModel(conn),
		SysDictTypeModel:      sys.NewSysDictTypeModel(conn),
		SysLogininforModel:    sys.NewSysLogininforModel(conn),
		SysMenuModel:          sys.NewSysMenuModel(conn),
		SysNoticeModel:        sys.NewSysNoticeModel(conn),
		SysOperLogModel:       sys.NewSysOperLogModel(conn),
		SysOssModel:           sys.NewSysOssModel(conn),
		SysOssConfigModel:     sys.NewSysOssConfigModel(conn),
		SysPostModel:          sys.NewSysPostModel(conn),
		SysRoleModel:          sys.NewSysRoleModel(conn),
		SysRoleDeptModel:      sys.NewSysRoleDeptModel(conn),
		SysRoleMenuModel:      sys.NewSysRoleMenuModel(conn),
		SysSocialModel:        sys.NewSysSocialModel(conn),
		SysTenantModel:        sys.NewSysTenantModel(conn),
		SysTenantPackageModel: sys.NewSysTenantPackageModel(conn),
		SysUserModel:          sys.NewSysUserModel(conn),
		SysUserPostModel:      sys.NewSysUserPostModel(conn),
		SysUserRoleModel:      sys.NewSysUserRoleModel(conn),

		// 初始化OSS客户端管理器
		OssManager: oss.NewOssManager(sys.NewSysOssConfigModel(conn)),
	}
}

// checkMySQLConnection 检查MySQL连接
func checkMySQLConnection(conn sqlx.SqlConn) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 执行一个简单的查询来验证连接
	_, err := conn.ExecCtx(ctx, "SELECT 1")
	if err != nil {
		return fmt.Errorf("无法连接到MySQL数据库: %w", err)
	}
	return nil
}

// checkRedisConnection 检查Redis连接
func checkRedisConnection(rds *redis.Redis) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 执行PING命令来验证连接
	result := rds.PingCtx(ctx)
	if !result {
		return fmt.Errorf("无法连接到Redis: PING命令返回失败")
	}
	return nil
}
