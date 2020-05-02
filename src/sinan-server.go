// Sinan - 书签导航服务 API

package main

import (
    "github.com/emicklei/go-restful"

	"log"
    "net/http"
    "io/ioutil"
    "encoding/json"
    "strconv"
)

type BookmarkModel struct {
    Id int `json:"id"`
    Name string `json:"name"`
    Uri string `json:"uri"`
    Icon string `json:"icon"`
    Description string `json:"description"`
}

type BookmarkList struct {
    Models []BookmarkModel `json:"items"`
}

func (list *BookmarkList) AddModel(model *BookmarkModel) {
    l := len(list.Models);
    if l == 0 {
        model.Id = 1;
    } else {
        model.Id = list.Models[l-1].Id + 1;
    }

    list.Models = append(list.Models, *model);
}

func (list *BookmarkList) DeleteModel(id int) {
    list.Models = append(list.Models[:id], list.Models[id+1:]...);
    for i, m := range list.Models {
        if m.Id == id {
            list.Models = append(list.Models[:i], list.Models[i+1:]...);
            break;
        }
    }
}

func (list *BookmarkList) UpdateModel(id int, model *BookmarkModel) {
    for i, m := range list.Models {
        if m.Id == id {
            list.Models[i] = *model;
            break;
        }
    }
}

func (list *BookmarkList) QueryModel(id int, model *BookmarkModel) {
    for _, m := range list.Models {
        if m.Id == id {
            *model = m;
            break;
        }
    }
}

type BookmarkControler struct {
    modelList BookmarkList
}

func (c *BookmarkControler) List(request *restful.Request, response *restful.Response) {
    c.readFromFile();
    response.WriteAsJson(c.modelList);
}

func (c *BookmarkControler) Add(request *restful.Request, response *restful.Response) {
    c.readFromFile();
    model := BookmarkModel{};
    request.ReadEntity(&model);
    c.modelList.AddModel(&model);
    c.writeToFile();
}

func (c *BookmarkControler) Delete(request *restful.Request, response *restful.Response) {
    c.readFromFile();
    id := request.QueryParameter("id");
    i, err := strconv.Atoi(id);
    if err != nil {
        c.modelList.DeleteModel(i);
    }
    c.writeToFile()
}

func (c *BookmarkControler) Update(request *restful.Request, response *restful.Response) {
    c.readFromFile();
    id := request.PathParameter("id");
    model := BookmarkModel{};
    request.ReadEntity(&model);
    i, err := strconv.Atoi(id);
    if err != nil {
        c.modelList.UpdateModel(i, &model);
    }
    c.writeToFile();
}

func (c *BookmarkControler) Query(request *restful.Request, response *restful.Response) {
    c.readFromFile();
    id := request.PathParameter("id");
    i, err := strconv.Atoi(id);
    if err == nil {
        model := BookmarkModel{};
        c.modelList.QueryModel(i, &model);
        response.WriteAsJson(model);
    } else {
        response.AddHeader("Content-Type", "text/plain");
        response.WriteErrorString(http.StatusNotFound, "Id could not be found.");
    }
}

func (c *BookmarkControler) readFromFile() {
    path := "bookmark.json";
    bt, err := ioutil.ReadFile(path);
    if err != nil {
        log.Printf(err.Error());
        return;
    }

    err = json.Unmarshal(bt, &c.modelList);
    if err != nil {
        log.Printf(err.Error());
        return;
    }
}

func (c *BookmarkControler) writeToFile() {
    path := "bookmark.json";
    bt, err := json.MarshalIndent(c.modelList, "", "  ");
    if err != nil {
        log.Printf(err.Error());
        return;
    }

    ioutil.WriteFile(path, bt, 0644);
}

func main() {
    container := restful.NewContainer();
    webService := &restful.WebService{};

    bookmarkControler := BookmarkControler{};

    webService.Path("/").Produces(restful.MIME_JSON);

    webService.Route(webService.GET("/bookmark").
            To(bookmarkControler.List));
    webService.Route(webService.POST("/bookmark").
            To(bookmarkControler.Add));
    webService.Route(webService.DELETE("/bookmark/{id}").
            To(bookmarkControler.Delete).
            Param(webService.PathParameter("id", "identifier of the markbook").
                    DataType("string")));
    webService.Route(webService.PUT("/bookmark/{id}").
            To(bookmarkControler.Update).
            Param(webService.PathParameter("id", "identifier of the markbook").
                    DataType("string")));
    webService.Route(webService.GET("/bookmark/{id}").
            To(bookmarkControler.Query).
            Param(webService.PathParameter("id", "identifier of the markbook").
                    DataType("string")));

    container.Add(webService);
    
    log.Printf("start listening on localhost:8080");
    server := &http.Server{Addr: ":8080", Handler: container}
    log.Fatal(server.ListenAndServe());
}
