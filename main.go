package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

// App represents the application
type App struct {
	Router *mux.Router
}

// New creates a new instance of App
func (app *App) Initialize() error {
	app.Router = mux.NewRouter().StrictSlash(true)
	app.handleRoutes()
	return nil
}

// SetDB sets the database connection
func (app *App) Run(address string) {
	log.Fatal(http.ListenAndServe(address, app.Router))
}

const API = "https://www.mp3quran.net/api/"

// StripTashkeel removes Arabic diacritics from the input text.
func StripTashkeel(text string) string {
	diacritics := []rune{
		'\u064b', // FATHATAN
		'\u064c', // DAMMATAN
		'\u064d', // KASRATAN
		'\u064e', // FATHA
		'\u064f', // DAMMA
		'\u0650', // KASRA
		'\u0651', // SHADDA
		'\u0652', // SUKUN
	}

	var result strings.Builder
	for _, char := range text {
		if !contains(diacritics, char) {
			result.WriteRune(char)
		}
	}
	return result.String()
}

// contains checks if a slice contains a specific rune.
func contains(s []rune, r rune) bool {
	for _, a := range s {
		if a == r {
			return true
		}
	}
	return false
}

// SuraName represents the JSON structure returned by the API.
type SuraName struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// GetSuraName fetches and returns the name of the Sura for the given number.
func GetSuraName(suraNumber int) (string, error) {
	if suraNumber <= 0 || suraNumber > 114 {
		return "", fmt.Errorf("invalid sura number: sura not found '%d'", suraNumber)
	}

	response, err := http.Get(API + "_arabic_sura.php")
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	var suras struct {
		SurasName []SuraName `json:"Suras_Name"`
	}
	err = json.Unmarshal(body, &suras)
	if err != nil {
		return "", err
	}

	for _, sura := range suras.SurasName {
		if sura.ID == fmt.Sprintf("%d", suraNumber) {
			return StripTashkeel(strings.TrimSpace(sura.Name)), nil
		}
	}

	return "", fmt.Errorf("sura not found for number '%d'", suraNumber)
}

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

