package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"text/template"
)

type Tasks struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Complete bool   `json:"complete"`
}

var (
	tasks  []Tasks
	nextID int
)

const fileName = "tasks.json"

func main() {
	Loader()
	http.HandleFunc("/", Handler)
	http.HandleFunc("/add", Add)
	http.HandleFunc("/complete", Complete)
	http.HandleFunc("/remove", Remove)

	fmt.Println("sever runnint at: http://localhost:8080")
	http.ListenAndServe("localhost:8080", nil)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "ERROR WHILE PARSING THE FILE: ", http.StatusBadRequest)
	}
	temp.Execute(w, tasks)
}

func Loader() {
	file, err := os.ReadFile(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			tasks = []Tasks{}
			nextID = 1
			return
		}
		log.Fatalln("ERRROR WHILE READING THE FILE:", err)
		return
	}
	err = json.Unmarshal(file, &tasks)
	if err != nil {
		log.Fatalln("ERROR WHILE UNMARSHALLING:", err)
		return
	}
	for _, task := range tasks {
		if task.ID >= nextID {
			nextID = task.ID + 1
		}
	}
}

func Saver() {
	data, err := json.MarshalIndent(tasks, "", "")
	if err != nil {
		log.Fatalln("ERROR WHILE MARSHALLING: ", err)
		return
	}
	err = os.WriteFile(fileName, data, 0644)
	if err != nil {
		log.Fatalln("ERROR WHILE WRITING THE FILE:", err)
		return
	}
}

func Add(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		taskName := r.FormValue("name")
		if taskName == "" {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		task := Tasks{
			ID:       nextID,
			Name:     taskName,
			Complete: false,
		}
		tasks = append(tasks, task)
		nextID++
		Saver()
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func Complete(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		idstr := r.FormValue("id")
		id, err := strconv.Atoi(idstr)
		if err != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		for i, task := range tasks {
			if task.ID == id {
				tasks[i].Complete = true
				break
			}
		}
		Saver()
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func Remove(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		idstr := r.FormValue("id")
		id, err := strconv.Atoi(idstr)
		if err != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		for i, task := range tasks {
			if task.ID == id {
				tasks = append(tasks[:i], tasks[i+1:]...)
				break
			}
		}
		Saver()
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
