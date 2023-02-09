package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"
	"tugas-3/connection"

	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
)

var Data = map[string]interface{}{
	"Title": "Personal Web",
}

type Project struct {
	Id                   int
	Title                string
	Description          string
	StartDate            time.Time
	EndDate              time.Time
	IsUsingReact         bool
	IsUsingNode          bool
	IsUsingAngular       bool
	IsUsingJavascript    bool
	Formatted_Start_Date string
	Formatted_End_Date   string
}

var Projects []Project

func main() {

	router := mux.NewRouter()
	connection.DatabaseConnect()
	router.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public/"))))

	router.HandleFunc("/", home).Methods("GET")
	router.HandleFunc("/addProjectForm", addProjectForm).Methods("GET")
	router.HandleFunc("/contact", contact).Methods("GET")
	router.HandleFunc("/addProject", addProject).Methods("POST")
	router.HandleFunc("/update-project/{id}", updateProject).Methods("POST")
	router.HandleFunc("/project-detail/{id}", projectDetail).Methods("GET")
	router.HandleFunc("/edit-project/{id}", editProject).Methods("GET")
	router.HandleFunc("/deleteProject/{id}", deleteProject).Methods("GET")
	fmt.Println("Server running on port 3000")
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

	rows, _ := connection.Conn.Query(context.Background(), "SELECT id, title, description, date(start_date), date(end_date), is_using_react, is_using_node, is_using_angular, is_using_javascript FROM public.tb_user")

	var result []Project
	for rows.Next() {
		var each = Project{}

		var err = rows.Scan(&each.Id, &each.Title, &each.Description, &each.StartDate, &each.EndDate, &each.IsUsingReact, &each.IsUsingNode, &each.IsUsingAngular, &each.IsUsingJavascript)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		each.Formatted_Start_Date = each.StartDate.Format("2006-01-02")
		each.Formatted_End_Date = each.EndDate.Format("2006-01-02")

		result = append(result, each)
	}
	fmt.Println(result)
	resp := map[string]interface{}{
		"Title":    Data,
		"Projects": result,
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
	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO tb_user(title, description, start_date, end_date, is_using_react, is_using_node, is_using_angular, is_using_javascript) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)", title, description, startDate, endDate, isUsingReact, isUsingNode, isUsingAngular, isUsingJavascript)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message : " + err.Error()))
		return
	}

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
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	fmt.Println("deleteProject function called")
	fmt.Println("Request method:", r.Method)
	w.Header().Set("Content-type", "text/html; charset=utf-8")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	fmt.Println(id)

	_, err := connection.Conn.Exec(context.Background(), "DELETE FROM tb_user WHERE id=$1", id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message : " + err.Error()))
		return
	}

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

	conn, err := pgx.Connect(context.Background(), "postgres://postgres:Rakhapostgre@localhost:5432/Personal-Web")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}
	defer conn.Close(context.Background())

	row := conn.QueryRow(context.Background(), "SELECT title, description, is_using_react, is_using_node, is_using_angular, is_using_javascript FROM public.tb_user WHERE id=$1", id)
	projectData := Project{}
	err = row.Scan(&projectData.Title, &projectData.Description, &projectData.IsUsingReact, &projectData.IsUsingNode, &projectData.IsUsingAngular, &projectData.IsUsingJavascript)
	if err == pgx.ErrNoRows {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("message : project not found"))
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	resp := map[string]interface{}{
		"Data":    Data,
		"Project": projectData,
		"Index":   id,
	}
	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, resp)
}

func updateProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	r.ParseForm()

	// Execute the update statement
	title := r.FormValue("title")
	description := r.FormValue("description")

	isUsingReact := r.FormValue("is_using_react")
	if isUsingReact == "" {
		isUsingReact = "false"
	}
	isUsingNode := r.FormValue("is_using_node")
	if isUsingNode == "" {
		isUsingNode = "false"
	}
	isUsingAngular := r.FormValue("is_using_angular")
	if isUsingAngular == "" {
		isUsingAngular = "false"
	}
	isUsingJavascript := r.FormValue("is_using_javascript")
	if isUsingJavascript == "" {
		isUsingJavascript = "false"
	}

	_, err := connection.Conn.Exec(context.Background(), "UPDATE tb_user SET title=$1, description=$2, is_using_react=$3, is_using_node=$4, is_using_angular=$5, is_using_javascript=$6 WHERE id=$7", title, description, isUsingReact, isUsingNode, isUsingAngular, isUsingJavascript, id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	// Redirect back to the project list page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
