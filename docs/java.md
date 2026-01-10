curl 'http://192.168.3.216:58888/system/user/getInfo' \
-H 'Accept: application/json, text/plain, */*' \
-H 'Accept-Language: zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7,ja;q=0.6,zh-TW;q=0.5,fr;q=0.4' \
-H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOjEsInVzZXJuYW1lIjoiYWRtaW4iLCJ0ZW5hbnRJZCI6IjAwMDAwMCIsImV4cCI6MTc2ODYyMjY1NSwibmJmIjoxNzY4MDE3ODU1LCJpYXQiOjE3NjgwMTc4NTV9.2yC8plxIZwLrxIe-P6j16gOu5wSNssIpi1W3yrm3qbc' \
-H 'Connection: keep-alive' \
-H 'Content-Language: zh_CN' \
-H 'Origin: http://localhost' \
-H 'Referer: http://localhost/' \
-H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36' \
-H 'clientid: e5cd7e4891bf95d1d19206ce24a7b32e' \
--insecure

```json
{
  "code": 200,
  "msg": "操作成功",
  "data": {
    "user": {
      "userId": 1,
      "tenantId": "000000",
      "deptId": 103,
      "userName": "admin",
      "nickName": "疯狂的狮子Li",
      "userType": "sys_user",
      "email": "crazyLionLi@163.com",
      "phonenumber": "15888888888",
      "sex": "1",
      "avatar": null,
      "status": "0",
      "loginIp": "192.168.3.216",
      "loginDate": "2026-01-10 13:41:43",
      "remark": "管理员",
      "createTime": "2026-01-02 21:32:02",
      "deptName": "研发部门",
      "roles": [
        {
          "roleId": 1,
          "roleName": "超级管理员",
          "roleKey": "superadmin",
          "roleSort": 1,
          "dataScope": "1",
          "menuCheckStrictly": null,
          "deptCheckStrictly": null,
          "status": "0",
          "remark": null,
          "createTime": null,
          "flag": false,
          "superAdmin": true
        }
      ],
      "roleIds": null,
      "postIds": null,
      "roleId": null
    },
    "permissions": [
      "*:*:*"
    ],
    "roles": [
      "superadmin"
    ]
  }
}
```

curl 'http://192.168.3.216:58888/system/menu/getRouters' \
-H 'Accept: application/json, text/plain, */*' \
-H 'Accept-Language: zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7,ja;q=0.6,zh-TW;q=0.5,fr;q=0.4' \
-H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOjEsInVzZXJuYW1lIjoiYWRtaW4iLCJ0ZW5hbnRJZCI6IjAwMDAwMCIsImV4cCI6MTc2ODYyMjY1NSwibmJmIjoxNzY4MDE3ODU1LCJpYXQiOjE3NjgwMTc4NTV9.2yC8plxIZwLrxIe-P6j16gOu5wSNssIpi1W3yrm3qbc' \
-H 'Connection: keep-alive' \
-H 'Content-Language: zh_CN' \
-H 'Origin: http://localhost' \
-H 'Referer: http://localhost/' \
-H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36' \
-H 'clientid: e5cd7e4891bf95d1d19206ce24a7b32e' \
--insecure

