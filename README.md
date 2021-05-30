# Sinan web server

基于 Golang 的网络书签服务。采用 RESTful 风格 API 接口。

## 接口定义

GET /index

    获取所有书签

GET /index/{id}

    获取指定 id 的书签

POST /index

    提交一个新书签

PUT /index/{id}

    提交一个指定 id 的书签

DELETE /index/{id}

    删除指定 id 的书签

## Database

使用 go-sqlite3 模块操作 sqlite 作为数据库。
