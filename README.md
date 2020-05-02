# Sinan server golang

基于 Golang 的网络书签服务。采用 RESTful 风格 API 接口。

## 接口定义

GET /bookmark

    获取所有书签

GET /bookmark/{id}

    获取指定 id 的书签

POST /bookmark

    提交一个新书签

PUT /bookmark/{id}

    提交一个指定 id 的书签

DELETE /bookmark/{id}

    删除指定 id 的书签

## Database

目前使用 json 文件。

## Todo

* 使用 memcache 优化缓存性能
* 使用 swagger 作为 web 框架
