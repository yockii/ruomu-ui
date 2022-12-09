# ruomu-ui
界面中心

提供基础界面管理、菜单管理

依赖于 ruomu-core

# 注入信息
## 模块名称
ui

## 模块注入点
| 代码             | 注入类型        | 说明    |
|----------------|-------------|-------|
| /menu/add      | HTTP-POST   | 新增菜单  |
| /menu/update   | HTTP-PUT    | 修改菜单  |
| /menu/delete   | HTTP-DELETE | 删除菜单  |
| /menu/list     | HTTP-GET    | 菜单列表  |
| /page/add      | HTTP-POST   | 新增页面  |
| /page/update   | HTTP-PUT    | 修改页面  |
| /page/delete   | HTTP-DELETE | 删除页面  |
| /page/list     | HTTP-GET    | 页面列表  |
| /page/instance | HTTP-GET    | 单页面详情 |



## 注入请求
POST /module/add
```json
{
    "name": "界面中心",
    "code": "ui",
    "cmd": "./plugins/ruomu_ui.exe",
    "status": 1,
    "needDb": true,
    "needUserTokenExpire": true,
    "dependencies": [],
    "injects": [
        {
            "name": "新增菜单",
            "type": 2,
            "injectCode": "/menu/add",
            "authorizationCode": "menu:add"
        },
        {
            "name": "修改菜单",
            "type": 3,
            "injectCode": "/menu/update",
            "authorizationCode": "menu:update"
        },
        {
            "name": "删除菜单",
            "type": 4,
            "injectCode": "/menu/delete",
            "authorizationCode": "menu:delete"
        },
        {
            "name": "菜单列表",
            "type": 1,
            "injectCode": "/menu/list",
            "authorizationCode": "menu:list"
        },
        {
            "name": "新增页面",
            "type": 2,
            "injectCode": "/page/add",
            "authorizationCode": "page:add"
        },
        {
            "name": "修改页面",
            "type": 3,
            "injectCode": "/page/update",
            "authorizationCode": "page:update"
        },
        {
            "name": "删除页面",
            "type": 4,
            "injectCode": "/page/delete",
            "authorizationCode": "page:delete"
        },
        {
            "name": "页面列表",
            "type": 1,
            "injectCode": "/page/list",
            "authorizationCode": "page:list"
        },
        {
          "name": "页面详情",
          "type": 1,
          "injectCode": "/page/instance",
          "authorizationCode": "page:instance"
        }
    ]
}
```



[//]: # (goreleaser release --skip-publish --rm-dist --snapshot)
