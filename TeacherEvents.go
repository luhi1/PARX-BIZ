package main

import (
	"fmt"
	"net/http"
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

func (p *PartnerInfo) GETHandler(writer http.ResponseWriter, request *http.Request) {
	partners := []PartnerInfo{}

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
	err = tplExec(writer, "teacher_partners.gohtml", partners)
	//@TODO: REMOVE
	if err != nil {
		return
	}
}

/*
func (e *PartnerInfo) POSTHandler(writer http.ResponseWriter, request *http.Request) {
	err := tplExec(writer, "teacher_create_event.gohtml", nil)
	//@TODO: REMOVE
	if err != nil {
		return
	}
}

func (e *PartnerInfo) valHandler(writer http.ResponseWriter, request *http.Request) {
	var err error
	err = request.ParseForm()
	if err != nil {
		http.Redirect(writer, request, "./error", 303)
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
	for i := 0; i < len(request.Form["Attendee"]); i++ {
		currentHomie, _ := strconv.Atoi(request.Form["Attendee"][i])
		e.Attendance = append(e.Attendance, StudentAttendance{StudentNumber: currentHomie, Attended: "true"})
	}
	if e.dataVal(strings.TrimPrefix(request.URL.Path, "/eventValidation/")) {
		check := db.QueryRow("select ID from sports where SportName = ?", e.Sport)
		var sportID int
		err := check.Scan(&sportID)
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
		sID := strconv.Itoa(sportID)
		points := strconv.Itoa(e.Points)
		roomNumber := strconv.Itoa(e.RoomNumber)

		result, err := db.Exec("update events set events.Points = ?, EventDescription = ?, EventDate = ?, RoomNumber = ?, Advisors = ?, Location = ?, LocationDescription = ?, SportID = ? where events.EventName = ?",
			points, e.EventDescription, e.EventDate, roomNumber, e.AdvisorNames, e.Location, e.LocationDescription, sID, e.EventName,
		)
		if err != nil {
			fmt.Println(err)
		}

		insert := db.QueryRow("select EventID from events where EventName = ?;", e.EventName)
		insert.Scan(&e.EventID)

		minion, _ := db.Exec("update userevents set Attended = 'false' where EventID = ?;", e.EventID)
		fmt.Println(minion.RowsAffected())
		for i := 0; i < len(e.Attendance); i++ {
			//Change it to add Points when the homies sign up for an event.
			vector, _ := db.Exec("update userevents set Attended = 'true' where EventID = ? and UserID = ?", e.EventID, e.Attendance[i].StudentNumber)
			fmt.Println(vector.RowsAffected())
		}
		fmt.Println(e.Attendance)
		fmt.Println(result.RowsAffected())
	}
	http.Redirect(writer, request, "../teacherEvents", 307)
}

func (e *PartnerInfo) dataVal(requestMethod string) bool {
	fmt.Println(e)
	if e.Points < 0 || e.EventDescription == "" || e.EventDate == "" || e.RoomNumber < 1 || e.AdvisorNames == "" || e.Location == "" {
		return false
	}
	for i := 0; i < len(e.Attendance); i++ {
		if e.Attendance[i].StudentNumber < 1 {
			return false
		}
	}
	return true
}

func (e *PartnerInfo) removeHandler(writer http.ResponseWriter, request *http.Request) {
	e.EventName = request.FormValue("EventName")
	eventID := db.QueryRow("select EventID from events where EventName = ?", e.EventName)
	eventID.Scan(&e.EventID)
	exec, err := db.Exec("update events set Active = 0 where EventID = ?", e.EventID)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(exec.RowsAffected())

	insert, _ := db.Query("select UserID from userevents where Attended = 'false' and EventID = ?", e.EventID)

	var subtracters []int
	for insert.Next() {
		var currentSubtracter int
		insert.Scan(&currentSubtracter)
		subtracters = append(subtracters, currentSubtracter)
	}
	for i := 0; i < len(subtracters); i++ {
		addition, _ := db.Exec("update users set Points = Points-10 where userID = ?", subtracters[i])
		fmt.Println(addition.RowsAffected())
	}
	http.Redirect(writer, request, "./teacherEvents", 307)
}

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
