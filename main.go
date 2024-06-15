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

// scrapQuran function to fetch and parse the Quran verses
func scrapQuran(surahNumber int) (string, error) {
	url := fmt.Sprintf("https://api.globalquran.com/surah/%d/quran-uthmani-hafs", surahNumber)

	// Make the HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("error making HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	// Parse the JSON response
	var quranResponse QuranResponse
	err = json.Unmarshal(body, &quranResponse)
	if err != nil {
		return "", fmt.Errorf("error parsing JSON: %w", err)
	}

	// Combine the verses into a single string
	combinedVerses := ""
	for _, verse := range quranResponse.Quran.UthmaniHafs {
		combinedVerses += verse.Verse + " "
	}

	return combinedVerses, nil
}

func getQuran(w http.ResponseWriter, req *http.Request) {
	// vars := mux.Vars(req)
	// SurahID := vars["id"]
	// url := "https://api.globalquran.com/surah/1/quran-uthmani-hafs"
	// scrapdata, err := ScrapeQuranData(url)
	// if err != nil {
	// 	return
	// }
	surahNumber := 1
	combinedVerses, err := scrapQuran(surahNumber)
	fmt.Println(combinedVerses)

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

	// Prepare the payload as a map
	// fmt.Println(scrapdata)
	Data := map[string]string{
		"SurahName":  string(surahname),
		"FinaleData": QuranfinaleData,
		"audio":      string(straudio),
		"scrapdata":  combinedVerses,
	}

	// Send the response with the template
	sendResponse(w, http.StatusOK, Data, "templates/GetSurah.html", "text/html; charset=UTF-8")
	// sendResponse(w, http.StatusOK, QuranDataStr, "templates/index.html", "text/html; charset=UTF-8")

}

func main() {
	http.HandleFunc("/", getQuran)
	log.Println("Server starting on :9090")
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
