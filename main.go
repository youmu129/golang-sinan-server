// Sinan - 书签导航服务 API

package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/emicklei/go-restful"
	_ "github.com/mattn/go-sqlite3"
)

type IndexModel struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Url         string `json:"url"`
	Icon        string `json:"icon"`
	Description string `json:"description"`
}

type IndexList struct {
	Models []IndexModel `json:"index"`
}

type IndexService struct {
	db *sql.DB
}

func (s *IndexService) Init() {
	db, err := sql.Open("sqlite3", "./sinan.db")
	if err != nil {
		panic(err)
	}
	s.db = db

	sqlCreateTable := `
    CREATE TABLE IF NOT EXISTS sinan_index(
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        title TEXT NOT NULL,
        url TEXT NOT NULL,
        icon TEXT,
        description TEXT
    );
    `
	_, err = s.db.Exec(sqlCreateTable)
	if err != nil {
		panic(err)
	}
}

func (s *IndexService) GetIndex(request *restful.Request, response *restful.Response) {
	sqlQuery := `
    SELECT * FROM sinan_index
    `
	rows, err := s.db.Query(sqlQuery)

	if err != nil {
		log.Fatal(err)
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, "query database failed.")
		return
	}

	var list IndexList
	for rows.Next() {
		var model IndexModel
		err = rows.Scan(&model.Id, &model.Title, &model.Url, &model.Icon,
			&model.Description)
		if err != nil {
			log.Fatal(err)
			response.WriteError(500, err)
			return
		}
		list.Models = append(list.Models, model)
	}
	rows.Close()
	response.WriteAsJson(list)
}

func (s *IndexService) PostIndex(request *restful.Request, response *restful.Response) {
	var list IndexList
	err := request.ReadEntity(&list)
	if err != nil {
		log.Fatal(err)
		response.WriteError(500, err)
		return
	}

	sqlInsert := `
    INSERT INTO sinan_index(title, url, icon, description) VALIES(?, ?, ?)
    `
	stmt, err := s.db.Prepare(sqlInsert)
	if err != nil {
		log.Fatal(err)
		response.WriteError(500, err)
		return
	}

	for _, model := range list.Models {
		_, err := stmt.Exec(model.Title, model.Url,
			model.Icon, model.Description)
		if err != nil {
			log.Fatal(err)
			response.WriteError(500, err)
			return
		}
	}
}

func (s *IndexService) DeleteIndexById(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("id")
	_, err := strconv.Atoi(id)
	if err != nil {
		log.Fatal(err)
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusBadRequest, "id not found.")
	}
}

func (s *IndexService) PostIndexById(request *restful.Request, response *restful.Response) {

}

func main() {
	webService := &restful.WebService{}
	webService.Path("/").Produces(restful.MIME_JSON)

	indexService := IndexService{}
	indexService.Init()
	webService.Route(webService.GET("/index").
		To(indexService.GetIndex))
	webService.Route(webService.POST("/index").
		To(indexService.PostIndex))
	webService.Route(webService.DELETE("/index/{id}").
		To(indexService.DeleteIndexById).
		Param(webService.PathParameter("id", "identifier of the markbook").
			DataType("string")))
	webService.Route(webService.POST("/index/{id}").
		To(indexService.PostIndexById).
		Param(webService.PathParameter("id", "identifier of the markbook").
			DataType("string")))

	container := restful.NewContainer()
	container.Add(webService)

	log.Printf("start listening on localhost:8080")
	server := &http.Server{Addr: ":8080", Handler: container}
	log.Fatal(server.ListenAndServe())
}
