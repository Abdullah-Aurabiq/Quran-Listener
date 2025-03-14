package main

import (
	"database/sql"
	"fmt"
	"time"
)

// Define structures to hold the JSON data
type QuranResponse struct {
	Quran QuranData `json:"quran"`
}

type QuranData struct {
	UthmaniHafs map[string]Verse `json:"quran-uthmani-hafs"`
}

type Verse struct {
	Ayah  int    `json:"ayah"`
	ID    int    `json:"id"`
	Surah int    `json:"surah"`
	Verse string `json:"verse"`
}

// type Verse struct {
// 	Ayah  int    `json:"ayah"`
// 	Verse string `json:"text"`
// }

type Users struct {
	UserName string `json:"username"`
	Gender   string `json:"gender"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

// uint8 to time.Time

func Uint8ToTime(u []uint8) (time.Time, error) {
	timeString := string(u)
	layout := time.RFC1123Z // Adjust according to your datetime format
	parsedTime, err := time.Parse(layout, timeString)
	if err != nil {
		return time.Time{}, err
	}
	return parsedTime, nil
}

func (st *Users) CreateUserAccount(db *sql.DB) error {
	query := fmt.Sprintf("INSERT INTO users (username, gender, password, email) VALUES ('%v','%v','%v','%v')", st.UserName, st.Gender, st.Password, st.Email)
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (st *Users) UserLogin(db *sql.DB) error {
	query := fmt.Sprintf("SELECT UserName, Gender, email, password FROM users WHERE email='%v' and password='%v' LIMIT 1", st.Email, st.Password)
	row := db.QueryRow(query)
	err := row.Scan(&st.UserName, &st.Gender, &st.Email, &st.Password)
	return err
}
