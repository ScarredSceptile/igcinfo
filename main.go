package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
	"strings"
	"net/http"
	"github.com/marni/goigc"
)

type URL struct {
	url string `json:"url"`
}

type Track struct {
	id string `json:"id"`
	Track igc.Track
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

	uptime := calcDuration(int(time.Now().Unix()) - timeStart)

	Response := "{"
	Response += "\"uptime\": \"" + uptime + "\","
	Response += "\"info\": \"Service for IGC tracks.\","
	Response += "\"version\": \"v1\""
	Response += "}"

	fmt.Fprintln(w, Response)
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

		var trackUrl URL
		decode := json.NewDecoder(r.Body)
		decode.Decode(&trackUrl)

		var track igc.Track
		var err error
		track, err = igc.ParseLocation(trackUrl.url)
		if err != nil {
			error400(w)
			return
		}

		//Check if track is new
		if !inList(track.UniqueID) {
			urlAmount += 1
			urlList[track.UniqueID] = Track{"igc" + strconv.Itoa(urlAmount), track}
			response := "{"
			response += "\"id:\" \"" + urlList[track.UniqueID].id + "\""
			response += "}"
			fmt.Fprintln(w,response)
		} else {
			error400(w)
			return
		}

	} else if r.Method == "GET" { //returns all the tracks ids
		w.Header().Set("Content-Type", "application/json")

		x := 0 //variable as i can't be compared to int

		response := "[ "
		for i := range urlList {
			response += urlList[i].id
			if x < len(urlList) {
				response += ", "
			}
			x++
		}

		response += " ]"
		fmt.Fprintln(w, response)
	}
}

//Gives all the information about a certain track by given ID
func getTrackById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	parts := strings.Split(r.URL.Path, "/")

	track := parts[len(parts)-1]
	if track != "" {
		response := "{"
		response += "\"H_date\": " + "\"" + urlList[track].Track.Date.String() + "\","
		response += "\"pilot\": " + "\"" + urlList[track].Track.Pilot + "\","
		response += "\"glider\": " + "\"" + urlList[track].Track.GliderType + "\","
		response += "\"glider_id\": " + "\"" + urlList[track].Track.GliderID + "\","
		response += "\"track_length\": " + "\"" + strconv.FormatFloat(igc.Distance(urlList[track].Track.Task), 'E', -1, 64) + "\","
		response += "}"

		fmt.Fprint(w, response)
	} else {
		error404(w, r)
	}


}

//Given id and field, returns the field of given id
func getTrackField(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.URL.Path, "/")

		id := parts[len(parts)-2]
		field:= parts[len(parts)-1]

		if id != "" && field != "" {
			switch field {
			case "pilot": fmt.Fprint(w, urlList[id].Track.Pilot)
			case "glider": fmt.Fprint(w, urlList[id].Track.GliderType)
			case "glider_id": fmt.Fprint(w, urlList[id].Track.GliderID)
			case "track_length": fmt.Fprint(w, igc.Distance(urlList[id].Track.Task))
			case "H_date": fmt.Fprint(w, urlList[id].Track.Date.String())
			}
		} else {
			error404(w, r)
		}
}

//Runs the application
func main(){
	http.HandleFunc("/igcinfo/api/igc/", manageTrack)
	http.HandleFunc("/igcinfo/api/", getMetaInfo)
	http.HandleFunc("/igcinfo/api/igc/{[0-9A-Za-z]+}/", getTrackById)
	http.HandleFunc("/igcinfo/api/igc/{[0-9]+}/{[A-Za-z]+}/", getTrackField)
	http.HandleFunc("/igcinfo/", error404)
	http.HandleFunc("/", error404)
	http.ListenAndServe(":" + os.Getenv("PORT"), nil)
}
