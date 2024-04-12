package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type PartnerInfo struct {
	ID           int
	Name         string
	Type         string
	Email        string
	Phone_Number string
	Resources    []string
	Active       int8
}

var partners []PartnerInfo

func (p *PartnerInfo) GETHandler(writer http.ResponseWriter, request *http.Request) {

	if (userInfo != UserData{} && strings.TrimPrefix(request.URL.Path, "/") != "teacherPartners") {
		http.Redirect(writer, request, "./teacherPartners", 303)
	} else if (userInfo == UserData{} && strings.TrimPrefix(request.URL.Path, "/") != "home") {
		http.Redirect(writer, request, "./home", 303)
	}

	partners = []PartnerInfo{}
	rows, err := db.Query("select Partners.id, Partners.name, Representatives.email, Representatives.phone, Partner_Types.name, Partners.active from Partners join Partner_Types on Partners.type = Partner_Types.id join Representatives on Partners.representative = Representatives.id")

	if err != nil {
		fmt.Println(err)
		return
	}
	for rows.Next() {
		*p = PartnerInfo{}
		rows.Scan(&p.ID, &p.Name, &p.Email, &p.Phone_Number, &p.Type, &p.Active)
		if p.Active == 0 {
			continue
		}

		rowResources, err := db.Query("select Resources.info from Resources where Resources.partner = ?", p.ID)
		if err != nil {
			fmt.Println(err)
			return
		}

		resources := []string{}
		for rowResources.Next() {
			var resource string
			rowResources.Scan(&resource)
			resources = append(resources, resource)
		}
		p.Resources = resources
		partners = append(partners, *p)
	}
	err = tplExec(writer, strings.TrimPrefix(request.URL.Path, "/")+".gohtml", partners)
	//@TODO: REMOVE
	if err != nil {
		return
	}
}
func (p *PartnerInfo) POSTHandler(writer http.ResponseWriter, request *http.Request) {
	err := tplExec(writer, "create.gohtml", nil)
	//@TODO: REMOVE
	if err != nil {
		return
	}
}

func (p *PartnerInfo) valHandler(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		http.Redirect(writer, request, "./error", 303)
		return
	}

	p.Name = request.FormValue("Name of Organization")
	p.Type = request.FormValue("Type of Organization")
	p.Email = request.FormValue("Email")
	p.Phone_Number = request.FormValue("Phone Number")
	//p.Resources
	p.Active = 1

	index := 0
	var resources []string
	for request.FormValue("Resource"+strconv.Itoa(index)) != "" {
		resources = append(resources, request.FormValue("Resource"+strconv.Itoa(index)))
		index++
	}
	p.Resources = resources

	if p.dataVal() {
		check := db.QueryRow("select id from Partner_Types where name = ?", p.Type)
		var partnerTypesID int
		err := check.Scan(&partnerTypesID)
		if err != nil {
			fmt.Println(err)
			return
		}

		check = db.QueryRow("select id from Representatives where email = ? and phone = ?", p.Email, p.Phone_Number)
		var repID int64
		err = check.Scan(&repID)
		if err != nil {
			repID = -1
		}
		if repID == -1 {
			result, err := db.Exec("insert into Representatives(email, phone) values (?, ?)",
				p.Email, p.Phone_Number,
			)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(result.RowsAffected())
			repID, _ = result.LastInsertId()
		}

		check = db.QueryRow("select id from Partners where name = ?", p.Name)
		var partID int64
		err = check.Scan(&partID)
		if err != nil {
			partID = -1
		}

		if partID == -1 {
			result, err := db.Exec("insert into Partners(`name`, representative, `type`, `active`) values (?, ?, ?, ?)",
				p.Name, repID, partnerTypesID, p.Active,
			)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(result.RowsAffected())
			partID, _ = result.LastInsertId()
		}

		for i := 0; i < len(p.Resources); i++ {
			check = db.QueryRow("select id from Resources where info = ?", p.Resources[i])
			var resID int
			err = check.Scan(&resID)
			if err != nil {
				vector, _ := db.Exec("insert into Resources(partner, info) values(?, ?)",
					partID, p.Resources[i])
				fmt.Println(vector.RowsAffected())
			}
		}
	}
	http.Redirect(writer, request, "../teacherPartners", 307)
}

func (p *PartnerInfo) dataVal() bool {
	if p.Name == "" || p.Type == "" || p.Email == "" || p.Phone_Number == "" {
		return false
	}
	for i := 0; i < len(p.Resources); i++ {
		if p.Resources[i] == "" {
			return false
		}
	}
	return true
}

func (p *PartnerInfo) removeHandler(writer http.ResponseWriter, request *http.Request) {
	var err error
	err = request.ParseForm()
	if err != nil {
		http.Redirect(writer, request, "./error", 303)
		return
	}

	p.Name = request.FormValue("Name of Organization")
	p.Type = request.FormValue("Type of Organization")
	p.Email = request.FormValue("Email")
	p.Phone_Number = request.FormValue("Phone Number")
	p.Active = 1

	check := db.QueryRow("select id from Partner_Types where name = ?", p.Type)
	var partnerTypesID int
	err = check.Scan(&partnerTypesID)
	if err != nil {
		fmt.Println(err)
		return
	}
	check = db.QueryRow("select id from Representatives where email = ? and phone = ?", p.Email, p.Phone_Number)
	var repID int64
	err = check.Scan(&repID)
	if err != nil {
		repID = -1
	}

	for i := 0; i < len(partners); i++ {
		if partners[i].Name == p.Name {
			p.ID = partners[i].ID
		}
	}

	exec, err := db.Exec("update Partners set active = 0 where id = ?", p.ID)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(exec.RowsAffected())

	resources, err := db.Exec("delete from Resources where partner = ?",
		p.ID)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(resources.RowsAffected())
	http.Redirect(writer, request, "../teacherPartners", 307)
}

/*
func (e *PartnerInfo) createEvent(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		fmt.Println(err)
		fmt.Println(request.Form)
		return
	}
	e.EventName = request.FormValue("EventName")
	e.Points, _ = strconv.Atoi(request.FormValue("Points"))
	e.EventDescription = request.FormValue("EventDescription")
	e.EventDate = request.FormValue("EventDate")
	e.RoomNumber, _ = strconv.Atoi(request.FormValue("RoomNumber"))
	e.AdvisorNames = request.FormValue("AdvisorNames")
	e.Location = request.FormValue("Location")
	e.LocationDescription = request.FormValue("LocationDescription")
	e.Sport = request.FormValue("Sport")
	e.SportDescription = request.FormValue("SportDescription")
	e.Active = true
	if e.dataVal("") {
		check := db.QueryRow("select ID from sports where SportName = ?", e.Sport)
		var sportID int
		check.Scan(&sportID)
		if err != nil {
			sportID = -1
		}
		if sportID == -1 {
			insert, _ := db.Exec("insert into sports(sportname, sportdescription) values(?, ?);", e.Sport, e.SportDescription)
			fmt.Println(insert.RowsAffected())
			getSportID, err := insert.LastInsertId()
			if err != nil {
				return
			}
			sportID = int(getSportID)
		}
		result, err := db.Exec("insert into events(eventname, Points, eventdescription, eventdate, roomnumber, advisors, location, locationdescription, sportid, active) VALUES (?,?,?,?,?,?,?,?,?,?)",
			e.EventName, e.Points, e.EventDescription, e.EventDate, e.RoomNumber, e.AdvisorNames, e.Location, e.LocationDescription, sportID, e.Active)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(result.RowsAffected())

	}
	http.Redirect(writer, request, "../teacherEvents", 307)
}

*/
