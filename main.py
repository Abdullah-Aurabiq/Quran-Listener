# Python program to translate
# speech to text and text to speech


import speech_recognition as sr
import pyttsx3 

# Initialize the recognizer 
print("hello")
# r = sr.Recognizer()

# Function to convert text to
# speech
def SpeakText(command):
	
	# Initialize the engine
	engine = pyttsx3.init()

	engine.say(command) 
	engine.runAndWait()
	
	
SpeakText("Abdullah, How are you?")