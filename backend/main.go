package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type PostBody struct {
	Message string `json:message`
}

type ID struct {
	ID string `json:ID`
}

type Merchandise struct {
	ID        string  `json:ID`
	Name      string  `json:Name`
	Category  string  `json:Category`
	Condition string  `json:Condition`
	Price     float64 `json:Price`
}

type Database struct {
	Item []Merchandise
}

var (
	dataBase Database
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
}

func main() {
	r := mux.NewRouter()
	// r.HandleFunc("/", frontpage)
	r.HandleFunc("/AddItem", addItem).Methods("POST")
	r.HandleFunc("/EditItem", editItem).Methods("POST")
	r.HandleFunc("/RemoveItem", removeItem).Methods("POST")
	r.HandleFunc("/ListItem", listItem).Methods("GET")
	r.HandleFunc("/GetItemByID", getItemByID).Methods("GET")

	log.Info().Msg("Server Starting...")
	http.ListenAndServe(":8080", r)
}

func GetHostName() string {
	name, err := os.Hostname()
	if err != nil {
		log.Error().Err(err).Msg("Failed to retrieve MachineHost Name.")
	}
	return name
}

func addItem(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error().Err(err).Str("Machine", GetHostName()).Str("Route", r.RequestURI).Msg("Request Error.")
	}
	var data Merchandise
	json.Unmarshal([]byte(body), &data)
	dataBase.Item = append(dataBase.Item, data)
	log.Info().Str("Machine", GetHostName()).Str("Route", r.RequestURI).Msg("Item " + data.ID + " added to database.")

	var reply PostBody
	reply.Message = "Item Added..."
	json, _ := json.Marshal(reply)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func editItem(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error().Err(err).Str("Machine", GetHostName()).Str("Route", r.RequestURI).Msg("Request Error.")
	}
	var data Merchandise
	json.Unmarshal([]byte(body), &data)

	var updBlock int

	for i, v := range dataBase.Item {
		if data.ID != "" {
			if data.ID == v.ID {
				updBlock = i
				break
			}
		}
	}

	var replacement Merchandise

	replacement.ID = dataBase.Item[updBlock].ID

	if data.Name != "" {
		replacement.Name = data.Name
	} else {
		replacement.Name = dataBase.Item[updBlock].Name
	}
	if data.Category != "" {
		replacement.Category = data.Category
	} else {
		replacement.Category = dataBase.Item[updBlock].Category
	}
	if data.Condition != "" {
		replacement.Condition = data.Condition
	} else {
		replacement.Condition = dataBase.Item[updBlock].Condition
	}
	if data.Price != 0 {
		replacement.Price = data.Price
	} else {
		replacement.Price = dataBase.Item[updBlock].Price
	}
	dataBase.Item[updBlock] = replacement
	log.Info().Str("Machine", GetHostName()).Str("Route", r.RequestURI).Msg("Item " + dataBase.Item[updBlock].ID + " was edited.")

	var reply PostBody
	reply.Message = "Item Updated..."
	json, _ := json.Marshal(reply)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func removeItem(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error().Err(err).Str("Machine", GetHostName()).Str("Route", r.RequestURI).Msg("Request Error.")
	}
	var itemID ID
	json.Unmarshal([]byte(body), &itemID)

	var placeholder Database

	for i, v := range dataBase.Item {
		if itemID.ID != "" {
			if itemID.ID != v.ID {
				placeholder.Item = append(placeholder.Item, dataBase.Item[i])
			}
		}
	}

	dataBase = placeholder

	log.Info().Str("Machine", GetHostName()).Str("Route", r.RequestURI).Msg("Item " + itemID.ID + " was deleted.")

	var reply PostBody
	reply.Message = "Item Removed..."
	json, _ := json.Marshal(reply)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func listItem(w http.ResponseWriter, r *http.Request) {
	json, _ := json.Marshal(dataBase.Item)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
	log.Info().Str("Machine", GetHostName()).Str("Route", r.RequestURI).Msg("Item list requested.")
}

func getItemByID(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error().Err(err).Str("Machine", GetHostName()).Str("Route", r.RequestURI).Msg("Request Error. " + r.RequestURI)
	}

	var data ID
	json.Unmarshal([]byte(body), &data)

	for _, v := range dataBase.Item {
		if v.ID == data.ID {
			json, _ := json.Marshal(v)
			w.Header().Set("Content-Type", "application/json")
			w.Write(json)
			log.Info().Str("Machine", GetHostName()).Str("Route", r.RequestURI).Msg("Item info requested.")
			break
		}
	}
}

// func frontpage(w http.ResponseWriter, r *http.Request) {
// 	if r.URL.Path != "/" {
// 		http.Error(w, "404 not found.", http.StatusNotFound)
// 		return
// 	}

// 	switch r.Method {
// 	case "GET":
// 		http.ServeFile(w, r, "./frontend/frontpage.html")
// 	case "POST":
// 		body, err := ioutil.ReadAll(r.Body)
// 		if err != nil {
// 			log.Println(err)
// 		}
// 		var data PostBody
// 		json.Unmarshal([]byte(body), &data)
// 		data.Message = reverseBody(data.Message)
// 		json, _ := json.Marshal(data)
// 		w.Write(json)
// 	default:
// 		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
// 	}

// }

// func reverseBody(body string) (rString string) {
// 	for _, v := range body {
// 		rString = string(v) + rString
// 	}
// 	return
// }
