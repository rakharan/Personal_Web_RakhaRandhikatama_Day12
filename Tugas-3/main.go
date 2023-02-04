package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/gorilla/mux"
)

var Data = map[string]interface{}{
	"Title": "Personal Web",
}

type Project struct {
	Title             string
	Description       string
	StartDate         string
	EndDate           string
	IsUsingReact      bool
	IsUsingNode       bool
	IsUsingAngular    bool
	IsUsingJavascript bool
}

var Projects = []Project{
	{
		Title:             "Kotakode Clone 1",
		Description:       "Kotakode merupakan platform komunitas bagi para pegiat IT di Indonesia dimana programmer dapat belajar dan berbagi wawasan seputar dunia IT terkini untuk mendukung memberikan pertumbuhan perekonomian di Indonesia.",
		StartDate:         time.Now().Format("2006-01-02"),
		EndDate:           time.Now().Format("2006-01-02"),
		IsUsingReact:      true,
		IsUsingNode:       true,
		IsUsingAngular:    true,
		IsUsingJavascript: true,
	},
}

func main() {

	router := mux.NewRouter()

	router.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public/"))))

	router.HandleFunc("/", home).Methods("GET")
	router.HandleFunc("/addProjectForm", addProjectForm).Methods("GET")
	router.HandleFunc("/contact", contact).Methods("GET")
	router.HandleFunc("/addProject", addProject).Methods("POST")
	router.HandleFunc("/update-project", updateProject).Methods("POST")
	router.HandleFunc("/delete-project/{id}", deleteProject).Methods("GET")
	router.HandleFunc("/project-detail/{id}", projectDetail).Methods("GET")
	router.HandleFunc("/edit-project/{id}", editProject).Methods("GET")
	fmt.Println("Server running on port 5000")
	http.ListenAndServe("localhost:3000", router)
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("pages/index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	resp := map[string]interface{}{
		"Title":    Data,
		"Projects": Projects,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, resp)
}

func addProjectForm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("pages/addProject.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, Data)
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("pages/contact.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, Data)
}

func addProject(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	title := r.PostForm.Get("title")
	description := r.PostForm.Get("description")
	startDate := r.PostForm.Get("startDate")
	endDate := r.PostForm.Get("endDate")
	isUsingReact := false
	isUsingNode := false
	isUsingJavascript := false
	isUsingAngular := false

	// Checked Tech Logic
	if r.FormValue("react") != "" {
		isUsingReact = true
	}
	if r.FormValue("node") != "" {
		isUsingNode = true
	}
	if r.FormValue("javascript") != "" {
		isUsingJavascript = true
	}
	if r.FormValue("angular") != "" {
		isUsingAngular = true
	}

	var newProject = Project{
		Title:             title,
		Description:       description,
		StartDate:         startDate,
		EndDate:           endDate,
		IsUsingReact:      isUsingReact,
		IsUsingAngular:    isUsingAngular,
		IsUsingNode:       isUsingNode,
		IsUsingJavascript: isUsingJavascript,
	}

	Projects = append(Projects, newProject)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func projectDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	// parsing template html
	var tmpl, err = template.ParseFiles("pages/projectDetail.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	ProjectDetail := Project{}
	for i, data := range Projects {
		if i == id {
			ProjectDetail = Project{
				Title:             data.Title,
				Description:       data.Description,
				StartDate:         data.StartDate,
				EndDate:           data.EndDate,
				IsUsingReact:      data.IsUsingReact,
				IsUsingNode:       data.IsUsingNode,
				IsUsingAngular:    data.IsUsingAngular,
				IsUsingJavascript: data.IsUsingJavascript,
			}
		}
	}

	resp := map[string]interface{}{
		"Data":    Data,
		"Project": ProjectDetail,
	}
	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, resp)
}

func deleteProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	fmt.Println(id)

	Projects = append(Projects[:id], Projects[id+1:]...)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func editProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var tmpl, err = template.ParseFiles("pages/editProject.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	ProjectData := Project{}
	for i, data := range Projects {
		if i == id {
			ProjectData = Project{
				Title:             data.Title,
				Description:       data.Description,
				StartDate:         data.StartDate,
				EndDate:           data.EndDate,
				IsUsingReact:      data.IsUsingReact,
				IsUsingNode:       data.IsUsingNode,
				IsUsingAngular:    data.IsUsingAngular,
				IsUsingJavascript: data.IsUsingJavascript,
			}
		}
	}

	resp := map[string]interface{}{
		"Data":    Data,
		"Project": ProjectData,
	}
	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, resp)

}

func updateProject(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	title := r.PostForm.Get("title")
	description := r.PostForm.Get("description")
	startDate := r.PostForm.Get("startDate")
	endDate := r.PostForm.Get("endDate")
	isUsingReact := false
	isUsingNode := false
	isUsingJavascript := false
	isUsingAngular := false

	// Checked Tech Logic
	if r.FormValue("react") != "" {
		isUsingReact = true
	}
	if r.FormValue("node") != "" {
		isUsingNode = true
	}
	if r.FormValue("javascript") != "" {
		isUsingJavascript = true
	}
	if r.FormValue("angular") != "" {
		isUsingAngular = true
	}

	pr := &Projects[id]
	(*pr).Title = title
	(*pr).Description = description
	(*pr).StartDate = startDate
	(*pr).EndDate = endDate
	(*pr).IsUsingReact = isUsingReact
	(*pr).IsUsingNode = isUsingNode
	(*pr).IsUsingJavascript = isUsingJavascript
	(*pr).IsUsingAngular = isUsingAngular
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}
