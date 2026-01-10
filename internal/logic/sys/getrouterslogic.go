// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sys

import (
	"context"
	"database/sql"
	"strconv"
	"strings"
	"unicode"

	"gozero-ruoyi-vue-plus/internal/model/sys"
	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"
	"gozero-ruoyi-vue-plus/internal/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRoutersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// MenuNode 菜单节点（带子节点）
type MenuNode struct {
	*sys.SysMenu
	Children []*MenuNode
}

// 获取路由信息
func NewGetRoutersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRoutersLogic {
	return &GetRoutersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetRoutersLogic) GetRouters() (resp *types.RouterResp, err error) {
	// 1. 从 JWT token 中获取用户ID
	userId, err := util.GetUserIdFromContext(l.ctx)
	if err != nil {
		return &types.RouterResp{
			BaseResp: types.BaseResp{
				Code: 401,
				Msg:  "未授权，请先登录",
			},
			Data: []types.RouterVo{},
		}, nil
	}

	// 2. 查询用户的菜单树
	menus, err := l.selectMenuTreeByUserId(userId)
	if err != nil {
		l.Errorf("查询菜单树失败: %v", err)
		return &types.RouterResp{
			BaseResp: types.BaseResp{
				Code: 500,
				Msg:  "查询菜单树失败",
			},
			Data: []types.RouterVo{},
		}, err
	}

	// 3. 构建路由树
	routers := l.buildMenus(menus)

	return &types.RouterResp{
		BaseResp: types.BaseResp{
			Code: 200,
			Msg:  "操作成功",
		},
		Data: routers,
	}, nil
}

// selectMenuTreeByUserId 根据用户ID查询菜单树
func (l *GetRoutersLogic) selectMenuTreeByUserId(userId int64) ([]*MenuNode, error) {
	var menus []*sys.SysMenu

	// 检查是否为超级管理员
	isSuperAdmin := l.isSuperAdmin(userId)
	if isSuperAdmin {
		// 超级管理员返回所有菜单
		menus = l.selectMenuTreeAll()
	} else {
		// 普通用户根据权限查询
		var err error
		menus, err = l.selectMenuListByUserId(userId)
		if err != nil {
			return nil, err
		}
	}

	// 构建树形结构
	return l.getChildPerms(menus, 0), nil
}

// isSuperAdmin 检查用户是否为超级管理员
func (l *GetRoutersLogic) isSuperAdmin(userId int64) bool {
	// 查询用户的角色
	query := `
		SELECT DISTINCT r.role_id, r.role_key
		FROM sys_role r
		INNER JOIN sys_user_role ur ON r.role_id = ur.role_id
		WHERE ur.user_id = ? 
		  AND r.status = '0'
		  AND r.del_flag = '0'
	`

	type roleRow struct {
		RoleId  int64  `db:"role_id"`
		RoleKey string `db:"role_key"`
	}

	var rows []roleRow
	err := l.svcCtx.DB.QueryRowsPartialCtx(l.ctx, &rows, query, userId)
	if err != nil {
		return false
	}

	for _, row := range rows {
		if row.RoleId == 1 || row.RoleKey == "superadmin" {
			return true
		}
	}
	return false
}

// selectMenuTreeAll 查询所有菜单（超级管理员）
func (l *GetRoutersLogic) selectMenuTreeAll() []*sys.SysMenu {
	query := `
		SELECT menu_id, menu_name, parent_id, order_num, path, component, query_param, 
		       is_frame, is_cache, menu_type, visible, status, perms, icon, 
		       create_dept, create_by, create_time, update_by, update_time, remark
		FROM sys_menu
		WHERE menu_type IN ('M', 'C')
		  AND status = '0'
		ORDER BY parent_id ASC, order_num ASC
	`

	var menus []*sys.SysMenu
	err := l.svcCtx.DB.QueryRowsPartialCtx(l.ctx, &menus, query)
	if err != nil {
		l.Errorf("查询所有菜单失败: %v", err)
		return []*sys.SysMenu{}
	}
	return menus
}

// selectMenuListByUserId 根据用户ID查询菜单列表（非超级管理员）
func (l *GetRoutersLogic) selectMenuListByUserId(userId int64) ([]*sys.SysMenu, error) {
	// 构建查询 SQL，通过用户角色查询菜单
	query := `
		SELECT DISTINCT m.menu_id, m.menu_name, m.parent_id, m.order_num, m.path, m.component, m.query_param,
		       m.is_frame, m.is_cache, m.menu_type, m.visible, m.status, m.perms, m.icon,
		       m.create_dept, m.create_by, m.create_time, m.update_by, m.update_time, m.remark
		FROM sys_menu m
		INNER JOIN sys_role_menu rm ON m.menu_id = rm.menu_id
		INNER JOIN sys_user_role ur ON rm.role_id = ur.role_id
		INNER JOIN sys_role r ON ur.role_id = r.role_id
		WHERE ur.user_id = ?
		  AND m.menu_type IN ('M', 'C')
		  AND m.status = '0'
		  AND r.status = '0'
		ORDER BY m.parent_id ASC, m.order_num ASC
	`

	var menus []*sys.SysMenu
	err := l.svcCtx.DB.QueryRowsPartialCtx(l.ctx, &menus, query, userId)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return menus, nil
}

// getChildPerms 根据父节点ID构建菜单树
func (l *GetRoutersLogic) getChildPerms(menus []*sys.SysMenu, parentId int64) []*MenuNode {
	var result []*MenuNode
	for _, menu := range menus {
		if menu.ParentId == parentId {
			node := &MenuNode{
				SysMenu:  menu,
				Children: l.getChildPerms(menus, menu.MenuId),
			}
			result = append(result, node)
		}
	}
	return result
}

// buildMenus 构建路由树
func (l *GetRoutersLogic) buildMenus(menus []*MenuNode) []types.RouterVo {
	var routers []types.RouterVo

	for _, menu := range menus {
		// 构建路由名称：routeName + menuId
		routeName := l.getRouteName(menu)
		name := routeName + int64ToString(menu.MenuId)

		router := types.RouterVo{
			Hidden:    menu.Visible == "1",
			Name:      name,
			Path:      l.getRouterPath(menu),
			Component: l.getComponentInfo(menu),
			Query:     l.getQueryParam(menu),
		}

		// 构建 Meta 信息
		meta := types.MetaVo{
			Title: menu.MenuName,
			Icon:  menu.Icon,
		}
		if menu.IsCache == 1 {
			meta.NoCache = true
		}
		if menu.Path != "" && l.isHttp(menu.Path) {
			meta.Link = menu.Path
		}
		if menu.Remark != "" {
			meta.Remark = menu.Remark
		}
		router.Meta = meta

		// 处理子菜单
		if len(menu.Children) > 0 && menu.MenuType == "M" {
			// 目录类型，有子菜单
			router.AlwaysShow = true
			router.Redirect = "noRedirect"
			router.Children = l.buildMenus(menu.Children)
		} else if l.isMenuFrame(menu) {
			// 一级菜单且是外链
			router.Meta = types.MetaVo{} // 清空 meta
			frameName := capitalize(menu.Path) + int64ToString(menu.MenuId)
			children := types.RouterVo{
				Path:      menu.Path,
				Component: l.getComponent(menu),
				Name:      frameName,
				Query:     l.getQueryParam(menu),
			}
			childrenMeta := types.MetaVo{
				Title: menu.MenuName,
				Icon:  menu.Icon,
			}
			if menu.IsCache == 1 {
				childrenMeta.NoCache = true
			}
			if menu.Path != "" && l.isHttp(menu.Path) {
				childrenMeta.Link = menu.Path
			}
			if menu.Remark != "" {
				childrenMeta.Remark = menu.Remark
			}
			children.Meta = childrenMeta
			router.Children = []types.RouterVo{children}
		} else if menu.ParentId == 0 && l.isInnerLink(menu) {
			// 一级菜单且是内链
			router.Meta = types.MetaVo{
				Title: menu.MenuName,
				Icon:  menu.Icon,
			}
			router.Path = "/"
			routerPath := l.innerLinkReplaceEach(menu.Path)
			innerLinkName := capitalize(routerPath) + int64ToString(menu.MenuId)
			children := types.RouterVo{
				Path:      routerPath,
				Component: "InnerLink",
				Name:      innerLinkName,
				Meta: types.MetaVo{
					Title: menu.MenuName,
					Icon:  menu.Icon,
					Link:  menu.Path,
				},
			}
			router.Children = []types.RouterVo{children}
		}

		routers = append(routers, router)
	}

	return routers
}

// getRouteName 获取路由名称
func (l *GetRoutersLogic) getRouteName(menu *MenuNode) string {
	if l.isMenuFrame(menu) {
		return ""
	}
	return capitalize(menu.Path)
}

// getRouterPath 获取路由路径
func (l *GetRoutersLogic) getRouterPath(menu *MenuNode) string {
	if l.isInnerLink(menu) {
		return l.innerLinkReplaceEach(menu.Path)
	}
	// 非外链并且是一级目录（类型为目录）
	if menu.ParentId == 0 && menu.MenuType == "M" && menu.IsFrame == 1 {
		return "/" + menu.Path
	}
	// 非外链并且是一级目录（类型为菜单）
	if l.isMenuFrame(menu) {
		return "/"
	}
	return menu.Path
}

// getComponentInfo 获取组件信息
func (l *GetRoutersLogic) getComponentInfo(menu *MenuNode) string {
	if menu.Component.Valid && menu.Component.String != "" && !l.isMenuFrame(menu) {
		return menu.Component.String
	}
	if l.isInnerLink(menu) {
		return "InnerLink"
	}
	if (!menu.Component.Valid || menu.Component.String == "") && l.isParentView(menu) {
		return "ParentView"
	}
	return ""
}

// getComponent 获取组件
func (l *GetRoutersLogic) getComponent(menu *MenuNode) string {
	if menu.Component.Valid {
		return menu.Component.String
	}
	return ""
}

// getQueryParam 获取查询参数
func (l *GetRoutersLogic) getQueryParam(menu *MenuNode) string {
	if menu.QueryParam.Valid {
		return menu.QueryParam.String
	}
	return ""
}

// isMenuFrame 是否是一级菜单外链
func (l *GetRoutersLogic) isMenuFrame(menu *MenuNode) bool {
	return menu.ParentId == 0 && menu.MenuType == "C" && menu.IsFrame == 1
}

// isInnerLink 是否为内链组件
func (l *GetRoutersLogic) isInnerLink(menu *MenuNode) bool {
	return menu.IsFrame == 0 && l.isHttp(menu.Path)
}

// isParentView 是否为父视图
func (l *GetRoutersLogic) isParentView(menu *MenuNode) bool {
	return menu.ParentId != 0 && menu.MenuType == "M"
}

// isHttp 判断是否为 http(s):// 开头
func (l *GetRoutersLogic) isHttp(path string) bool {
	return strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://")
}

// innerLinkReplaceEach 内链域名特殊字符替换
func (l *GetRoutersLogic) innerLinkReplaceEach(path string) string {
	result := strings.ReplaceAll(path, "http://", "")
	result = strings.ReplaceAll(result, "https://", "")
	result = strings.ReplaceAll(result, "www.", "")
	result = strings.ReplaceAll(result, ".", "/")
	result = strings.ReplaceAll(result, ":", "/")
	return result
}

// capitalize 首字母转大写
func capitalize(s string) string {
	if s == "" {
		return ""
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

// int64ToString 将 int64 转换为 string
func int64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}
