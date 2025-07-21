import requests
url = "https://api.globalquran.com/surah/{id}/quran-uthmani-hafs?key=123"
for i in range(3, 115):
    response = requests.get(url.format(id=i))
    print(response.json())
    # Save it as json file
    if i < 10:
        i = f"00{i}"
    elif i < 100:
        i = f"0{i}"
    with open(f"{i}.json", "w") as f:
        f.write(response.text)