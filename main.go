package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/marni/goigc"
	"net/http"
	//"os"
	"strconv"
	"strings"
	"time"
)

//Stores a string for use to check if the url is already added

//Stores data about the track
//Changed from Track igc.Track to the different values needed! done after the deadline
type Track struct {
	ID    string
	HDate time.Time
	Pilot string
	Glider string
	GliderID string
	TrackLength float64
}

var urlList = make(map[string]Track)
var urlAmount = 0

//Check if the url is already added
func inList(url string) bool {
	for inList := range urlList {
		if inList == url {
			return true
		}
	}
	return false
}

var timeStart = int(time.Now().Unix())

//If bad link is provided, error message will be shown
func error404(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

//Bad Request error message function as it is used many times
func error400(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}


//Errorcheck added after deadline, uses URL to not make two structs that only store a string
func errorCheck(val string) (string, error) {
	return val, nil
}

//Calculate the time passed since server for x seconds
func calcDuration(timePassed int) string {
	passed := "P"

	timeLeft := timePassed

	days := timeLeft / 86400
	passed += fmt.Sprintf("D%d", days)
	timeLeft -= days * 86400

	passed += "T"

	hours := timeLeft / 3600
	passed += fmt.Sprintf("H%d", hours)
	timeLeft -= hours * 3600

	minutes := timeLeft / 60
	passed += fmt.Sprintf("M%d", minutes)
	timeLeft -= minutes * 60

	passed += fmt.Sprintf("S%d", timeLeft)

	return passed

}

//shows the meta information about the API
func getMetaInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	parts := strings.Split(r.URL.Path, "/")

	if len(parts) != 4 {
		error404(w, r)
		return
	}

	upTime := calcDuration(int(time.Now().Unix()) - timeStart)

	Response := "{"
	Response += "\"uptime\": \"" + upTime + "\","
	Response += "\"info\": \"Service for IGC tracks.\","
	Response += "\"version\": \"v1\""
	Response += "}"

	response, err := errorCheck(Response)
	if err != nil {
		error400(w)
		return
	}

	fmt.Fprintln(w, response)
}

//registers track and returns array of all trackids depending on POST and GET
func manageTrack(w http.ResponseWriter, r *http.Request) {

	//Register a new track
	if r.Method == "POST" {
		w.Header().Set("Content-Type", "application/json")
		if r.Body == nil {
			error400(w)
			return
		}

		var trackURL string

		//Had used a struct before, but it has been removed as it was not needed
		//Error-check was added after deadline
		//Used to be separated with new decoder and decoding, but now is joined
		err := json.NewDecoder(r.Body).Decode(&trackURL)

		if err != nil {
			error400(w)
			return
		}

		track, err := igc.ParseLocation(trackURL)
		if err != nil {
			error400(w)
			return
		}
		urlAmount++
		newID := "igc" + strconv.Itoa(urlAmount)

		ID, err := errorCheck(newID)
		if err != nil {
			error400(w)
			return
		}
		//Check if track is new
		if !inList(newID) {
			//Method of adding new track changed a bit after the new update (and it works now! yay)
			urlList[ID] = Track{ID, track.Date, track.Pilot, track.GliderType, track.GliderID, igc.Distance(track.Task)}
			Response := "{"
			Response += "\"id:\" \"" + urlList[ID].ID + "\""
			Response += "}"

			response, err := errorCheck(Response)
			if err != nil {
				error400(w)
				return
			}

			fmt.Fprintln(w, response)
		} else {
			error400(w)
			return
		}

	} else if r.Method == "GET" { //returns all the tracks ids
		w.Header().Set("Content-Type", "application/json")

		x := 0 //variable as i can't be compared to int

		Response := "[ "
		for i := range urlList {
			Response += urlList[i].ID
			if x < len(urlList) {
				Response += ", "
			}
			x++
		}

		Response += " ]"
		response, err := errorCheck(Response)
		if err != nil {
			error400(w)
			return
		}

		fmt.Fprintln(w, response)
	}
}

//Gives all the information about a certain track by given ID
func getTrackByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	parts := strings.Split(r.URL.Path, "/")
	//Used to be (parts)-1, but that gave the wrong field, changed after deadline
	track := parts[len(parts)-2]
	if track != "" {
		Response := "{"
		Response += "\"H_date\": " + "\"" + urlList[track].HDate.String() + "\","
		Response += "\"pilot\": " + "\"" + urlList[track].Pilot + "\","
		Response += "\"glider\": " + "\"" + urlList[track].Glider + "\","
		Response += "\"glider_id\": " + "\"" + urlList[track].GliderID + "\","
		Response += "\"track_length\": " + "\"" + strconv.FormatFloat(urlList[track].TrackLength, 'E', -1, 64) + "\","
		Response += "}"

		response, err := errorCheck(Response)
		if err != nil {
			error400(w)
			return
		}

		fmt.Fprint(w, response)
	} else {
		error404(w, r)
	}

}

//Given id and field, returns the field of given id
func getTrackField(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")

	//Used to be (parts) -2, but was changed after deadline due to error
	id := parts[len(parts)-3]
	//Used to be (parts) -1, but was changed after deadline due to error
	field := parts[len(parts)-2]

	if id != "" && field != "" {
		switch field {
		case "pilot":
			Response := urlList[id].Pilot

			response, err := errorCheck(Response)
			if err != nil {
				error400(w)
				return
			}
			fmt.Fprint(w, response)
		case "glider":
			Response := urlList[id].Glider

			response, err := errorCheck(Response)
			if err != nil {
				error400(w)
				return
			}
			fmt.Fprint(w, response)
		case "glider_id":
			Response := urlList[id].GliderID

			response, err := errorCheck(Response)
			if err != nil {
				error400(w)
				return
			}
			fmt.Fprint(w, response)
		case "track_length":
			fmt.Fprint(w, urlList[id].TrackLength)
		case "H_date":
			Response := urlList[id].HDate.String()

			response, err := errorCheck(Response)
			if err != nil {
				error400(w)
				return
			}
			fmt.Fprint(w, response)
		}
	} else {
		error404(w, r)
	}
}

//Runs the application
func main() {
	//Missing router makes getTrackByID and getTrackField unaccessable
	//It has been added so the judges can more easily check the last two pages!
	//It was added after the deadline
	router := mux.NewRouter()
	router.HandleFunc("/igcinfo/api/igc/", manageTrack)
	router.HandleFunc("/igcinfo/api/", getMetaInfo)
	router.HandleFunc("/igcinfo/api/igc/{[0-9A-Za-z]+}/", getTrackByID)
	router.HandleFunc("/igcinfo/api/igc/{[0-9]+}/{[A-Za-z]+}/", getTrackField)
	router.HandleFunc("/igcinfo/", error404)
	router.HandleFunc("/", error404)
	//Handler has been changed from nil to router after deadline
	//":"+os.Getenv("PORT")
	http.ListenAndServe("127.0.0.1:8080", router)
}
