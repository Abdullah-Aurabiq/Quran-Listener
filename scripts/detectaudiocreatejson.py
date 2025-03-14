import whisper
import json


def generate_word_subtitles(audio_path):
    """Generate subtitles word by word using OpenAI Whisper."""
    model = whisper.load_model("large")  # Use a larger model for better accuracy
    result = model.transcribe(audio_path)
    words = []
    for segment in result["segments"]:
        start = segment["start"]
        end = segment["end"]
        text = segment["text"]
        duration = end - start

        # Split the text into words and distribute timing proportionally
        word_list = text.split()
        num_words = len(word_list)
        word_duration = duration / num_words

        for i, word in enumerate(word_list):
            words.append(
                {
                    "word": word,
                    "start": start + i * word_duration,
                    "end": start + (i + 1) * word_duration,
                }
            )
    print(words)
    return words


def create_transcription_json(audio_path, json_path):
    words = generate_word_subtitles(audio_path)
    with open(json_path, "w", encoding="utf-8") as f:
        json.dump(words, f, ensure_ascii=False, indent=4)


if __name__ == "__main__":
    audio_path = "108.mp3"
    json_path = "108_transcription.json"
    create_transcription_json(audio_path, json_path)
    print(f"Transcription saved to {json_path}")
