// Lanzs Client for recycleLah. Handle all bins details and feedback.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
)

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}

var tpl *template.Template

type RecycleBinDetails struct {
	ID              int     `json:"id"`              // auto increm and primary ID.
	BinID           string  `json:"binid"`           // need to assign "A00001" A for HDB recycling bin
	BinType         string  `json:"bintype"`         // A: Common Bins, E : E waste, C: Recycling center, M: Mix Bins , W: Workplace Bins
	BinLocationLat  float32 `json:"binlocationlat"`  // HC: 311.364587
	BinLocationLong float32 `json:"binlocationlong"` // HC: 1.364587
	BinAddress      string  `json:"locdescription"`  // Postcode 123456
	Postcode        string  `json:"postcode"`        // string but need to conv to int.
	UserID          string  `json:"userid"`          // from main site HC: Lanzshot
	FBoptions       string  `json:"fboption"`        // "Bin Fullness Status"
	ColorCode       string  `json:"colorcode"`       // "Yellow Half Full"
	Remarks         string  `json:"remarks"`         // "Please clear asap."
	BinStatusUpdate string  `json:"binstatusupdate"` // Completed / Rejected / Submitted
}

// var binFeedbacks recycleBinDetails

// WEB Server port and url.
const baseURL = "http://localhost:5000/api/v1/recyclebindetails"

func Index(res http.ResponseWriter, req *http.Request) {
	tpl.ExecuteTemplate(res, "index.gohtml", nil)
}

// Get and send user feedback.
func recycleBinFB(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {

		binFeedbacks := RecycleBinDetails{
			ID:              0,
			BinID:           "NIL",
			BinType:         "NIL",
			BinLocationLat:  0,
			BinLocationLong: 0,
			BinAddress:      req.FormValue("Binaddress"),
			Postcode:        req.FormValue("postcode"),
			UserID:          "Lanzs", //To be puck in.
			FBoptions:       req.FormValue("FBoptions"),
			ColorCode:       req.FormValue("binfull"),
			Remarks:         req.FormValue("remarks"),
			BinStatusUpdate: "Submitted",
		}
		fmt.Println("User binFeedbacks : ", binFeedbacks)

		jsonString, err := json.Marshal(binFeedbacks)
		if err != nil {
			fmt.Println("Json Mashal error :", err)
		}
		// os.Stdout.Write(jsonString)

		// apiCode := binFeedbacks.BinID
		fmt.Println("Sending User FB Via POST")
		response, err := http.Post(baseURL+"/feedback", "application/json", bytes.NewBuffer(jsonString)) //POST to create course.
		// response, err := http.Post(baseURL+"/"+apiCode, "application/json", bytes.NewBuffer(jsonString)) //POST to create course.

		if err != nil {
			fmt.Printf("The HTTP POST request failed with error %s\n", err)
		} else {
			defer response.Body.Close()
			data, _ := ioutil.ReadAll(response.Body)
			// fmt.Println("Add POST Course data:", data)
			// fmt.Println("Status Code : ", response.StatusCode)
			fmt.Println("Bin Details added :\n", string(data))
		}
	}
	tpl.ExecuteTemplate(res, "recycleBinOptions.gohtml", nil)
}

// Get user FB from DB.
func queryFB(res http.ResponseWriter, req *http.Request) {
	fmt.Println("Show user past Feedback with status.")
	var feedBacks []RecycleBinDetails
	userID := "Lanzs"

	fmt.Println(baseURL + "/feedback/" + userID)
	response, err := http.Get(baseURL + "/feedback/" + userID)

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		fmt.Println("FB Query Status Code : ", response.StatusCode) //200 OK.
	}
	defer response.Body.Close()

	JsonByteData, _ := ioutil.ReadAll(response.Body)
	fmt.Println("User FB Recieved from Server in Json:", JsonByteData)

	// convert JSON to object
	json.Unmarshal(JsonByteData, &feedBacks)
	fmt.Println("User FB  Details Recieved from Server in String:", feedBacks)

	tpl.ExecuteTemplate(res, "showUserFB.gohtml", feedBacks)

}

// show only recyclebins.
func showRecycleBins(res http.ResponseWriter, req *http.Request) {

	response, err := http.Get(baseURL)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		fmt.Println("Show all Bins Status Code : ", response.StatusCode) //200 OK.
	}
	defer response.Body.Close()

	JsonByteData, _ := ioutil.ReadAll(response.Body)
	fmt.Println("ALL Bin Details Recieved from Server in Json:", JsonByteData)
	allBinsDataJson := make([]RecycleBinDetails, 0)

	// convert JSON to object
	json.Unmarshal(JsonByteData, &allBinsDataJson)
	fmt.Println("ALL Bin Details Recieved from Server in String:", allBinsDataJson)

	tpl.ExecuteTemplate(res, "showRecycleBinsDetails.gohtml", allBinsDataJson)
}

func serverOption() http.Handler {
	router := mux.NewRouter()
	// router.PathPrefix("/RecycleLah/")

	// recycleIndex
	router.HandleFunc("/", Index)
	// recyclebins user feedback page
	router.HandleFunc("/recyclebinsFB", recycleBinFB)
	// Get user query pass feedback
	router.HandleFunc("/queryFB", queryFB)
	// show ALL recyclebins and FeedBack.
	router.HandleFunc("/showrecyclebins", showRecycleBins)
	return router
}

func main() {
	routerOption := serverOption()

	// handle to file serve for images dir "/images/"
	// http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("images"))))

	fmt.Println("Start client at Port : 5221")
	http.ListenAndServe(":5221", routerOption)
}
