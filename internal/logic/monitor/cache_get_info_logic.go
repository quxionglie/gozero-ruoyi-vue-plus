// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package monitor

import (
	"context"
	"strings"

	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CacheGetInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取缓存监控列表
func NewCacheGetInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CacheGetInfoLogic {
	return &CacheGetInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CacheGetInfoLogic) CacheGetInfo() (resp *types.CacheListResp, err error) {
	// 1. 使用 EvalCtx 执行 Lua 脚本获取 Redis INFO 信息
	// eval "return redis.call('info')" 0
	infoScript := "return redis.call('info')"
	infoResult, err := l.svcCtx.RedisConn.EvalCtx(l.ctx, infoScript, []string{}, []string{})
	if err != nil {
		l.Errorf("执行Redis INFO命令失败: %v", err)
		return &types.CacheListResp{
			BaseResp: types.BaseResp{
				Code: 500,
				Msg:  "获取Redis信息失败",
			},
		}, err
	}

	// 将结果转换为字符串
	infoStr, ok := infoResult.(string)
	if !ok {
		l.Errorf("Redis INFO返回结果格式错误")
		return &types.CacheListResp{
			BaseResp: types.BaseResp{
				Code: 500,
				Msg:  "Redis INFO返回结果格式错误",
			},
		}, nil
	}

	// 解析 INFO 结果
	infoMap := l.parseInfoString(infoStr)

	// 2. 使用 EvalCtx 执行 Lua 脚本获取数据库大小
	// eval "return redis.call('dbsize')" 0
	dbsizeScript := "return redis.call('dbsize')"
	dbSizeResult, err := l.svcCtx.RedisConn.EvalCtx(l.ctx, dbsizeScript, []string{}, []string{})
	if err != nil {
		l.Errorf("执行Redis DBSIZE命令失败: %v", err)
		// DBSIZE 失败不影响整体返回，设为0
		dbSizeResult = int64(0)
	}

	// 将结果转换为 int64
	var dbSize int64
	switch v := dbSizeResult.(type) {
	case int64:
		dbSize = v
	case int:
		dbSize = int64(v)
	default:
		l.Errorf("Redis DBSIZE返回结果格式错误: %T", v)
		dbSize = 0
	}

	// 3. 使用 EvalCtx 执行 Lua 脚本获取命令统计
	// eval "return redis.call('info','commandstats')" 0
	commandStatsScript := "return redis.call('info','commandstats')"
	commandStatsResult, err := l.svcCtx.RedisConn.EvalCtx(l.ctx, commandStatsScript, []string{}, []string{})
	if err != nil {
		l.Errorf("执行Redis INFO commandstats命令失败: %v", err)
		commandStatsResult = ""
	}

	// 将结果转换为字符串
	commandStatsStr := ""
	if resultStr, ok := commandStatsResult.(string); ok {
		commandStatsStr = resultStr
	}

	// 解析命令统计
	commandStats := l.parseCommandStats(commandStatsStr)

	return &types.CacheListResp{
		BaseResp: types.BaseResp{
			Code: 200,
			Msg:  "获取成功",
		},
		Data: types.CacheListInfoVo{
			Info:         infoMap,
			DbSize:       dbSize,
			CommandStats: commandStats,
		},
	}, nil
}

// parseInfoString 解析 Redis INFO 命令返回的字符串
func (l *CacheGetInfoLogic) parseInfoString(infoStr string) map[string]string {
	infoMap := make(map[string]string)
	// 统一处理不同格式的换行符（\r\n 或 \n）
	infoStr = strings.ReplaceAll(infoStr, "\r\n", "\n")
	lines := strings.Split(infoStr, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			infoMap[key] = value
		}
	}
	return infoMap
}

// parseCommandStats 解析 Redis INFO commandstats 命令返回的字符串
func (l *CacheGetInfoLogic) parseCommandStats(commandStatsStr string) []map[string]string {
	var commandStats []map[string]string
	if commandStatsStr == "" {
		return commandStats
	}

	// 统一处理不同格式的换行符（\r\n 或 \n）
	commandStatsStr = strings.ReplaceAll(commandStatsStr, "\r\n", "\n")
	lines := strings.Split(commandStatsStr, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// 格式: cmdstat_get:calls=123,usec=456,usec_per_call=3.70
		// 或者: cmdstat_get:calls=123
		if strings.HasPrefix(line, "cmdstat_") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				// 提取命令名（去掉 cmdstat_ 前缀）
				cmdName := strings.TrimPrefix(parts[0], "cmdstat_")
				// 解析统计信息
				stats := parts[1]
				// 提取 calls 值
				var calls string
				callsStart := strings.Index(stats, "calls=")
				if callsStart >= 0 {
					callsEnd := strings.Index(stats[callsStart:], ",")
					if callsEnd > 0 {
						calls = stats[callsStart+6 : callsStart+callsEnd]
					} else {
						calls = stats[callsStart+6:]
					}
				} else {
					calls = "0"
				}

				statMap := map[string]string{
					"name":  cmdName,
					"value": calls,
				}
				commandStats = append(commandStats, statMap)
			}
		}
	}
	return commandStats
}
