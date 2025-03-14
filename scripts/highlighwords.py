import sys
import json
from PyQt6.QtWidgets import (
    QApplication,
    QWidget,
    QVBoxLayout,
    QPushButton,
    QLabel,
    QFileDialog,
    QSlider,
)
from PyQt6.QtCore import *
from PyQt6.QtMultimedia import *
from PyQt6.QtMultimediaWidgets import *


class AudioHighlighter(QWidget):
    def __init__(self):
        super().__init__()
        self.initUI()

    def initUI(self):
        self.setWindowTitle("Audio Highlighter")
        self.setGeometry(100, 100, 1600, 900)  # Set a reasonable default size

        self.layout = QVBoxLayout()

        self.openAudioButton = QPushButton("Open Audio File")
        self.openAudioButton.clicked.connect(self.openAudioFile)
        self.layout.addWidget(self.openAudioButton)

        self.openTranscriptionButton = QPushButton("Open Transcription File")
        self.openTranscriptionButton.clicked.connect(self.openTranscriptionFile)
        self.layout.addWidget(self.openTranscriptionButton)

        self.openSurahButton = QPushButton("Open Surah File")
        self.openSurahButton.clicked.connect(self.openSurahFile)
        self.layout.addWidget(self.openSurahButton)

        self.audioLabel = QLabel("Audio File: None")
        self.layout.addWidget(self.audioLabel)

        self.transcriptionLabel = QLabel("Transcription File: None")
        self.layout.addWidget(self.transcriptionLabel)

        self.surahLabel = QLabel("Surah File: None")
        self.layout.addWidget(self.surahLabel)

        self.audioSlider = QSlider(Qt.Orientation.Horizontal)
        self.audioSlider.setRange(0, 100)
        self.audioSlider.sliderMoved.connect(self.setPosition)
        self.layout.addWidget(self.audioSlider)

        self.wordLabel = QLabel("")
        self.wordLabel.setAlignment(Qt.AlignmentFlag.AlignCenter)
        self.wordLabel.setStyleSheet("font-size: 24px;")
        self.layout.addWidget(self.wordLabel)

        self.setLayout(self.layout)

        self.player = QMediaPlayer()
        self.audioOutput = QAudioOutput()
        self.player.setAudioOutput(self.audioOutput)
        self.player.positionChanged.connect(self.updatePosition)
        self.player.durationChanged.connect(self.updateDuration)

        self.transcription = None
        self.currentWordIndex = 0

        self.timer = QTimer()
        self.timer.timeout.connect(self.highlightWord)

        self.surahData = None
        self.verseLabels = []
        self.surahWords = []

    def openAudioFile(self):
        audioFile, _ = QFileDialog.getOpenFileName(
            self, "Open Audio File", "", "Audio Files (*.mp3 *.wav)"
        )
        if audioFile:
            self.audioLabel.setText(f"Audio File: {audioFile}")
            self.player.setSource(QUrl.fromLocalFile(audioFile))
            self.player.play()
            self.timer.start(100)

    def openTranscriptionFile(self):
        transcriptionFile, _ = QFileDialog.getOpenFileName(
            self, "Open Transcription File", "", "JSON Files (*.json)"
        )
        if transcriptionFile:
            self.transcriptionLabel.setText(f"Transcription File: {transcriptionFile}")
            with open(transcriptionFile, "r", encoding="utf-8") as f:
                self.transcription = json.load(f)

    def openSurahFile(self):
        surahFile, _ = QFileDialog.getOpenFileName(
            self, "Open Surah File", "", "JSON Files (*.json)"
        )
        if surahFile:
            self.surahLabel.setText(f"Surah File: {surahFile}")
            with open(surahFile, "r", encoding="utf-8") as f:
                self.surahData = json.load(f)["ayahs"]
            self.extractSurahWords()
            self.displaySurah()

    def extractSurahWords(self):
        self.surahWords = []
        for ayah in self.surahData:
            words = ayah["text"].split()
            self.surahWords.extend(words)

    def displaySurah(self):
        for ayah in self.surahData:
            verseLabel = QLabel(ayah["text"])
            verseLabel.setAlignment(Qt.AlignmentFlag.AlignRight)
            verseLabel.setWordWrap(True)
            verseLabel.setTextFormat(Qt.TextFormat.RichText)
            verseLabel.hide()  # Hide the verse labels
            self.layout.addWidget(verseLabel)
            self.verseLabels.append(verseLabel)

    def setPosition(self, position):
        self.player.setPosition(position)

    def updatePosition(self, position):
        self.audioSlider.setValue(position)

    def updateDuration(self, duration):
        self.audioSlider.setRange(0, duration)

    def highlightWord(self):
        if self.transcription:
            currentTime = self.player.position() / 1000.0  # Convert to seconds
            while (
                self.currentWordIndex < len(self.transcription)
                and currentTime > self.transcription[self.currentWordIndex]["end"]
            ):
                self.currentWordIndex += 1
            if (
                self.currentWordIndex < len(self.transcription)
                and currentTime >= self.transcription[self.currentWordIndex]["start"]
            ):
                self.showCurrentWord(self.currentWordIndex)

    def showCurrentWord(self, wordIndex):
        if wordIndex < len(self.surahWords):
            currentWord = self.surahWords[wordIndex]
            self.wordLabel.setText(
                currentWord
            )  # Display the current word in the wordLabel


if __name__ == "__main__":
    app = QApplication(sys.argv)
    ex = AudioHighlighter()
    ex.showFullScreen()
    sys.exit(app.exec())
