# Sinan web server

基于 Golang 的网络书签服务。采用 RESTful 风格 API 接口。

## 接口定义

- GET /index/{workspace}

    获取所有书签

- GET /index/{workspace}/{id}

    获取指定 id 的书签

- POST /index

    提交一个新书签

- POST /index/{workspace}/{id}

    更新一个指定 id 的书签

- DELETE /index/{workspace}/{id}

    删除指定 id 的书签

## Database

使用 go-sqlite3 模块操作 sqlite 作为数据库。

## TODO

考虑使用 Beego 框架，使用其提供的 ORM 和 Swagger 集成。
