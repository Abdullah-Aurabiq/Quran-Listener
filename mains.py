# from PyQt6.QtCore import Qt
# from PyQt6.QtWidgets import *
# from PyQt6.QtGui import QKeySequence, QAction
# from PyQt6.QtWidgets import QWidget


# class MainWindow(QMainWindow):
#     def __init__(self) -> None:
#         super().__init__()
#         self.setStyleSheet("background-color:white;")
#         self.setWindowOpacity(0.3)


# app = QApplication([])
# app.setApplicationName("Text Editor")
# window = MainWindow()
# window.showFullScreen()
# app.exec()

import requests
from transformers import T5ForConditionalGeneration, T5Tokenizer

# Load T5 model and tokenizer
model_name = "t5-base"
tokenizer = T5Tokenizer.from_pretrained(model_name)
model = T5ForConditionalGeneration.from_pretrained(model_name)


# Define the search function using the Quran.com API
def search_quran(prompt, language="en"):
    api_url = "https://api.quran.com/api/v4/search"
    params = {"q": prompt, "language": language}
    response = requests.get(api_url, params=params)
    if response.status_code == 200:
        return response.json()["search"]
    else:
        return None


# Define a function to generate a response using T5
def generate_response(prompt, verses):
    context = "\n".join([f"Quran verse: {verse['text']}" for verse in verses])
    full_prompt = f"Generate a response based on the following Para:"

    inputs = tokenizer.encode(
        full_prompt, return_tensors="pt", max_length=512, truncation=True
    )
    outputs = model.generate(inputs, max_length=150, num_beams=4, early_stopping=True)
    response = tokenizer.decode(outputs[0], skip_special_tokens=True)

    return response


# User prompt
user_prompt = "Tell me what Quran says about not to fear anything except Allah"

# Search for relevant verses
results = search_quran(user_prompt)

# Respond to user
if results:
    print(results)
    verses = results["results"][0]["translations"]
    response = generate_response(user_prompt, verses)
    print("Here is the response based on the Quran:")
    print(response)
else:
    print("No relevant verses found.")
