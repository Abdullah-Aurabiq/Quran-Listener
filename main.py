import sys
import speech_recognition as sr
from PyQt5.QtWidgets import QApplication, QLabel, QVBoxLayout, QWidget
import pyttsx3


class SpeechRecognizer(QWidget):
    def __init__(self):
        super().__init__()

        # Set up the UI
        self.setWindowTitle("Urdu Speech Recognition")
        self.setGeometry(100, 100, 400, 200)
        self.layout = QVBoxLayout()
        self.label = QLabel("Listening...")
        self.layout.addWidget(self.label)
        self.setLayout(self.layout)

        # Set up the recognizer and microphone
        self.recognizer = sr.Recognizer()
        self.microphone = sr.Microphone()

        # Initialize the TTS engine
        self.tts_engine = pyttsx3.init()

        # Start listening in a background thread
        self.listen_in_background()

    def listen_in_background(self):
        with self.microphone as source:
            self.recognizer.adjust_for_ambient_noise(source)
        self.stop_listening = self.recognizer.listen_in_background(
            self.microphone, self.callback
        )

    def say_text(self, text):
        # text = "text"
        # Initialize the engin
        self.engine = pyttsx3.init()
        self.engine.say(text)
        self.engine.runAndWait()

    def callback(self, recognizer, audio):
        try:
            # Recognize speech using Google Web Speech API
            text = recognizer.recognize_google(audio, language="ur")
            self.label.setText(f"You said: {text}")

            # Speak the recognized text
            self.say_text(text)

            if "الله" in text:
                print("Yes")
        except sr.UnknownValueError:
            print("Could not understand the audio")
        except sr.RequestError as e:
            self.label.setText(f"Could not request results; {e}")


if __name__ == "__main__":
    app = QApplication(sys.argv)
    recognizer = SpeechRecognizer()
    recognizer.show()
    sys.exit(app.exec_())
