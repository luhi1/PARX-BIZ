package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	"html/template"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type DisplayError struct {
	ErrorDescription string
}

// TeacherPageHandlers Consider creating a generic interface for both teacher and student to implement.
type TeacherPageHandlers interface {
	GETHandler(writer http.ResponseWriter, request *http.Request)
	POSTHandler(writer http.ResponseWriter, request *http.Request)
	valHandler(writer http.ResponseWriter, request *http.Request)
	dataVal(requestMethod string) bool
}

var db *sql.DB
var userInfo = UserData{}

// Start server run, files, and other shit.
func main() {
	partnerInfo := PartnerInfo{}
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/login", userInfo.GETHandler)

	http.HandleFunc("/signup", userInfo.POSTHandler)

	http.HandleFunc("/userValidation/", userInfo.valHandler)

	http.HandleFunc("/logout", func(writer http.ResponseWriter, request *http.Request) {
		userInfo = UserData{}
		http.Redirect(writer, request, "./home", 307)
	})

	http.HandleFunc("/teacherPartners", partnerInfo.GETHandler)
	http.HandleFunc("/create", partnerInfo.POSTHandler)
	http.HandleFunc("/create/submit", partnerInfo.valHandler)
	http.HandleFunc("/create/remove", partnerInfo.removeHandler)
	http.HandleFunc("/create/update", partnerInfo.updateHandler)
	http.HandleFunc("/home", partnerInfo.GETHandler)

	/*
		http.HandleFunc("/teacherCreateEvent", eventInfo.POSTHandler)

		http.HandleFunc("/winners", winners.GETHandler)
		http.HandleFunc("/prizes", prize.GETHandler)
		http.HandleFunc("/eventValidation/", eventInfo.valHandler)
		http.HandleFunc("/removeEvent", eventInfo.removeHandler)
		http.HandleFunc("/teacherCreateEvent/createEvent", eventInfo.createEvent)
		http.HandleFunc("/reroll", winners.valHandler)
		http.HandleFunc("/prizeChecking", prize.valHandler)
		http.HandleFunc("/createPrize", prize.POSTHandler)
		http.HandleFunc("/createPrizes", prize.createPrize)
		http.HandleFunc("/studentEvents", studentEventInformation.GETStudentHandler)
		http.HandleFunc("/dropOut", studentEventInformation.dropOutHandler)
		http.HandleFunc("/home", homeData.GETStudentHandler)
		http.HandleFunc("/quarterReport", winners.report)
		http.HandleFunc("/studentSignupEvent", studentEventInformation.studentSignupEventHandler)

		http.HandleFunc("/qna", func(writer http.ResponseWriter, request *http.Request) {
			multiTplExec(writer, "qna.gohtml", nil, "teacher_partners.gohtml")
		})
		http.HandleFunc("/bugs", func(writer http.ResponseWriter, request *http.Request) {
			request.ParseForm()
			db.Exec("insert into bugs(bugs) values(?)", request.FormValue("ProblemDesc"))
			http.Redirect(writer, request, "/home", 307)
		})


	*/
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/" {
			err := tplExec(writer, "error.gohtml", DisplayError{"Error 404"})
			if err != nil {
				tplExec(writer, "error.gohtml", DisplayError{err.Error()})
			}
		} else {
			http.Redirect(writer, request, "./home", 301)
		}
	})

	initDB, err := sql.Open("mysql", "root:test@tcp(127.0.0.1:3306)/fblacnp")

	if err != nil {
		panic(err)
	}

	defer db.Close()

	db = initDB
	fmt.Println("Connected to DB")

	fmt.Println("Server is running on port 8082")
	err = http.ListenAndServe(":8082", nil)
	if err != nil {
		fmt.Println("Error starting server, aborting tasks")
		panic(err)
	}
}

func tplExec(w http.ResponseWriter, filename string, information any) error {
	temp := template.Must(template.ParseFiles("./WebPages/" + filename))

	err := temp.Execute(w, information)
	//@TODO: REMOVE
	if err != nil {
		return err
	}
	return nil
}

func hashPswd(pwd string) string {
	hasher := sha256.New()
	hasher.Write([]byte(pwd))
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return sha
}
