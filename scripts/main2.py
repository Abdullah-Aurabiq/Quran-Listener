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
        self.label = QLabel("Processing audio file...")
        self.layout.addWidget(self.label)
        self.setLayout(self.layout)

        # Set up the recognizer
        self.recognizer = sr.Recognizer()

        # Initialize the TTS engine
        self.tts_engine = pyttsx3.init()

        # Load the audio file
        self.audio_file = sr.AudioFile(
            r"C:\Users\Black\Downloads\Audio-Introduction-0.1.wav"
        )

        # Process the audio file
        self.process_audio_file()

    def process_audio_file(self):
        with self.audio_file as source:
            audio = self.recognizer.record(source)
            try:
                # Recognize speech using Google Web Speech API
                text = self.recognizer.recognize_google(audio, language="en")
                self.label.setText(f"Audio file transcription: {text}")

                # Speak the recognized text
                self.say_text(text)

                if "code a amazon scraper using Python" in text:
                    print("Yes")
            except sr.UnknownValueError:
                print("Could not understand the audio")
            except sr.RequestError as e:
                self.label.setText(f"Could not request results; {e}")

    def say_text(self, text):
        # text = "text"
        # Initialize the engin
        self.engine = pyttsx3.init()
        self.engine.say(text)
        self.engine.runAndWait()


if __name__ == "__main__":
    app = QApplication(sys.argv)
    recognizer = SpeechRecognizer()
    recognizer.show()
    sys.exit(app.exec_())
