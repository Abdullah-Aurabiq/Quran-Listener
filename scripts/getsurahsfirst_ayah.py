import requests
import json

# API endpoint to fetch Surah names and their first verses
SURAH_NAMES_API = "https://api.alquran.cloud/v1/surah"
FIRST_VERSE_API = "https://api.alquran.cloud/v1/ayah/{}/en"


def fetch_surah_data():
    response = requests.get(SURAH_NAMES_API)
    if response.status_code != 200:
        print("Failed to fetch Surah names")
        return []

    surah_data = response.json()["data"]
    surahs = []

    for surah in surah_data:
        surah_id = surah["number"]
        surah_name = surah["englishName"]
        surah_arabic_name = surah["name"]
        surah_english_meaning = surah["englishNameTranslation"]
        total_verses = surah["numberOfAyahs"]

        # Fetch the first verse of the Surah
        first_verse_response = requests.get(FIRST_VERSE_API.format(surah_id))
        if first_verse_response.status_code != 200:
            print(f"Failed to fetch first verse for Surah {surah_id}")
            continue

        first_verse = first_verse_response.json()["data"]["text"]

        surahs.append(
            {
                "id": surah_id,
                "englishName": surah_name,
                "arabicName": surah_arabic_name,
                "englishMeaning": surah_english_meaning,
                "totalVerses": total_verses,
                "startingVerses": first_verse,
            }
        )

    return surahs


def save_surah_data_to_json(surahs, filename="surahs.json"):
    with open(filename, "w", encoding="utf-8") as f:
        json.dump(surahs, f, ensure_ascii=False, indent=4)


if __name__ == "__main__":
    surah_data = fetch_surah_data()
    if surah_data:
        save_surah_data_to_json(surah_data)
        print("Surah data saved to surahs.json")
    else:
        print("No Surah data to save")
