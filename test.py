import requests

url = "https://stoplight.io/mocks/sunnah/api/352643496/collections"

headers = {"Accept": "application/json", "X-API-Key": "123"}

response = requests.get(url, headers=headers)

print(response.json())
exit()
import requests, sys, ast


class QG:
    SURAH_LAST = 114
    SURAH_FIRST = 1

    def __init__(
        self,
        url="http://api.globalquran.com/surah/",
        token="",
        lg_codes={},
        cache_size=0,
    ):
        """url and token are API url and token. lg_codes is a dict which maps
        2 letter language codes to their fullname. like:
        {"en": "en.sahih", "ar": "quran-simple"}
        """

        self.url = url
        self.token = token
        self.lg_codes = lg_codes
        self.cache_size = cache_size
        self.cache = list()

    def __update_cache(self):
        if not self.cache_size:
            return
        if len(self.cache) > self.cache_size:
            while len(self.cache) > self.cache_size:
                self.cache.pop()

    def getAudio(self, surah):
        """Download a Aduio of given url" """
        url = f"https://api.quran.com/api/v4/chapter_recitations/10/{surah}"
        # url = f"https://api.quran.com/api/v4/chapter_recitations/2/{surah}"

        payload = {}
        headers = {"Accept": "application/json"}

        response = requests.request("GET", url, headers=headers, data=payload)
        dictres = ast.literal_eval(response.text)
        result = dictres["audio_file"]["audio_url"]
        with open("audio.txt", "w", encoding="utf-8") as file:
            import json

            json.dump(result, file, ensure_ascii=False)
        # s = requests.get(result)
        # open('audio.mp3', 'wb').write(s.content)

    def getAyah(self, surah, ayah, lang):
        """Returns a dict containing the verse itself, ayah and surah number
        and id. example:
            {
            "surah": 1,
            "ayah": 1,
            "id": 1,
            "verse": "In the name of Allah, the Entirely Merciful, the Especially Merciful."
            }
        """
        if len(lang) == 2:
            if lang in self.lg_codes:
                lang = self.lg_codes[lang]
            else:
                raise ValueError(lang + " is not supported using 2letter codes")
        if not self.SURAH_FIRST <= surah <= self.SURAH_LAST:
            raise ValueError(
                "surah(chapter) must be between "
                + str(self.SURAH_FIRST)
                + " and "
                + str(self.SURAH_LAST)
            )

        # TODO: check if it's in cache, if so use cache, if not get and add to
        # cache
        if self.cache_size:
            for key, json_ in self.cache:  # FIXME: better name than json_
                if (surah, ayah, lang) == key:
                    self.cache.remove((key, json_))
                    self.cache.insert(0, (key, json_))
                    return json_

        req_url = self.url + str(surah) + "/" + lang
        ayah_json = requests.get(req_url, params={"key": self.token}).json()
        while len(ayah_json) == 1:
            ayah_json = ayah_json[next(iter(ayah_json))]

        # if ayah_json["surah"] != surah:
        #     raise ValueError("Invalid ayah number for this surah")

        self.cache.insert(0, ((surah, lang), ayah_json))
        self.__update_cache()
        return ayah_json


if __name__ == "__main__":
    from quran_suras import QuranSuras

    quran_suras = QuranSuras()
    insurah = int(sys.argv[1])

    # audio = quran_suras.get_sura_by_number(sura_number=surah, amount=1)['result'][-1]['url']
    # suras = requests.get(suras)
    # lang = "quran-uthmani"
    # quran-wordbyword
    # en.qaribullah
    lang = "quran-uthmani-hafs"
    lang_Translation = "en.qaribullah"
    Quran = QG()
    Quran.getAudio(surah=insurah)
    surahname = quran_suras.get_sura_name(sura_number=insurah)
    Qsurah = Quran.getAyah(insurah, insurah, lang)
    Tsurah = Quran.getAyah(insurah, insurah, lang_Translation)

    Qoutput = ""
    Toutput = ""

    s = int(list(Tsurah.keys())[0])
    for i in range(len(Qsurah)):
        Qayahnum = Qsurah[f"{s}"]["ayah"]
        Qoutput += Qsurah[f"{s}"]["verse"] + f"[{Qayahnum}]A8ea8"
        s += 1

    s = int(list(Qsurah.keys())[0])
    for i in range(len(Tsurah)):
        ayahnum = Tsurah[f"{s}"]["ayah"]
        Toutput += Tsurah[f"{s}"]["verse"] + f"[{ayahnum}]A8ea8"
        s += 1

    def Reverse(lst):
        new_lst = lst[::-1]
        return new_lst

    Qlist = Qoutput.split("A8ea8")
    Qlist.pop()

    Tlist = Toutput.split("A8ea8")
    Tlist.pop()
    a = {"ar": Qlist, "en": Tlist}

    with open("finale.txt", "w", encoding="utf-8") as file:
        import json

        json.dump(a, file, ensure_ascii=False)

    with open("surahname.txt", "w", encoding="utf-8") as file:
        import json

        json.dump(surahname, file, ensure_ascii=False)

    # with open("output.txt", "w", encoding="utf-8") as file:
    #     import json
    #     json.dump(Qoutput, file, ensure_ascii=False)

    # with open("translation.txt", "w", encoding="utf-8") as file:
    #     import json
    #     json.dump(Toutput, file, ensure_ascii=False)
