package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

// App represents the application

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

var store = sessions.NewCookieStore([]byte("`giytIBFi0<.-3U,y$<!fdQFx$?@)}"))

// New creates a new instance of App
func (app *App) Initialize(Dbuser string, Dbpassword string, Dbname string) error {
	connectionString := fmt.Sprintf("%v:%v@tcp(127.0.0.1:3306)/%v", Dbuser, Dbpassword, Dbname)
	var err error
	app.DB, err = sql.Open("mysql", connectionString)
	if err != nil {
		return err
	}

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

func (app *App) GetStudentClasses(w http.ResponseWriter, r *http.Request) {

	classes, err := GetClasses(app.DB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Send the response with the template
	sendResponse(w, http.StatusOK, classes, "templates/classes.html", "text/html; charset=UTF-8")

}
func (app *App) CreateStudentClasses(w http.ResponseWriter, r *http.Request) {

	var cl Class
	if r.FormValue("coursename") != "" {

		TeacherID, _ := strconv.Atoi(r.FormValue("teacherid"))
		CourseName := r.FormValue("coursename")
		startTimeStr := r.FormValue("starttime")
		endTimeStr := r.FormValue("endtime")

		cl.TeacherID = TeacherID
		cl.CourseName = CourseName
		cl.Start_Time = startTimeStr
		cl.End_Time = endTimeStr
		fmt.Println(cl)
		err := cl.CreateClass(app.DB)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// sendResponse(w, http.StatusOK, nil, "templates/classcreated.html", "text/html; charset=UTF-8")

	}
	classes, err := GetClasses(app.DB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sendResponse(w, http.StatusOK, classes, "templates/createclass.html", "text/html; charset=UTF-8")

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
		http.Redirect(w, r, "/dashboard?check=1", http.StatusFound)
		return
	} else {
		sendResponse(w, http.StatusOK, nil, "templates/intro.html", "text/html; charset=UTF-8")
	}
}
func contact_us(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	auth, ok := session.Values["authenticated"].(bool)
	if !ok || !auth {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	data := map[string]interface{}{
		"user": session.Values,
	}
	sendResponse(w, http.StatusOK, data, "templates/contact-us.html", "text/html; charset=UTF-8")
}
func about_us(w http.ResponseWriter, r *http.Request) {
	sendResponse(w, http.StatusOK, nil, "templates/about-us.html", "text/html; charset=UTF-8")
}
func (app *App) StudentStart(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	auth, ok := session.Values["authenticated"].(bool)
	if !ok || !auth {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	// Get all courses
	courses, err := GetAllCourses(app.DB)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	fmt.Println(courses)
	fmt.Println(session.Values)
	// Create a new payload
	data := map[string]interface{}{
		"courses": courses,
		"user":    session.Values,
	}
	// Send the response with the template
	sendResponse(w, http.StatusOK, data, "templates/start.html", "text/html; charset=UTF-8")
}
func (app *App) MyStudents(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	auth, ok := session.Values["authenticated"].(bool)
	if !ok || !auth {
		http.Redirect(w, r, "/login/teacher", http.StatusFound)
		return
	}
	// Get all courses
	students, err := GetAllStudents(app.DB)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	fmt.Println(students)
	// Send the response with the template
	sendResponse(w, http.StatusOK, students, "templates/mystudents.html", "text/html; charset=UTF-8")
}

func (app *App) StudentSignup(w http.ResponseWriter, r *http.Request) {
	var st Student
	// Get the component from url
	// vars := mux.Vars(r)
	// fmt.Println(vars)

	// // Validate the input
	name := r.FormValue("name")
	age := r.FormValue("age")
	gender := r.FormValue("gender")
	email := r.FormValue("email")
	password := r.FormValue("password")
	conmfirm_password := r.FormValue("confirmpassword")
	fmt.Println(name, email, password, conmfirm_password)
	st.Name = name
	st.Email = email
	st.Password = password
	st.Gender = gender
	st.Age, _ = strconv.Atoi(age)
	if password != conmfirm_password {
		sendError(w, http.StatusBadRequest, "Passwords do not match")
		return
	}

	if name != "" || email != "" || password != "" {
		// Create Student account
		err := st.CreateStudentAccount(app.DB)
		if err != nil {
			sendError(w, http.StatusInternalServerError, err.Error())
			return
		}
		// Create a session for the user
		session, _ := store.Get(r, "session")
		session.Values["email"] = st.Email
		session.Values["name"] = st.Name
		session.Values["authenticated"] = true
		err = session.Save(r, w)
		if err != nil {
			sendError(w, http.StatusInternalServerError, "Failed to save session")
			return
		}
		http.Redirect(w, r, "/start/", http.StatusFound)
	}

	sendResponse(w, http.StatusOK, nil, "templates/signup.html", "text/html; charset=UTF-8")
}

func (app *App) StudentLogin(w http.ResponseWriter, r *http.Request) {
	// Parse form values
	email := r.FormValue("email")
	password := r.FormValue("password")

	// Validate input
	if email == "" || password == "" {
		sendResponse(w, http.StatusOK, nil, "templates/login.html", "text/html; charset=UTF-8")
		// sendError(w, http.StatusBadRequest, "Email and password are required")
		return
	}
	var st Student
	st.Email = email
	st.Password = password
	err := st.LoginStudent(app.DB)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Create a session for the user
	session, _ := store.Get(r, "session")
	session.Values["email"] = st.Email
	session.Values["name"] = st.Name
	session.Values["authenticated"] = true
	err = session.Save(r, w)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to save session")
		return
	}
	http.Redirect(w, r, "/start/", http.StatusFound)

}

func (app *App) TeacherLogin(w http.ResponseWriter, r *http.Request) {
	// Parse form values
	email := r.FormValue("email")
	password := r.FormValue("password")

	// Validate input
	if email == "" || password == "" {
		sendResponse(w, http.StatusOK, nil, "templates/tlogin.html", "text/html; charset=UTF-8")
		// sendError(w, http.StatusBadRequest, "Email and password are required")
		return
	}
	var t Teacher
	t.Email = email
	t.Password = password
	err := t.LoginTeacher(app.DB)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Create a session for the user
	session, _ := store.Get(r, "session")
	session.Values["email"] = t.Email
	session.Values["name"] = t.Name
	session.Values["authenticated"] = true
	err = session.Save(r, w)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to save session")
		return
	}
	http.Redirect(w, r, "/dashboard/teacher", http.StatusFound)

}

func (app *App) Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	delete(session.Values, "authenticated")
	session.Save(r, w)
	http.Redirect(w, r, "/login", http.StatusFound)
}

func (app *App) StudentDashboard(w http.ResponseWriter, r *http.Request) {
	check := r.FormValue("check")
	session, _ := store.Get(r, "session")
	auth, ok := session.Values["authenticated"].(bool)
	if !ok || !auth {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	var st Student
	st.Name = session.Values["name"].(string)
	st.Email = session.Values["email"].(string)
	student, err := st.GetStudent(app.DB)
	if err.Error() == "sql: no rows in result set" && check == "1" {
		http.Redirect(w, r, "/dashboard/teacher", http.StatusFound)
	}
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		fmt.Println(err.Error())
		return
	}
	data := map[string]interface{}{
		"user":    session.Values,
		"student": student,
	}
	sendResponse(w, http.StatusOK, data, "templates/dashboard.html", "text/html; charset=UTF-8")

}
func (app *App) TeacherDashboard(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	auth, ok := session.Values["authenticated"].(bool)
	if !ok || !auth {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	var t Teacher
	t.Name = session.Values["name"].(string)
	t.Email = session.Values["email"].(string)
	teacher, err := t.GetTeacher(app.DB)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	var cl Class
	if r.FormValue("coursename") != "" {

		TeacherID, _ := strconv.Atoi(r.FormValue("teacherid"))
		CourseName := r.FormValue("coursename")
		startTimeStr := r.FormValue("starttime")
		endTimeStr := r.FormValue("endtime")

		cl.TeacherID = TeacherID
		cl.CourseName = CourseName
		cl.Start_Time = startTimeStr
		cl.End_Time = endTimeStr
		fmt.Println(cl)
		err := cl.CreateClass(app.DB)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// sendResponse(w, http.StatusOK, nil, "templates/classcreated.html", "text/html; charset=UTF-8")

	}
	classes, err := GetClasses(app.DB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := map[string]interface{}{
		"user":    session.Values,
		"teacher": teacher,
		"classes": classes,
	}

	sendResponse(w, http.StatusOK, data, "templates/tdashboard.html", "text/html; charset=UTF-8")

}

func (app *App) TeacherStart(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	auth, ok := session.Values["authenticated"].(bool)
	if !ok || !auth {
		http.Redirect(w, r, "/login/teacher", http.StatusFound)
		return
	}
	// Get all courses
	courses, err := GetAllCourses(app.DB)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	fmt.Println(courses)
	fmt.Println(session.Values)
	// Create a new payload
	data := map[string]interface{}{
		"courses": courses,
		"user":    session.Values,
	}
	// Send the response with the template
	sendResponse(w, http.StatusOK, data, "templates/tstart.html", "text/html; charset=UTF-8")
}
func (app *App) Blogs(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	auth, ok := session.Values["authenticated"].(bool)
	if !ok || !auth {
		http.Redirect(w, r, "/home", http.StatusFound)
		return
	}
	// Create a new payload
	data := map[string]interface{}{
		"user": session.Values,
	}
	// Send the response with the template
	sendResponse(w, http.StatusOK, data, "templates/blogs.html", "text/html; charset=UTF-8")
}

func (app *App) handleRoutes() {
	app.Router.HandleFunc("/quran/{id}", getQuran).Methods("GET")
	app.Router.HandleFunc("/quran/", Home)
	app.Router.HandleFunc("/search/", Search)
	app.Router.HandleFunc("/hadith/", Hadith)
	app.Router.HandleFunc("/home/", Intro)
	app.Router.HandleFunc("/contact-us/", contact_us)
	app.Router.HandleFunc("/about-us/", about_us)
	app.Router.HandleFunc("/start/", app.StudentStart)
	app.Router.HandleFunc("/start/teacher", app.TeacherStart)
	app.Router.HandleFunc("/mystudents/", app.MyStudents)
	app.Router.HandleFunc("/register", app.StudentSignup)
	app.Router.HandleFunc("/login", app.StudentLogin)
	app.Router.HandleFunc("/login/teacher", app.TeacherLogin)
	app.Router.HandleFunc("/logout", app.Logout)
	app.Router.HandleFunc("/dashboard", app.StudentDashboard)
	app.Router.HandleFunc("/dashboard/teacher", app.TeacherDashboard)
	app.Router.HandleFunc("/cl", app.GetStudentClasses)
	app.Router.HandleFunc("/addclass", app.CreateStudentClasses)
	app.Router.HandleFunc("/blog", app.Blogs)
	// app.Router.HandleFunc("/start/", Start)

	// Host the static folder for resources route
	fs := http.FileServer(http.Dir("./static/"))
	app.Router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

}