func getQuran(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	SurahID, e := strconv.Atoi(vars["id"])
	if e != nil {
		sendResponse(w, http.StatusBadRequest, "Invalid surah number", "", "application/json")
		return
	}

	surahname, err := GetSuraName(SurahID)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Sura Name:", surahname)
	}

	quranen := NewQG("http://api.globalquran.com/surah/", "", map[string]string{"ar": "en.qaribullah", "en": "en.qaribullah"}, 10)
	surah := SurahID // example surah number
	ayah := 1        // example ayah number
	language := "en" // example language code
	// "quran-uthmani-hafs"

	ayahInfosen, err := quranen.getAyah(surah, ayah, language)
	if err != nil {
		fmt.Println(err)
		return
	}
	ayahInfoseng, _ := json.Marshal(ayahInfosen)

	quranar := NewQG("http://api.globalquran.com/surah/", "", map[string]string{"ar": "quran-uthmani-hafs", "en": "quran-uthmani-hafs"}, 10)
	language = "ar" // example language code
	ayahInfosar, err := quranar.getAyah(surah, ayah, language)
	if err != nil {
		fmt.Println(err)
		return
	}
	ayahInfosara, _ := json.Marshal(ayahInfosar)
	// keys := reflect.ValueOf(ayahInfosen["quran"]["en.qaribullah"]).MapKeys()
	English_try := string(ayahInfoseng)
	English_try = strings.ReplaceAll(English_try, "map[quran:map[en.qaribullah:map[", "s")
	English_try = strings.ReplaceAll(English_try, "]]]]", "")
	// fmt.Println(English_try)
	Arabic_try := string(ayahInfosara)
	Arabic_try = strings.ReplaceAll(Arabic_try, "map[quran:map[en.qaribullah:map[", "s")
	Arabic_try = strings.ReplaceAll(Arabic_try, "]]]]", "")
	// fmt.Println(Arabic_try)

	var Englishdata map[string]map[string]map[int]Verse
	err = json.Unmarshal([]byte(English_try), &Englishdata)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Create a slice to store the verses
	var English_verses []Verse
	for _, innerMap := range Englishdata["quran"]["en.qaribullah"] {
		English_verses = append(English_verses, innerMap)
	}

	// Sort the English_verses by Ayah number

	sort.Slice(English_verses, func(i, j int) bool {
		return English_verses[i].Ayah < English_verses[j].Ayah
	})

	// Display the sorted English_verses
	var English_output string
	for _, verse := range English_verses {
		English_output += fmt.Sprintf("%s[%d]A8ea8", verse.Verse, verse.Ayah)
		// fmt.Printf("English: Verse %d: %s\n", verse.Ayah, verse.Verse)
	}

	var arabic_data map[string]map[string]map[int]Verse
	err = json.Unmarshal([]byte(Arabic_try), &arabic_data)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Create a slice to store the verses
	var arabic_verses []Verse
	for _, innerMap := range arabic_data["quran"]["quran-uthmani-hafs"] {
		arabic_verses = append(arabic_verses, innerMap)
	}

	// Sort the arabic_verses by Ayah number
	sort.Slice(arabic_verses, func(i, j int) bool {
		return arabic_verses[i].Ayah < arabic_verses[j].Ayah
	})

	// Display the sorted arabic_verses
	var arabic_output string
	for _, verse := range arabic_verses {
		arabic_output += fmt.Sprintf("%s[%d]A8ea8", verse.Verse, verse.Ayah)
		// fmt.Printf("Arabic: Verse %d: %s\n", verse.Ayah, verse.Verse)
	}
	ar := strings.Split(arabic_output, "A8ea8")
	en := strings.Split(English_output, "A8ea8")

	// fmt.Println("English: ", en)
	// fmt.Println("Arabic: ", ar)
	a := make(map[string][]string)
	a["ar"] = ar
	a["en"] = en
	// fmt.Println(a)
	// str := fmt.Sprintf("%v", a)
	as, _ := json.Marshal(a)
	ar_en := string(as)
	ar_en = strings.ReplaceAll(ar_en, "map[", "")
	ar_en = strings.ReplaceAll(ar_en, "]]", "")
	// fmt.Println(ar_en)
	audioUrl, _ := getAudio(1)
	var SurahIDs string
	if SurahID < 10 {
		SurahIDs = "00" + strconv.Itoa(SurahID)
	} else if SurahID < 100 {
		SurahIDs = "0" + strconv.Itoa(SurahID)
	} else if SurahID >= 100 {
		SurahIDs = strconv.Itoa(SurahID)
	}
	fmt.Println(SurahIDs)
	Data := map[string]string{
		"SurahName":  string(surahname),
		"FinaleData": ar_en,
		"audio":      audioUrl,
		"id":         SurahIDs,
		// "englishQuran": English_output,
		// "arabicQuran":  arabic_output,
	}

	// Send the response with the template
	sendResponse(w, http.StatusOK, Data, "templates/GetSurah.html", "text/html; charset=UTF-8")
	// sendResponse(w, http.StatusOK, QuranDataStr, "templates/index.html", "text/html; charset=UTF-8")

}

func getAudio(surah int) (string, error) {
	url := fmt.Sprintf("https://api.quran.com/api/v4/chapter_recitations/10/%d", surah)
	// url := fmt.Sprintf("https://api.quran.com/api/v4/chapter_recitations/2/%d", surah)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var dictres map[string]interface{}
	err = json.Unmarshal(body, &dictres)
	if err != nil {
		return "", err
	}

	audioURL, ok := dictres["audio_file"].(map[string]interface{})["audio_url"].(string)
	if !ok {
		return "", fmt.Errorf("audio url not found in response")
	}

	return audioURL, nil
}

func Search(w http.ResponseWriter, r *http.Request) {
	// Get the id from the url parameter

	query := r.URL.Query().Get("q")

	url := fmt.Sprintf("https://api.quran.com/api/v4/search?q=%v&size=200&page=1&language=en", query)
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

func Hadith(w http.ResponseWriter, r *http.Request) {
	url := "https://hadithapi.com/api/sahih-bukhari/chapters?apiKey=$2y$10$u6K80SDvlCph1KgbQOOq0uaC68QRd1JwsESIYRZwOvc9ARow1TZq"
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

	payload := map[string]string{
		"body": strbordy,
		// "scrapdata":  result["verse"],
	}

	sendResponse(w, http.StatusOK, payload, "templates/hadith.html", "text/html; charset=UTF-8")

}
func (app *App) handleRoutes() {
	app.Router.HandleFunc("/quran/{id}", getQuran).Methods("GET")
	app.Router.HandleFunc("/home/", Home)
	app.Router.HandleFunc("/search/", Search)
	app.Router.HandleFunc("/hadith/", Hadith)

}

func main() {
	app := App{}
	log.Println("Server starting on :1481")
	app.Initialize()
	app.Run("localhost:1481")
}
