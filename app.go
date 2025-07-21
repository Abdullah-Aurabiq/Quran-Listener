package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"golang.ngrok.com/ngrok/v2"
)

// App represents the application

type App struct {
	Router *mux.Router
	// DB     *sql.DB
}

var store = sessions.NewCookieStore([]byte("`giytIBFi0<.-3U,y$<!fdQFx$?@)}"))

// SurahCardData represents the data structure for a Surah card
type SurahCardData struct {
	ID             int    `json:"id"`
	EnglishName    string `json:"englishName"`
	ArabicName     string `json:"arabicName"`
	EnglishMeaning string `json:"englishMeaning"`
	TotalVerses    int    `json:"totalVerses"`
	StartingVerses string `json:"startingVerses"`
}

// New creates a new instance of App
func (app *App) Initialize() error {
	// connectionString := fmt.Sprintf("%v:%v@tcp(127.0.0.1:3306)/%v", Dbuser, Dbpassword, Dbname)
	// var err error
	// app.DB, err = sql.Open("mysql", connectionString)
	// if err != nil {
	// 	return err
	// }

	app.Router = mux.NewRouter().StrictSlash(true)
	app.handleRoutes()
	return nil
}

// SetDB sets the database connection
func (app *App) Run(address string) {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file")
	}

	// ✅ Run ngrok in background
	go func() {
		time.Sleep(2 * time.Second)
		if err := run(context.Background()); err != nil {
			log.Fatalf("Failed to start ngrok agent: %v", err)
		}
	}()

	// ✅ Start actual server
	handler := c.Handler(app.Router)
	log.Println("Starting server on", address)
	log.Fatal(http.ListenAndServe(address, handler))
}

// const trafficPolicy = `
// on_http_request:
//   - actions:
//   - type: oauth
//     config:
//     provider: google
//
// `
const address = "localhost:1481"

func run(ctx context.Context) error {
	agent, err := ngrok.NewAgent(ngrok.WithAuthtoken(os.Getenv("NGROK_AUTHTOKEN")))
	if err != nil {
		return err
	}
	time.Sleep(2 * time.Second) // Wait for the agent to start
	ln, err := agent.Forward(ctx,
		ngrok.WithUpstream(address),
		ngrok.WithURL("electric-mistakenly-rat.ngrok-free.app"),
		// ngrok.WithTrafficPolicy(trafficPolicy),
	)

	if err != nil {
		fmt.Println("Error", err)
		return err
	}

	fmt.Println("Endpoint online: forwarding from", ln.URL(), "to", address)

	// Explicitly stop forwarding; otherwise it runs indefinitely
	<-ln.Done()
	return nil
}

