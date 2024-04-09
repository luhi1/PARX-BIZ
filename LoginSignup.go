package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type UserData struct {
	ID           int
	Username     string
	Password     string
	Real_Name    string
	Program_Area string
	valid        DisplayError
}

func (u *UserData) GETHandler(writer http.ResponseWriter, request *http.Request) {
	err := tplExec(writer, "login.gohtml", u.valid)

	if err != nil {
		u.valid = DisplayError{err.Error()}
		return
	}
	u.valid = DisplayError{""}
}

func (u *UserData) POSTHandler(writer http.ResponseWriter, request *http.Request) {
	err := tplExec(writer, "signup.gohtml", u.valid)
	if err != nil {
		u.valid = DisplayError{err.Error()}
		return
	}
	u.valid = DisplayError{""}
}

func (u *UserData) valHandler(writer http.ResponseWriter, request *http.Request) {
	var err error
	err = request.ParseForm()
	if err != nil {
		http.Redirect(writer, request, "./error", 303)
		return
	}
	u.ID, err = strconv.Atoi(request.FormValue("ID"))
	u.Username = request.FormValue("Username")
	u.Password = hashPswd(request.FormValue("Password"))
	u.Real_Name = request.FormValue("Real_Name")
	u.Program_Area = request.FormValue("Program_Area")
	if strings.TrimPrefix(request.URL.Path, "/userValidation/") != "login" && err != nil {
		http.Redirect(writer, request, "../error", 303)
		return
	}

	if u.dataVal(strings.TrimPrefix(request.URL.Path, "/userValidation/")) {
		insert := db.QueryRow("select Users.id, Users.username, Users.password, Users.real_name,"+
			" Program_Areas.name from Users join Program_Areas on Users.program_area = Program_Areas.id where Users.id = ? && Users.username = ? && Users.password = ?;", strconv.Itoa(u.ID), u.Username, u.Password)
		insert.Scan(&u.ID, &u.Username, &u.Password, u.Real_Name, u.Program_Area)
		if u.ID == 0 && u.Username == "" && u.Password == "" && u.Real_Name == "" && u.Program_Area == "" {
			u.valid = DisplayError{"Invalid Credentials"}
			if strings.TrimPrefix(request.URL.Path, "/userValidation/") == "signup" {
				http.Redirect(writer, request, "../signup", 303)
			} else {
				http.Redirect(writer, request, "../login", 303)
			}
		} else {
			if u.ID == 1 {
				http.Redirect(writer, request, "../teacherEvents", 307)
			} else {
				http.Redirect(writer, request, "../home", 307)
			}
		}
	} else {
		u.valid = DisplayError{"Invalid Credentials"}
		if strings.TrimPrefix(request.URL.Path, "/userValidation/") == "signup" {
			http.Redirect(writer, request, "../signup", 303)
		} else {
			http.Redirect(writer, request, "../login", 303)
		}
	}
}

func (u *UserData) dataVal(requestMethod string) bool {
	valid := false

	//A bit of crazy code here, but this first conditional serves to make sure no one can access webpages early!
	//The latter half allows the user to login and populates their data from the DB
	if (*u != UserData{}) &&
		(u.Username != "" &&
			u.Password != hashPswd("")) {

		valid = true
		if requestMethod == "login" {
			updateUser := db.QueryRow("select Users.id, Users.real_name, Program_Areas.name from Users join Program_Areas on Users.program_area = Program_Areas.id where Users.username = ? && Users.password = ?;", u.Username, u.Password)
			err := updateUser.Scan(&u.ID, &u.Real_Name, &u.Program_Area)
			if err != nil {
				return false
			}
			return true
		}
	}

	programAreaInt, err := strconv.Atoi(u.Program_Area)
	if err != nil {
		return false
	}

	if requestMethod == "signup" && ((u.ID <= 1 || u.ID > 9999999) || u.Real_Name == "" ||
		(programAreaInt > 7 || programAreaInt < 1)) {
		valid = false
	}

	if valid && requestMethod == "signup" {
		result, err := db.Exec(
			"insert into Users(id, username, password, real_name, program_area) values(?, ?, ?, ?, ?);",
			u.ID,
			u.Username,
			u.Password,
			u.Real_Name,
			u.Program_Area,
		)
		if err != nil {
			return false
		}

		getPA := db.QueryRow("select name from Program_Areas where id = ?;", programAreaInt)
		getPA.Scan(&u.Program_Area)
		fmt.Println(result.RowsAffected())
	}
	return valid
}
