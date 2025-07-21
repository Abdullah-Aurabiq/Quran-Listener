import React, { useState, useEffect, useRef } from 'react';
import axios from 'axios';
import VerseCard from './VerseCard';
import CustomDropdown from './CustomDropdown';
import SettingsMenu from './SettingsMenu';
import Modal from './Modal';
import './VoiceChatCard.css';

const Quran = ({ surahId }) => {
  const [verses, setVerses] = useState([]);
  const [reciters, setReciters] = useState([]);
  const [selectedReciter, setSelectedReciter] = useState(null);
  const [settings, setSettings] = useState({ glow: true, fontSize: 16, translationSize: 14 });
  const [showSettings, setShowSettings] = useState(false);
  const [showVoiceCard, setShowVoiceCard] = useState(false);
  const [currentWord, setCurrentWord] = useState('');
  const [showModal, setShowModal] = useState(false);
  const audioRef = useRef(null);

  const arabicVersion = "quran-uthmani-hafs";
  const englishVersion = "en.qaribullah";

  useEffect(() => {
    async function fetchData() {
      try {
        const paddedId = String(surahId).padStart(3, '0');

        // Arabic: from static local JSON
        const arRes = await fetch(`/quranar/surahs/${paddedId}.json`);
        const arData = await arRes.json();
        const arabicVerses = Object.values(arData.quran[arabicVersion]);

        // English: from Global Quran API
        const enRes = await fetch(`https://api.globalquran.com/surah/${surahId}/${englishVersion}`);
        const enData = await enRes.json();
        const englishVerses = Object.values(enData.quran[englishVersion]);

        // Merge
        const combined = englishVerses.map((en, i) => ({
          ar: arabicVerses[i]?.verse || '',
          en: en.verse || ''
        }));

        setVerses(combined);
      } catch (error) {
        console.error("Error fetching Quran data:", error);
      }
    }

    fetchData();
  }, [surahId]);

  useEffect(() => {
    async function fetchReciters() {
      try {
        const res = await axios.get('https://api.quran.com/api/v4/resources/recitations');
        const list = res.data.recitations;
        setReciters(list);
        setSelectedReciter(list[0]); // default
      } catch (err) {
        console.error("Failed to fetch reciters:", err);
      }
    }

    fetchReciters();
  }, []);

  const audioUrl = selectedReciter
    ? `${selectedReciter.audio_url}?chapter=${surahId}`
    : '';

  return (
    <div>
      <h1 id="surahname" style={{ color: "white", justifyContent: "center", display: "flex" }} dir="rtl">
        {String(surahId).padStart(3, '0')}surah
      </h1>

      <div>
        <button onClick={() => setShowSettings(!showSettings)}>Settings</button>
        {showSettings && <SettingsMenu settings={settings} setSettings={setSettings} />}
      </div>

      <div id="QuranVerse" style={{ color: 'white', margin: '10px', textWrap: 'wrap', pointerEvents: 'all' }} dir="rtl">
        {verses.map((v, i) => (
          <VerseCard
            key={i}
            arabic={v.ar}
            english={v.en}
            settings={settings}
            verseNumber={i + 1}
            currentWord={currentWord}
          />
        ))}
      </div>

      <button className="toggle-voice-card" onClick={() => setShowVoiceCard(!showVoiceCard)}>
        {showVoiceCard ? 'Hide Audio' : 'Show Audio'}
      </button>

      <div className={`voice-chat-card ${showVoiceCard ? 'show' : 'hide'}`}>
        <div className="voice-chat-card-header">
          <img className="avatar" />
          <div className="username" style={{ fontFamily: "surahs" }}>Surah {surahId}</div>
          <div className="username" style={{ fontFamily: "surahs" }}>{selectedReciter?.name}</div>
        </div>

        <div className="voice-chat-card-body">
          <div className="audio-container">
            <audio controls src={audioUrl} ref={audioRef}></audio>
          </div>
          <label>
            <CustomDropdown
              options={reciters.map(r => r.name)}
              selectedOption={selectedReciter?.name}
              onChange={(name) => setSelectedReciter(reciters.find(r => r.name === name))}
            />
          </label>
        </div>
      </div>

      <Modal show={showModal} onClose={() => setShowModal(false)}>
        <h2>Current Word</h2>
        <p>{currentWord}</p>
      </Modal>
    </div>
  );
};

export default Quran;