const API = "https://www.mp3quran.net/api/_arabic_sura.php"

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

	client := &http.Client{}
	req, err := http.NewRequest("GET", API, nil)
	if err != nil {
		return "", err
	}

	// Add headers to the request
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Referer", "https://www.mp3quran.net/")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cache-Control", "no-cache")

	response, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch sura name: %s", response.Status)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	var suras struct {
		SurasName []SuraName `json:"Suras_Name"`
	}
	err = json.Unmarshal(body, &suras)
	if err != nil {
		return "", fmt.Errorf("{error: 'failed to unmarshal JSON'} %v", err)
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
		tmpl, err := template.ParseFiles(htmlfilename, "templates/navbar.html")
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
		// fmt.Printf("English: Verse %E: %s\n", verse.Ayah, verse.Verse)
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
		// fmt.Printf("Arabic: Verse %E: %s\n", verse.Ayah, verse.Verse)
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
	audioUrl, _ := getAudio(SurahID)
	var SurahIDs string
	if SurahID < 10 {
		SurahIDs = "00" + strconv.Itoa(SurahID)
	} else if SurahID < 100 {
		SurahIDs = "0" + strconv.Itoa(SurahID)
	} else if SurahID >= 100 {
		SurahIDs = strconv.Itoa(SurahID)
	}
	fmt.Println(SurahIDs)
	session, _ := store.Get(req, "session")
	fmt.Println(session.Values)
	qu := map[string]string{

		"SurahName":  string(surahname),
		"FinaleData": ar_en,
		"audio":      audioUrl,
		"id":         SurahIDs,
	}
	Data := map[string]interface{}{
		"user":  session.Values,
		"quran": qu,
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

func sendError(w http.ResponseWriter, statusCode int, err string) {
	errormessage := map[string]string{"error": err}
	sendResponse(w, statusCode, errormessage, "", "text/html; charset=UTF-8")
}

func Intro(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	auth, ok := session.Values["authenticated"].(bool)
	if ok || auth {
		http.Redirect(w, r, "/quran/", http.StatusFound)
		return
	} else {
		sendResponse(w, http.StatusOK, nil, "templates/intro.html", "text/html; charset=UTF-8")
	}
}
func contact_us(w http.ResponseWriter, r *http.Request) {
	// session, _ := store.Get(r, "session")
	// auth, ok := session.Values["authenticated"].(bool)
	// if !ok || !auth {
	// 	http.Redirect(w, r, "/login", http.StatusFound)
	// 	return
	// }
	data := map[string]interface{}{
		"user": map[string]string{"username": "sd"},
	}
	sendResponse(w, http.StatusOK, data, "templates/contact-us.html", "text/html; charset=UTF-8")
}
func about_us(w http.ResponseWriter, r *http.Request) {
	sendResponse(w, http.StatusOK, nil, "templates/about-us.html", "text/html; charset=UTF-8")
}

// func (app *App) SignupAPI(w http.ResponseWriter, r *http.Request) {
// 	var st Users
// 	username := r.FormValue("username")
// 	gender := r.FormValue("gender")
// 	email := r.FormValue("email")
// 	password := r.FormValue("password")
// 	confirmPassword := r.FormValue("confirmpassword")
// 	fmt.Println(username, email, password, confirmPassword)
// 	st.UserName = username
// 	st.Email = email
// 	st.Password = password
// 	st.Gender = gender
// 	if password != confirmPassword {
// 		sendError(w, http.StatusBadRequest, "Passwords do not match")
// 		return
// 	}

// 	if username != "" && email != "" && password != "" {
// 		// Create User account
// 		err := st.CreateUserAccount(app.DB)
// 		if err != nil {
// 			sendError(w, http.StatusInternalServerError, err.Error())
// 			return
// 		}
// 		// Create a session for the user
// 		session, _ := store.Get(r, "session")
// 		session.Values["email"] = st.Email
// 		session.Values["username"] = st.UserName
// 		session.Values["authenticated"] = true
// 		err = session.Save(r, w)
// 		if err != nil {
// 			sendError(w, http.StatusInternalServerError, "Failed to save session")
// 			return
// 		}

// 		// Prepare the response data
// 		responseData := map[string]interface{}{
// 			"email":         st.Email,
// 			"username":      st.UserName,
// 			"authenticated": true,
// 		}

// 		// Convert response data to JSON
// 		jsonData, err := json.Marshal(responseData)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}

// 		// Set the content type to application/json and write the JSON response
// 		w.Header().Set("Content-Type", "application/json")
// 		w.WriteHeader(http.StatusOK)
// 		w.Write(jsonData)
// 		return
// 	}

// 	sendError(w, http.StatusBadRequest, "Invalid input")
// }

// func (app *App) LoginAPI(w http.ResponseWriter, r *http.Request) {
// 	// Parse form values
// 	email := r.FormValue("email")
// 	password := r.FormValue("password")

// 	// Validate input
// 	if email == "" || password == "" {
// 		http.Error(w, "Email and password are required", http.StatusBadRequest)
// 		return
// 	}

// 	var st Users
// 	st.Email = email
// 	st.Password = password
// 	err := st.UserLogin(app.DB)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	// Create a session for the user
// 	session, _ := store.Get(r, "session")
// 	session.Values["email"] = st.Email
// 	session.Values["username"] = st.UserName
// 	session.Values["authenticated"] = true
// 	err = session.Save(r, w)
// 	if err != nil {
// 		http.Error(w, "Failed to save session", http.StatusInternalServerError)
// 		return
// 	}

// 	// Prepare the response data
// 	responseData := map[string]interface{}{
// 		"email":         st.Email,
// 		"username":      st.UserName,
// 		"authenticated": true,
// 	}

// 	// Convert response data to JSON
// 	jsonData, err := json.Marshal(responseData)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	// Set the content type to application/json and write the JSON response
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	w.Write(jsonData)
// }

func (app *App) Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	delete(session.Values, "authenticated")
	session.Save(r, w)
	http.Redirect(w, r, "/login", http.StatusFound)
}
func getQuranAPI(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	SurahID, e := strconv.Atoi(vars["id"])
	if e != nil {
		sendResponse(w, http.StatusBadRequest, "Invalid surah number", "", "application/json")
		return
	}
	// Get query parameters for Arabic version and translation
	// arabicVersion := req.URL.Query().Get("ar")
	// if arabicVersion == "" {
	// 	arabicVersion = "quran-uthmani-hafs" // default value
	// }
	translation := req.URL.Query().Get("translation")
	if translation == "" {
		translation = "en.sahih" // default value
	}
	// Fetch Surah name
	// surahname, err := GetSuraName(SurahID)
	// if err != nil {
	// 	fmt.Println("Error:", err)
	// } else {
	// 	fmt.Println("Sura Name:", surahname)
	// }

	// Fetch English version from API
	quranen := NewQG("http://api.globalquran.com/surah/", "", map[string]string{"ar": translation, "en": translation}, 10)
	language := "en"
	ayahInfosen, err := quranen.getAyah(SurahID, 1, language)
	if err != nil {
		fmt.Println(err)
		return
	}
	ayahInfoseng, _ := json.Marshal(ayahInfosen)

	// Read Arabic version from local JSON file
	// filePath := fmt.Sprintf("D:/F(DRIVE/Quran Listener/static/quranar/surahs/%03d.json", SurahID)
	// request https://electric-mistakenly-rat.ngrok-free.app/static/quranar/surahs/%03d.json for filePath
	filePath := fmt.Sprintf("http://electric-mistakenly-rat.ngrok-free.app/static/quranar/surahs/%03d.json", SurahID)

	fileData, err := ioutil.ReadFile(filePath)
	// fmt.Println(string(fileData))
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to read surah data", http.StatusInternalServerError)
		return
	}

	var arabicData map[string]interface{}
	err = json.Unmarshal(fileData, &arabicData)
	if err != nil {
		http.Error(w, "Failed to parse surah data", http.StatusInternalServerError)
		return
	}

	// Process English version
	English_try := string(ayahInfoseng)
	English_try = strings.ReplaceAll(English_try, "map[quran:map[en.qaribullah:map[", "s")
	English_try = strings.ReplaceAll(English_try, "]]]]", "")

	var Englishdata map[string]map[string]map[int]Verse
	err = json.Unmarshal([]byte(English_try), &Englishdata)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var English_verses []Verse
	for _, innerMap := range Englishdata["quran"][translation] {
		English_verses = append(English_verses, innerMap)
	}

	sort.Slice(English_verses, func(i, j int) bool {
		return English_verses[i].Ayah < English_verses[j].Ayah
	})

	// Process Arabic version
	arabicVerses := arabicData["quran"].(map[string]interface{})["quran-uthmani-hafs"].(map[string]interface{})
	var Arabic_verses []Verse
	for _, v := range arabicVerses {
		verse := v.(map[string]interface{})
		Arabic_verses = append(Arabic_verses, Verse{
			Ayah:  int(verse["ayah"].(float64)),
			Verse: verse["verse"].(string),
		})
	}

	sort.Slice(Arabic_verses, func(i, j int) bool {
		return Arabic_verses[i].Ayah < Arabic_verses[j].Ayah
	})

	// Combine Arabic and English verses
	var combinedVerses []map[string]string
	for i := 0; i < len(English_verses) && i < len(Arabic_verses); i++ {
		combinedVerses = append(combinedVerses, map[string]string{
			"ar": Arabic_verses[i].Verse,
			"en": English_verses[i].Verse,
		})
	}

	// Convert combined verses to JSON
	jsonData, err := json.Marshal(combinedVerses)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	audioUrl, _ := getAudio(SurahID)
	var SurahIDs string
	if SurahID < 10 {
		SurahIDs = "00" + strconv.Itoa(SurahID)
	} else if SurahID < 100 {
		SurahIDs = "0" + strconv.Itoa(SurahID)
	} else if SurahID >= 100 {
		SurahIDs = strconv.Itoa(SurahID)
	}

	qu := map[string]string{
		// "SurahName":  string(surahname),
		"FinaleData": string(jsonData),
		"audio":      audioUrl,
		"id":         SurahIDs,
	}
	Data := map[string]interface{}{
		"quran": qu,
	}
	s, _ := json.Marshal(Data)

	// Set the content type to application/json and write the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(s)
}

