package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func sendResponse(w http.ResponseWriter, statusCode int, payload interface{}, htmlfilename string, contentype string) {
	response, _ := json.Marshal(payload)
	if htmlfilename != "" {
		response = nil
		tmpl, err := template.ParseFiles(htmlfilename)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, payload)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	w.Header().Set("Content-type", contentype)
	w.WriteHeader(statusCode)
	if response != nil {
		w.Write(response)

	}
}

// // scrapQuran function to fetch and parse the Quran verses
// func scrapQuran(surahNumber int) (string, error) {
// 	url := fmt.Sprintf("https://api.globalquran.com/surah/%d/quran-uthmani-hafs", surahNumber)

// 	// Make the HTTP GET request
// 	resp, err := http.Get(url)
// 	if err != nil {
// 		return "", fmt.Errorf("error making HTTP request: %w", err)
// 	}
// 	defer resp.Body.Close()

// 	// Read the response body
// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		return "", fmt.Errorf("error reading response body: %w", err)
// 	}

// 	// Parse the JSON response
// 	var quranResponse QuranResponse
// 	err = json.Unmarshal(body, &quranResponse)
// 	if err != nil {
// 		return "", fmt.Errorf("error parsing JSON: %w", err)
// 	}

// 	// Combine the verses into a single string
// 	combinedVerses := ""
// 	for _, verse := range quranResponse.Quran.UthmaniHafs {
// 		combinedVerses += verse.Verse + " "
// 	}

// 	return combinedVerses, nil
// }

func getQuran(w http.ResponseWriter, req *http.Request) {
	// vars := mux.Vars(req)
	// SurahID := vars["id"]
	// url := "https://api.globalquran.com/surah/1/quran-uthmani-hafs"
	// scrapdata, err := ScrapeQuranData(url)
	// if err != nil {
	// 	return
	// }
	// surahNumber := 1
	// combinedVerses, err := scrapQuran(surahNumber)
	// fmt.Println(combinedVerses)

	// Read the contents of the text file
	surahname, errs := ioutil.ReadFile("surahname.txt")
	if errs != nil {
		// Log the error but don't terminate the server
		// log.Printf("Failed to read file: %v", err)
		// fmt.Fprintf(w, "Error reading file")
		return
	}

	audio, errs := ioutil.ReadFile("audio.txt")
	if errs != nil {
		// Log the error but don't terminate the server
		// log.Printf("Failed to read file: %v", err)
		// fmt.Fprintf(w, "Error reading file")
		return
	}
	straudio := strings.ReplaceAll(string(audio), `"`, ``)

	// Read the contents of the text file
	finale, err := ioutil.ReadFile("finale.txt")
	if err != nil {
		// Log the error but don't terminate the server
		// log.Printf("Failed to read file: %v", err)
		// fmt.Fprintf(w, "Error reading file")
		return
	}
	finalestr := string(finale)
	// finalestr = strings.ReplaceAll(finalestr, "]", "]")
	// QuranfinaleData := replaceAtIndex(finalestr, ' ', 0)
	QuranfinaleData := finalestr
	// var result map[string][]string
	// err = json.Unmarshal([]byte(finale), &result)
	// if err != nil {
	// 	fmt.Println("Error:", err)
	// 	return
	// }
	quran := NewQG("http://api.globalquran.com/surah/", "", map[string]string{"en": "quran-simple", "ar": "quran-simple"}, 10)
	surah := 1       // example surah number
	ayah := 1        // example ayah number
	language := "en" // example language code

	ayahInfos, err := quran.getAyah(surah, ayah, language)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(ayahInfos)
	// for _, ayahInfo := range ayahInfos {
	// 	fmt.Printf("ID: %v, \n", ayahInfo)
	// }
	// Prepare the payload as a map
	// fmt.Println(scrapdata)
	Data := map[string]string{
		"SurahName":  string(surahname),
		"FinaleData": QuranfinaleData,
		"audio":      string(straudio),
		// "scrapdata":  result["verse"],
	}

	// Send the response with the template
	sendResponse(w, http.StatusOK, Data, "templates/GetSurah.html", "text/html; charset=UTF-8")
	// sendResponse(w, http.StatusOK, QuranDataStr, "templates/index.html", "text/html; charset=UTF-8")

}

func Search(w http.ResponseWriter, r *http.Request) {
	// Get the id from the url parameter

	query := r.URL.Query().Get("q")

	url := fmt.Sprintf("https://api.quran.com/api/v4/search?q=%v&size=20&page=1&language=en", query)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	strbordy := string(body)
	if query == "" {
		strbordy = "Try To Search Up There 👆"
		// http.Error(w, "Searh query parameter (q) is required", http.StatusBadRequest)
		// return
	}
	// fmt.Println(string(body))
	Data := map[string]string{
		"SearchData": strbordy,
		// "scrapdata":  result["verse"],
	}

	// Send the response with the template
	sendResponse(w, http.StatusOK, Data, "templates/Search.html", "text/html; charset=UTF-8")

}

func Home(w http.ResponseWriter, r *http.Request) {
	// Send the response with the template
	sendResponse(w, http.StatusOK, nil, "templates/index.html", "text/html; charset=UTF-8")

}

func main() {

	http.HandleFunc("/", getQuran)
	http.HandleFunc("/home", Home)
	http.HandleFunc("/search/", Search)
	log.Println("Server starting on :9090")
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
