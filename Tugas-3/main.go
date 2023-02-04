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
	Title       string
	Description string
}

var Projects = []Project{
	{
		Title:       "Kotakode Clone 1",
		Description: "Kotakode merupakan platform komunitas bagi para pegiat IT di Indonesia dimana programmer dapat belajar dan berbagi wawasan seputar dunia IT terkini untuk mendukung memberikan pertumbuhan perekonomian di Indonesia.",
	},
}

func main() {

	router := mux.NewRouter()

	router.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public/"))))

	router.HandleFunc("/", home).Methods("GET")
	router.HandleFunc("/addProjectForm", addProjectForm).Methods("GET")
	router.HandleFunc("/contact", contact).Methods("GET")
	router.HandleFunc("/addProject", addProject).Methods("POST")
	router.HandleFunc("/delete-project/{id}", deleteProject).Methods("GET")
	fmt.Println("Server running on port 5000")
	http.ListenAndServe("localhost:8000", router)
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
	fmt.Println(time.Now().String())

	var newProject = Project{
		Title:       title,
		Description: description,
		// startDate:   startDate,
		// endDate:     endDate,
	}

	Projects = append(Projects, newProject)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func deleteProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	fmt.Println(id)

	Projects = append(Projects[:id], Projects[id+1:]...)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}
