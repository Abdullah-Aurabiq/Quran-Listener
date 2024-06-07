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
func replaceAtIndex(input string, replacement rune, index int) string {
	out := []rune(input)
	out[index] = replacement
	return string(out)
}

func hello(w http.ResponseWriter, req *http.Request) {
	// Read the contents of the text file
	surahname, errs := ioutil.ReadFile("surahname.txt")
	if errs != nil {
		// Log the error but don't terminate the server
		// log.Printf("Failed to read file: %v", err)
		// fmt.Fprintf(w, "Error reading file")
		return
	}

	// Read the contents of the text file
	data, err := ioutil.ReadFile("output.txt")
	if err != nil {
		// Log the error but don't terminate the server
		// log.Printf("Failed to read file: %v", err)
		// fmt.Fprintf(w, "Error reading file")
		return
	}
	QuranDataStr := string(data)
	QuranDataStr = strings.ReplaceAll(QuranDataStr, "]", "]\n")
	QuranData := replaceAtIndex(QuranDataStr, ' ', 0)
	// Print the contents
	fmt.Println(QuranData)

	// Read the contents of the text file
	translation, erred := ioutil.ReadFile("translation.txt")
	if erred != nil {
		// Log the error but don't terminate the server
		// log.Printf("Failed to read file: %v", err)
		// fmt.Fprintf(w, "Error reading file")
		return
	}
	translationstr := string(translation)
	translationstr = strings.ReplaceAll(translationstr, "]", "]\n\n")
	QuranTranslationData := replaceAtIndex(translationstr, ' ', 0)
	// Print the contents
	fmt.Println(QuranTranslationData)

	// Prepare the payload as a map
	payload := map[string]string{
		"SurahName":   string(surahname),
		"QuranVerse":  QuranData,
		"Translation": QuranTranslationData,
	}

	// Send the response with the template
	sendResponse(w, http.StatusOK, payload, "templates/index.html", "text/html; charset=UTF-8")
	// sendResponse(w, http.StatusOK, QuranDataStr, "templates/index.html", "text/html; charset=UTF-8")

}

func main() {
	// Ensure the file exists and create it if it doesn't
	// data := []byte("Hello from Python!\n")
	// err := ioutil.WriteFile("output.txt", data, 0644)
	// if err != nil {
	// 	log.Fatalf("Failed to create file: %v", err)
	// }

	http.HandleFunc("/", hello)
	log.Println("Server starting on :8091")
	err := http.ListenAndServe(":1485", nil)
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
