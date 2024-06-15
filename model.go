package main

// Define structures to hold the JSON data
type QuranResponse struct {
	Quran QuranData `json:"quran"`
}

type QuranData struct {
	UthmaniHafs map[string]Verse `json:"quran-uthmani-hafs"`
}

type Verse struct {
	ID    int    `json:"id"`
	Surah int    `json:"surah"`
	Ayah  int    `json:"ayah"`
	Verse string `json:"verse"`
}