func convertSessionValues(values map[interface{}]interface{}) map[string]interface{} {
	converted := make(map[string]interface{})
	for k, v := range values {
		strKey, ok := k.(string)
		if !ok {
			continue
		}
		converted[strKey] = v
	}
	return converted
}
func SearchAPI(w http.ResponseWriter, r *http.Request) {
	// Get the query parameter from the URL
	query := r.URL.Query().Get("q")

	url := fmt.Sprintf("https://api.quran.com/api/v4/search?q=%v&size=200&page=1&language=en", query)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	strbody := string(body)
	if query == "" {
		strbody = "Try To Search Up There 👆"
	}

	Data := map[string]string{
		"SearchData": strbody,
	}

	// Convert Data map to JSON
	jsonData, err := json.Marshal(Data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the content type to application/json and write the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
func getSurahs(w http.ResponseWriter, req *http.Request) {
	var surahs []SurahCardData
	// var filteredSurahs []SurahCardData

	// Read the surahs.json file
	url := "https://electric-mistakenly-rat.ngrok-free.app/static/surahs.json"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		http.Error(w, "Failed to fetch surah data", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
		return
	}

	// Parse the JSON data
	err = json.Unmarshal(data, &surahs)
	if err != nil {
		fmt.Println("Error:", err)
		http.Error(w, "Failed to parse surah data", http.StatusInternalServerError)
		return
	}

	// Filter Surahs based on criteria (favorite, recently viewed, prayer times)
	// for _, surah := range surahs {
	// 	if isFavorite(surah.ID) || isRecentlyViewed(surah.ID) || isRelatedToPrayerTime(surah.ID) {
	// 		filteredSurahs = append(filteredSurahs, surah)
	// 	}
	// }

	// Convert to JSON and send response
	jsonResponse, err := json.Marshal(surahs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

// Helper functions to determine if a Surah meets the criteria
// func isFavorite(surahID int) bool {
// 	// Implement logic to check if the Surah is marked as favorite
// 	return false
// }

// func isRecentlyViewed(surahID int) bool {
// 	// Implement logic to check if the Surah is recently viewed
// 	return false
// }

// func isRelatedToPrayerTime(surahID int) bool {
// 	// Implement logic to check if the Surah is related to the current prayer time
// 	return false
// }

func (app *App) handleRoutes() {
	// app.Router.HandleFunc("/home/", Intro)
	app.Router.HandleFunc("/api/hadith/", Hadith)
	app.Router.HandleFunc("/contact-us/", contact_us)
	// app.Router.HandleFunc("/about-us/", about_us)
	// app.Router.HandleFunc("/api/register", app.SignupAPI)
	// app.Router.HandleFunc("/api/login", app.LoginAPI)
	// app.Router.HandleFunc("/api/logout", app.Logout)
	app.Router.HandleFunc("/api/quran/{id}", getQuranAPI).Methods("GET")
	app.Router.HandleFunc("/api/search/", SearchAPI)
	app.Router.HandleFunc("/api/surahs", getSurahs).Methods("GET")

	// Host the static folder for resources route
	fs := http.FileServer(http.Dir("./static/"))
	app.Router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

}
