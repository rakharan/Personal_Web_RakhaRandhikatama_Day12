package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"
	"tugas-3/connection"
	"tugas-3/middleware"

	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type MetaData struct {
	Title    string
	IsLogin  bool
	Username string
	IsAuthor []string
}

var Data = MetaData{
	Title: "Personal Web",
}

type Users struct {
	Id       int
	Name     string
	Email    string
	Password string
}

var user = Users{}

type Project struct {
	Id                   int
	Title                string
	Description          string
	StartDate            time.Time
	EndDate              time.Time
	Technologies         []string
	Formatted_Start_Date string
	Formatted_End_Date   string
	Author               string
	Image                string
	Time_Difference      int
}

var Projects []Project

func main() {

	router := mux.NewRouter()
	connection.DatabaseConnect()
	router.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public/"))))
	router.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads/"))))

	router.HandleFunc("/", home).Methods("GET")
	router.HandleFunc("/addProjectForm", addProjectForm).Methods("GET")
	router.HandleFunc("/contact", contact).Methods("GET")
	router.HandleFunc("/addProject", middleware.UploadFile(addProject)).Methods("POST")
	router.HandleFunc("/update-project/{id}", middleware.UploadFile(updateProject)).Methods("POST")
	router.HandleFunc("/project-detail/{id}", projectDetail).Methods("GET")
	router.HandleFunc("/edit-project/{id}", editProject).Methods("GET")
	router.HandleFunc("/deleteProject/{id}", deleteProject).Methods("GET")
	router.HandleFunc("/register", registerForm).Methods("GET")
	router.HandleFunc("/register", register).Methods("Post")
	router.HandleFunc("/login", loginForm).Methods("GET")
	router.HandleFunc("/login", login).Methods("POST")
	router.HandleFunc("/logout", logout).Methods("GET")
	router.HandleFunc("/profile", profile).Methods("GET")
	fmt.Println("Server running on port 3000")
	http.ListenAndServe("localhost:3000", router)
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("pages/index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	rows, _ := connection.Conn.Query(context.Background(), "SELECT tb_project.id, title, description, date(start_date), date(end_date), technologies, image,  tb_user.name as author FROM tb_project LEFT JOIN tb_user ON tb_project.author_id = tb_user.id ORDER BY id ASC")

	var result []Project
	for rows.Next() {
		var each = Project{}

		var err = rows.Scan(&each.Id, &each.Title, &each.Description, &each.StartDate, &each.EndDate, &each.Technologies, &each.Image, &each.Author)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		each.Formatted_Start_Date = each.StartDate.Format("2006-01-02")
		each.Formatted_End_Date = each.EndDate.Format("2006-01-02")
		diff := each.EndDate.Sub(each.StartDate)
		seconds := int(diff.Seconds())
		minutes := seconds / 60
		hours := minutes / 60
		days := hours / 24
		each.Time_Difference = days
		result = append(result, each)
	}
	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	if session.Values["IsLogin"] != true {
		Data.IsLogin = false
	} else {
		Data.IsLogin = session.Values["IsLogin"].(bool)
		Data.Username = session.Values["Name"].(string)
	}

	resp := map[string]interface{}{
		"Projects": result,
		"Data":     Data,
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
	dataContex := r.Context().Value("dataFile")
	image := dataContex.(string)

	title := r.PostForm.Get("title")
	description := r.PostForm.Get("description")
	startDate := r.PostForm.Get("startDate")
	endDate := r.PostForm.Get("endDate")
	technology := r.Form["technology"]

	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO tb_project(title, description, start_date, end_date, technologies, image, author_id) VALUES ($1, $2, $3, $4, $5, $6, $7)", title, description, startDate, endDate, technology, image, user.Id)
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
	err = connection.Conn.QueryRow(context.Background(), "SELECT title, description, date(start_date), date(end_date), technologies, image FROM public.tb_project WHERE id=$1", id).Scan(&ProjectDetail.Title, &ProjectDetail.Description, &ProjectDetail.StartDate, &ProjectDetail.EndDate, &ProjectDetail.Technologies, &ProjectDetail.Image)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message : " + err.Error()))
		return
	}
	ProjectDetail.Formatted_Start_Date = ProjectDetail.StartDate.Format("2006-01-02")
	ProjectDetail.Formatted_End_Date = ProjectDetail.EndDate.Format("2006-01-02")
	diff := ProjectDetail.EndDate.Sub(ProjectDetail.StartDate)
	seconds := int(diff.Seconds())
	minutes := seconds / 60
	hours := minutes / 60
	days := hours / 24
	resp := map[string]interface{}{
		"Data":           Data,
		"Project":        ProjectDetail,
		"Index":          id,
		"TimeDifference": days,
	}
	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, resp)
}

func deleteProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	w.Header().Set("Content-type", "text/html; charset=utf-8")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	_, err := connection.Conn.Exec(context.Background(), "DELETE FROM tb_project WHERE id=$1", id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func editProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
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

	row := conn.QueryRow(context.Background(), "SELECT title, description, date(start_date), date(end_date), technologies FROM public.tb_project WHERE id=$1", id)
	projectData := Project{}
	err = row.Scan(&projectData.Title, &projectData.Description, &projectData.StartDate, &projectData.EndDate, &projectData.Technologies)
	if err == pgx.ErrNoRows {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("message : project not found"))
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	type isUsingTechnology struct {
		IsUsingReact      bool
		IsUsingJavascript bool
		IsUsingAngular    bool
		IsUsingNode       bool
	}

	var technologyUsed = isUsingTechnology{}
	var techUsed []string
	for _, data := range projectData.Technologies {
		techUsed = append(techUsed, data)
	}

	for _, tech := range techUsed {
		if tech == "react" {
			technologyUsed.IsUsingReact = true
		}
		if tech == "js" {
			technologyUsed.IsUsingJavascript = true
		}
		if tech == "node" {
			technologyUsed.IsUsingNode = true
		}
		if tech == "angular" {
			technologyUsed.IsUsingAngular = true
		}
	}

	resp := map[string]interface{}{
		"Data":       Data,
		"Project":    projectData,
		"Index":      id,
		"Technology": technologyUsed,
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
	startDate := r.PostForm.Get("startDate")
	endDate := r.PostForm.Get("endDate")
	technologies := r.Form["technology"]
	dataContex := r.Context().Value("dataFile")
	image := dataContex.(string)

	_, err := connection.Conn.Exec(context.Background(), "UPDATE tb_project SET title=$1, description=$2, start_date=$3, end_date=$4, technologies=$5, image=$6 WHERE id=$7", title, description, startDate, endDate, technologies, image, id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	// Redirect back to the project list page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func registerForm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("pages/registerForm.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, Data)
}

func register(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	name := r.PostForm.Get("name")
	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")

	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO public.tb_user(name, email, password) VALUES ($1, $2, $3);", name, email, passwordHash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	http.Redirect(w, r, "/login", http.StatusMovedPermanently)
}

func loginForm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("pages/loginForm.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, Data)
}

func login(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")

	user = Users{}

	err = connection.Conn.QueryRow(context.Background(), "SELECT Id, email, name, password FROM tb_user WHERE email=$1", email).Scan(&user.Id, &user.Email, &user.Name, &user.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	session.Values["IsLogin"] = true
	session.Values["Name"] = user.Name
	session.Options.MaxAge = 10800

	session.AddFlash("Login succes", "message")
	session.Save(r, w)
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func logout(w http.ResponseWriter, r *http.Request) {
	log.Println("logout function called")
	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")
	session.Values["IsLogin"] = false
	session.Options.MaxAge = -1
	session.Save(r, w)

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func profile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("pages/profile.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}
	rows, _ := connection.Conn.Query(context.Background(), "SELECT tb_project.id, title, description, date(start_date), date(end_date), technologies, image,  tb_user.name as author FROM tb_project LEFT JOIN tb_user ON tb_project.author_id = tb_user.id ORDER BY id ASC")

	var result []Project
	for rows.Next() {
		var each = Project{}

		var err = rows.Scan(&each.Id, &each.Title, &each.Description, &each.StartDate, &each.EndDate, &each.Technologies, &each.Image, &each.Author)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		each.Formatted_Start_Date = each.StartDate.Format("2006-01-02")
		each.Formatted_End_Date = each.EndDate.Format("2006-01-02")
		result = append(result, each)

	}
	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	if session.Values["IsLogin"] != true {
		Data.IsLogin = false
	} else {
		Data.IsLogin = session.Values["IsLogin"].(bool)
		Data.Username = session.Values["Name"].(string)
	}

	resp := map[string]interface{}{
		"Projects": result,
		"Data":     Data,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, resp)
}
