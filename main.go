// Sinan - 书签导航服务 API

package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/emicklei/go-restful"
	_ "github.com/mattn/go-sqlite3"
)

type IndexModel struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Url         string `json:"url"`
	Tags        string `json:"tags"`
	Description string `json:"description"`
}

type IndexModelList struct {
	Models []IndexModel `json:"index"`
}

type BaseDao struct {
	db *sql.DB
}

func (dao *BaseDao) Init(workspace string) error {
	_dir := "./" + workspace
	_, err := os.Stat(_dir)
	if os.IsNotExist(err) {
		err = os.Mkdir(_dir, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	dao.db, err = sql.Open("sqlite3", _dir+"/sinan.db")
	return err
}

type IndexDao struct {
	dao *BaseDao
}

func (id *IndexDao) Init(dao *BaseDao) error {
	id.dao = dao

	sqlCreateTable := `
    CREATE TABLE IF NOT EXISTS sinan_index(
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        title TEXT NOT NULL,
        url TEXT NOT NULL,
        tags TEXT,
        description TEXT
    );
    `
	_, err := dao.db.Exec(sqlCreateTable)
	return err
}

func (id *IndexDao) QueryIndex() (IndexModelList, error) {
	var list IndexModelList
	sqlQuery := `
    SELECT * FROM sinan_index
    `
	rows, err := id.dao.db.Query(sqlQuery)

	if err != nil {
		return list, err
	}

	for rows.Next() {
		var model IndexModel
		err = rows.Scan(&model.Id, &model.Title, &model.Url, &model.Tags,
			&model.Description)
		if err != nil {
			return list, err
		}
		list.Models = append(list.Models, model)
	}
	rows.Close()
	return list, nil
}

func (id *IndexDao) InsertIndex(list *IndexModelList) error {
	sqlInsert := `
    INSERT INTO sinan_index(title, url, tags, description) VALUES(?, ?, ?, ?)
    `
	stmt, err := id.dao.db.Prepare(sqlInsert)
	if err != nil {
		return err
	}

	for _, model := range list.Models {
		_, err := stmt.Exec(model.Title, model.Url,
			model.Tags, model.Description)
		if err != nil {
			return err
		}
	}

	return nil
}

func (id *IndexDao) DeleteIndexById(id_num int) error {
	sqlDelete := `
		DELETE FROM sinan_index where id=?
		`
	_, err := id.dao.db.Exec(sqlDelete, id_num)
	return err
}

func (id *IndexDao) UpdateIndexById(model *IndexModel, id_num int) error {
	sqlUpdate := `
		UPDATE sinan_index set title=? url=? tags=? description=? WHERE id=?
		`
	_, err := id.dao.db.Exec(sqlUpdate, model.Title, model.Url,
		model.Tags, model.Description, model.Id)
	return err
}

type BizDao struct {
	baseDao  BaseDao
	indexDao IndexDao
}

func (bd *BizDao) Init(workspace string) error {
	err := bd.baseDao.Init(workspace)
	if err != nil {
		return err
	}
	err = bd.indexDao.Init(&bd.baseDao)
	if err != nil {
		return err
	}

	return nil
}

type SinanService struct {
	daoMap map[string]BizDao
}

func (ss *SinanService) Init() error {
	ss.daoMap = make(map[string]BizDao)
	return nil
}

func (ss *SinanService) getDao(workspace string) (BizDao, error) {
	dao, ok := ss.daoMap[workspace]
	if !ok {
		dao = BizDao{}
		err := dao.Init(workspace)
		if err != nil {
			return dao, err
		}
		ss.daoMap[workspace] = dao
	}
	return dao, nil
}

func (ss *SinanService) GetIndex(request *restful.Request, response *restful.Response) {
	ws := request.PathParameter("workspace")
	dao, err := ss.getDao(ws)
	if err != nil {
		log.Fatal(err)
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, "init database failed.")
		return
	}

	list, err := dao.indexDao.QueryIndex()

	if err != nil {
		log.Fatal(err)
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, "query database failed.")
		return
	}

	response.WriteAsJson(list)
}

func (ss *SinanService) PostIndex(request *restful.Request, response *restful.Response) {
	ws := request.PathParameter("workspace")
	dao, err := ss.getDao(ws)
	if err != nil {
		log.Fatal(err)
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, "init database failed.")
		return
	}

	var list IndexModelList
	err = request.ReadEntity(&list)
	if err != nil {
		log.Fatal(err)
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusBadRequest, "bad request model.")
		return
	}

	err = dao.indexDao.InsertIndex(&list)
	if err != nil {
		log.Fatal(err)
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, "insert database failed.")
		return
	}
}

func (ss *SinanService) DeleteIndexById(request *restful.Request, response *restful.Response) {
	ws := request.PathParameter("workspace")
	dao, err := ss.getDao(ws)
	if err != nil {
		log.Fatal(err)
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, "init database failed.")
		return
	}

	id := request.PathParameter("id")
	id_num, err := strconv.Atoi(id)
	if err != nil {
		log.Fatal(err)
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusBadRequest, "bad id.")
		return
	}

	err = dao.indexDao.DeleteIndexById(id_num)
	if err != nil {
		log.Fatal(err)
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, "delete database failed.")
		return
	}
}

func (ss *SinanService) PostIndexById(request *restful.Request, response *restful.Response) {
	ws := request.PathParameter("workspace")
	dao, err := ss.getDao(ws)
	if err != nil {
		log.Fatal(err)
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, "init database failed.")
		return
	}

	id := request.PathParameter("id")
	id_num, err := strconv.Atoi(id)
	if err != nil {
		log.Fatal(err)
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusBadRequest, "bad id.")
		return
	}

	var model IndexModel
	err = request.ReadEntity(&model)
	if err != nil {
		log.Fatal(err)
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusBadRequest, "bad request model.")
		return
	}

	err = dao.indexDao.UpdateIndexById(&model, id_num)
	if err != nil {
		log.Fatal(err)
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, "update database failed.")
		return
	}
}

func main() {
	webService := &restful.WebService{}
	webService.Path("/").Produces(restful.MIME_JSON)

	sinanService := SinanService{}
	err := sinanService.Init()
	if err != nil {
		panic(err)
	}
	webService.Route(webService.GET("/index/{workspace}").
		To(sinanService.GetIndex).
		Param(webService.PathParameter("workspace", "the place of database").
			DataType("string")))

	webService.Route(webService.POST("/index/{workspace}").
		To(sinanService.PostIndex).
		Param(webService.PathParameter("workspace", "the place of database").
			DataType("string")))

	webService.Route(webService.DELETE("/index/{workspace}/{id}").
		To(sinanService.DeleteIndexById).
		Param(webService.PathParameter("workspace", "the place of database").
			DataType("string")).
		Param(webService.PathParameter("id", "identifier of the index").
			DataType("integer")))

	webService.Route(webService.POST("/index/{workspace}/{id}").
		To(sinanService.PostIndexById).
		Param(webService.PathParameter("workspace", "the place of database").
			DataType("string")).
		Param(webService.PathParameter("id", "identifier of the index").
			DataType("integer")))

	container := restful.NewContainer()
	container.Add(webService)

	log.Printf("start listening on localhost:8080")
	server := &http.Server{Addr: ":8080", Handler: container}
	log.Fatal(server.ListenAndServe())
}
