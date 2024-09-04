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
| /project/add   | HTTP-POST   | 新增项目  |
| /project/update   | HTTP-PUT    | 修改项目  |
| /project/delete   | HTTP-DELETE | 删除项目  |
| /project/list     | HTTP-GET    | 项目列表  |
| /page/add      | HTTP-POST   | 新增页面  |
| /page/update   | HTTP-PUT    | 修改页面  |
| /page/delete   | HTTP-DELETE | 删除页面  |
| /page/list     | HTTP-GET    | 页面列表  |
| /page/instance | HTTP-GET    | 页面详情  |
| /materialLib/add | HTTP-POST   | 新增物料库 |
| /materialLib/update | HTTP-PUT    | 修改物料库 |
| /materialLib/delete | HTTP-DELETE | 删除物料库 |
| /materialLib/list | HTTP-GET    | 物料库列表 |
| /materialLib/instance | HTTP-GET    | 物料库详情 |
| /materialComponentGroup/add | HTTP-POST   | 新增物料组件分组 |
| /materialComponentGroup/update | HTTP-PUT    | 修改物料组件分组 |
| /materialComponentGroup/delete | HTTP-DELETE | 删除物料组件分组 |
| /materialComponentGroup/list | HTTP-GET    | 物料组件分组列表 |
| /materialComponentGroup/instance | HTTP-GET    | 物料组件分组详情 |
| /materialComponent/add | HTTP-POST   | 新增物料组件 |
| /materialComponent/update | HTTP-PUT    | 修改物料组件 |
| /materialComponent/delete | HTTP-DELETE | 删除物料组件 |
| /materialComponent/list | HTTP-GET    | 物料组件列表 |
| /materialComponent/instance | HTTP-GET    | 物料组件详情 |



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
            "name": "新增项目",
            "type": 2,
            "injectCode": "/project/add",
            "authorizationCode": "project:add"
        },
        {
            "name": "修改项目",
            "type": 3,
            "injectCode": "/project/update",
            "authorizationCode": "project:update"
        },
        {
            "name": "删除项目",
            "type": 4,
            "injectCode": "/project/delete",
            "authorizationCode": "project:delete"
        },
        {
            "name": "项目列表",
            "type": 1,
            "injectCode": "/project/list",
            "authorizationCode": "project:list"
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
        },
        {
          "name": "新增物料库",
          "type": 2,
          "injectCode": "/materialLib/add",
          "authorizationCode": "materialLib:add"
        },
        {
          "name": "修改物料库",
          "type": 3,
          "injectCode": "/materialLib/update",
          "authorizationCode": "materialLib:update"
        },
      {
        "name": "删除物料库",
        "type": 4,
        "injectCode": "/materialLib/delete",
        "authorizationCode": "materialLib:delete"
      },
      {
        "name": "物料库列表",
        "type": 1,
        "injectCode": "/materialLib/list",
        "authorizationCode": "materialLib:list"
      },
      {
        "name": "物料库详情",
        "type": 1,
        "injectCode": "/materialLib/instance",
        "authorizationCode": "materialLib:instance"
      },
      {
        "name": "新增物料组件分组",
        "type": 2,
        "injectCode": "/materialComponentGroup/add",
        "authorizationCode": "materialComponentGroup:add"
      },
      {
        "name": "修改物料组件分组",
        "type": 3,
        "injectCode": "/materialComponentGroup/update",
        "authorizationCode": "materialComponentGroup:update"
      },
      {
        "name": "删除物料组件分组",
        "type": 4,
        "injectCode": "/materialComponentGroup/delete",
        "authorizationCode": "materialComponentGroup:delete"
      },
      {
        "name": "物料组件分组列表",
        "type": 1,
        "injectCode": "/materialComponentGroup/list",
        "authorizationCode": "materialComponentGroup:list"
      },
      {
        "name": "物料组件分组详情",
        "type": 1,
        "injectCode": "/materialComponentGroup/instance",
        "authorizationCode": "materialComponentGroup:instance"
      },
      {
        "name": "新增物料组件",
        "type": 2,
        "injectCode": "/materialComponent/add",
        "authorizationCode": "materialComponent:add"
      },
      {
        "name": "修改物料组件",
        "type": 3,
        "injectCode": "/materialComponent/update",
        "authorizationCode": "materialComponent:update"
      },
      {
        "name": "删除物料组件",
        "type": 4,
        "injectCode": "/materialComponent/delete",
        "authorizationCode": "materialComponent:delete"
      },
      {
        "name": "物料组件列表",
        "type": 1,
        "injectCode": "/materialComponent/list",
        "authorizationCode": "materialComponent:list"
      },
      {
        "name": "物料组件详情",
        "type": 1,
        "injectCode": "/materialComponent/instance",
        "authorizationCode": "materialComponent:instance"
      }
    ]
}
```



[//]: # (goreleaser release --skip-publish --clean --snapshot)
