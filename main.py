
import requests, time
from selenium import webdriver
from bs4 import BeautifulSoup
from selenium.webdriver.chrome.options import Options
# html = requests.get("https://quran.com/1?translations=158%2C131")
html = requests.get("https://previous.quran.com/1")
soup = BeautifulSoup(html.text, "html.parser")

def get_surahname(surahname:str):
    # result = soup.find(class_="Link_base__9W8Qs")
    surahnames = soup.find("a", class_="Link_base__9W8Qs", href=f"/{surahname}")
    # surahname = surahnames.find("span")
    # print(surahname)
    if surahnames:
        for surah in surahnames:
            print(surah.text, end="\n"*2)


def scrape_quran():
    url = 'https://quran.com/1'
    opt = webdriver.EdgeOptions()
    opt.headless = True
    # opt.add_argument('--headless=new')
    # Initialize Selenium WebDriver (make sure you have appropriate drivers installed)
    driver = webdriver.Edge(opt)  # Or specify the path to your webdriver
    # driver.minimize_window()
    driver.get(url)
    # driver.execute_script("alert('Hello');")
    # Wait for the page to load (adjust the sleep time as needed)
    time.sleep(8)
    
    # Get the page source after JavaScript execution
    page_source = driver.page_source.encode('utf16')
    
    # Close the WebDriver
    driver.quit()
    
    # Parse the page source with BeautifulSoup
    soup = BeautifulSoup(page_source, 'html.parser')
    
    # Find the elements containing the Quranic text
    verse_elements = soup.find_all('div', class_='TranslationViewCell_cellContainer__rhs1_')
    
    # Iterate through the verse elements and extract the text
    for verse in verse_elements:
        for verse in verse_elements:
            verse_text = verse.text
            print(verse_text)
    
def get_surah():


    # Find all <div> tags with class "TranslationViewCell_cellContainer__rhs1_"
    div_tags = soup.find("div", class_="translation__text")
    # print(div_tags.prettify())
    # surahname = div_tags.find("span")
    print(div_tags.text + "\n\n sd")
    # for div_tag in div_tags:
    #     print(div_tag.text, end="\n\n")

def main():
    get_surahname("1")
    get_surah()
    scrape_quran()

if __name__ == "__main__":
    main()
