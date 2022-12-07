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
| /menu/update   | HTTP-POST   | 修改菜单  |
| /menu/delete   | HTTP-DELETE | 删除菜单  |
| /menu/list     | HTTP-GET    | 菜单列表  |
| /page/add      | HTTP-POST   | 新增页面  |
| /page/update   | HTTP-POST   | 修改页面  |
| /page/delete   | HTTP-DELETE | 删除页面  |
| /page/list     | HTTP-GET    | 页面列表  |
| /page/instance | HTTP-GET    | 单页面详情 |



## json规则



[//]: # (goreleaser release --skip-publish --rm-dist --snapshot)