```json
{
  "code": 200,
  "msg": "操作成功",
  "data": [
    {
      "name": "System1",
      "path": "/system",
      "hidden": false,
      "redirect": "noRedirect",
      "component": "Layout",
      "alwaysShow": true,
      "meta": {
        "title": "系统管理",
        "icon": "system",
        "noCache": false,
        "link": null,
        "activeMenu": null
      },
      "children": [
        {
          "name": "User100",
          "path": "user",
          "hidden": false,
          "component": "system/user/index",
          "meta": {
            "title": "用户管理",
            "icon": "user",
            "noCache": false,
            "link": null,
            "activeMenu": null
          }
        },
        {
          "name": "User-auth/role/:userId131",
          "path": "user-auth/role/:userId",
          "hidden": true,
          "component": "system/user/authRole",
          "meta": {
            "title": "分配角色",
            "icon": "#",
            "noCache": true,
            "link": null,
            "activeMenu": "/system/user"
          }
        },
        {
          "name": "Role101",
          "path": "role",
          "hidden": false,
          "component": "system/role/index",
          "meta": {
            "title": "角色管理",
            "icon": "peoples",
            "noCache": false,
            "link": null,
            "activeMenu": null
          }
        },
        {
          "name": "Role-auth/user/:roleId130",
          "path": "role-auth/user/:roleId",
          "hidden": true,
          "component": "system/role/authUser",
          "meta": {
            "title": "分配用户",
            "icon": "#",
            "noCache": true,
            "link": null,
            "activeMenu": "/system/role"
          }
        },
        {
          "name": "Menu102",
          "path": "menu",
          "hidden": false,
          "component": "system/menu/index",
          "meta": {
            "title": "菜单管理",
            "icon": "tree-table",
            "noCache": false,
            "link": null,
            "activeMenu": null
          }
        },
        {
          "name": "Dept103",
          "path": "dept",
          "hidden": false,
          "component": "system/dept/index",
          "meta": {
            "title": "部门管理",
            "icon": "tree",
            "noCache": false,
            "link": null,
            "activeMenu": null
          }
        },
        {
          "name": "Post104",
          "path": "post",
          "hidden": false,
          "component": "system/post/index",
          "meta": {
            "title": "岗位管理",
            "icon": "post",
            "noCache": false,
            "link": null,
            "activeMenu": null
          }
        },
        {
          "name": "Dict105",
          "path": "dict",
          "hidden": false,
          "component": "system/dict/index",
          "meta": {
            "title": "字典管理",
            "icon": "dict",
            "noCache": false,
            "link": null,
            "activeMenu": null
          }
        },
        {
          "name": "Dict-data/index/:dictId132",
          "path": "dict-data/index/:dictId",
          "hidden": true,
          "component": "system/dict/data",
          "meta": {
            "title": "字典数据",
            "icon": "#",
            "noCache": true,
            "link": null,
            "activeMenu": "/system/dict"
          }
        },
        {
          "name": "Config106",
          "path": "config",
          "hidden": false,
          "component": "system/config/index",
          "meta": {
            "title": "参数设置",
            "icon": "edit",
            "noCache": false,
            "link": null,
            "activeMenu": null
          }
        },
        {
          "name": "Notice107",
          "path": "notice",
          "hidden": false,
          "component": "system/notice/index",
          "meta": {
            "title": "通知公告",
            "icon": "message",
            "noCache": false,
            "link": null,
            "activeMenu": null
          }
        },
        {
          "name": "Log108",
          "path": "log",
          "hidden": false,
          "redirect": "noRedirect",
          "component": "ParentView",
          "alwaysShow": true,
          "meta": {
            "title": "日志管理",
            "icon": "log",
            "noCache": false,
            "link": null,
            "activeMenu": null
          },
          "children": [
            {
              "name": "Operlog500",
              "path": "operlog",
              "hidden": false,
              "component": "monitor/operlog/index",
              "meta": {
                "title": "操作日志",
                "icon": "form",
                "noCache": false,
                "link": null,
                "activeMenu": null
              }
            },
            {
              "name": "Logininfor501",
              "path": "logininfor",
              "hidden": false,
              "component": "monitor/logininfor/index",
              "meta": {
                "title": "登录日志",
                "icon": "logininfor",
                "noCache": false,
                "link": null,
                "activeMenu": null
              }
            }
          ]
        },
        {
          "name": "Oss118",
          "path": "oss",
          "hidden": false,
          "component": "system/oss/index",
          "meta": {
            "title": "文件管理",
            "icon": "upload",
            "noCache": false,
            "link": null,
            "activeMenu": null
          }
        },
        {
          "name": "Oss-config/index133",
          "path": "oss-config/index",
          "hidden": true,
          "component": "system/oss/config",
          "meta": {
            "title": "文件配置管理",
            "icon": "#",
            "noCache": true,
            "link": null,
            "activeMenu": "/system/oss"
          }
        },
        {
          "name": "Client123",
          "path": "client",
          "hidden": false,
          "component": "system/client/index",
          "meta": {
            "title": "客户端管理",
            "icon": "international",
            "noCache": false,
            "link": null,
            "activeMenu": null
          }
        }
      ]
    },
    {
      "name": "Monitor2",
      "path": "/monitor",
      "hidden": false,
      "redirect": "noRedirect",
      "component": "Layout",
      "alwaysShow": true,
      "meta": {
        "title": "系统监控",
        "icon": "monitor",
        "noCache": false,
        "link": null,
        "activeMenu": null
      },
      "children": [
        {
          "name": "Online109",
          "path": "online",
          "hidden": false,
          "component": "monitor/online/index",
          "meta": {
            "title": "在线用户",
            "icon": "online",
            "noCache": false,
            "link": null,
            "activeMenu": null
          }
        },
        {
          "name": "Cache113",
          "path": "cache",
          "hidden": false,
          "component": "monitor/cache/index",
          "meta": {
            "title": "缓存监控",
            "icon": "redis",
            "noCache": false,
            "link": null,
            "activeMenu": null
          }
        }
      ]
    }
  ]
}
```