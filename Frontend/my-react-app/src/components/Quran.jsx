import React, { useState, useEffect, useRef } from 'react';
import axios from 'axios';
import VerseCard from './VerseCard';
import CustomDropdown from './CustomDropdown';
import SettingsMenu from './SettingsMenu';
import Modal from './Modal';
import './VoiceChatCard.css'; // Import the CSS file for animations

const Quran = ({ surahId }) => {
    const [data, setData] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const [arabicVersion, setArabicVersion] = useState('quran-uthmani-hafs');
    const [translation, setTranslation] = useState('en.sahih');
    const [versions, setVersions] = useState([]);
    const [reciters, setReciters] = useState([]);
    const [selectedReciter, setSelectedReciter] = useState({});
    const [settings, setSettings] = useState({
        glow: true,
        fontSize: 16,
        translationSize: 14,
    });
    const [showSettings, setShowSettings] = useState(false);
    const [showVoiceCard, setShowVoiceCard] = useState(false); // State variable for voice card visibility
    const [currentVerseIndex, setCurrentVerseIndex] = useState(null);
    const [currentWord, setCurrentWord] = useState('');
    const [showModal, setShowModal] = useState(false);
    const audioRef = useRef(null);

    useEffect(() => {
        const fetchVersions = async () => {
            try {
                const response = await axios.get('http://localhost:1481/static/qurantext.json');
                setVersions(response.data.quranList);
            } catch (err) {
                console.error('Error fetching versions:', err);
            }
        };

        const fetchReciters = async () => {
            try {
                const response = await axios.get('http://localhost:1481/static/reciters.json');
                setReciters(response.data.reciters);
                setSelectedReciter(response.data.reciters[0]); // Set default reciter
            } catch (err) {
                console.error('Error fetching reciters:', err);
            }
        };

        fetchVersions();
        fetchReciters();
    }, []);

    useEffect(() => {
        const fetchData = async () => {
            try {
                const response = await axios.get(`http://localhost:1481/api/quran/${surahId}?ar=${arabicVersion}&translation=${translation}`);
                setData(response.data);
                console.log(response.data);
            } catch (err) {
                setError(err);
            } finally {
                setLoading(false);
            }
        };

        fetchData();
    }, [surahId, arabicVersion, translation]);

    useEffect(() => {
        const fetchTranscription = async () => {
            try {
                const transcription = await axios.get(`http://localhost:1481/static/067_transcription.json`);
                const words = transcription.data;
                console.log('Transcription data:', words);
            } catch (err) {
                console.error('Error fetching transcription:', err);
            }
        };

        fetchTranscription();

        const handleTimeUpdate = async () => {
            const currentTime = audioRef.current.currentTime;
            console.log(`Current time: ${currentTime}`);
            const transcription = await axios.get(`http://localhost:1481/static/108_transcription.json`);
            const words = transcription.data;

            for (let i = 0; i < words.length; i++) {
                if (currentTime >= words[i].start && currentTime <= words[i].end) {
                    console.log(`Current word: ${words[i].word}`);
                    setCurrentWord(words[i].word);
                    setCurrentVerseIndex(i);
                    setShowModal(true);
                    break;
                }
            }
        };

        const handlePlay = () => {
            console.log('Audio is playing');
        };

        const handlePause = () => {
            console.log('Audio is paused');
        };

        const handleEnded = () => {
            console.log('Audio has ended');
        };

        if (audioRef.current) {
            console.log('Adding event listeners to audio element');
            audioRef.current.addEventListener('timeupdate', handleTimeUpdate);
            audioRef.current.addEventListener('play', handlePlay);
            audioRef.current.addEventListener('pause', handlePause);
            audioRef.current.addEventListener('ended', handleEnded);
        } else {
            console.log('audioRef.current is null');
        }

        return () => {
            if (audioRef.current) {
                console.log('Removing event listeners from audio element');
                audioRef.current.removeEventListener('timeupdate', handleTimeUpdate);
                audioRef.current.removeEventListener('play', handlePlay);
                audioRef.current.removeEventListener('pause', handlePause);
                audioRef.current.removeEventListener('ended', handleEnded);
            }
        };
    }, [audioRef.current]);

    if (loading) return <div>Loading...</div>;
    if (error) return <div>Error: {error.message}</div>;

    const formatData = (finaleData) => {
        finaleData = JSON.parse(finaleData);
        return finaleData.map((item, index) => (
            <VerseCard key={index} arabic={item.ar} english={item.en} settings={settings} verseNumber={index + 1} currentWord={currentWord} />
        ));
    };

    const audioUrl = `${selectedReciter.Server}/${String(surahId).padStart(3, '0')}.mp3`;

    return (
        <div>

            <h1 id="surahname" style={{ fontFamily: "'surahs', sans-serif", color: "white", justifyContent: "center", display: "flex" }} dir="rtl">
                {data.quran.id}surah
            </h1>
            <div>
                <label>
                    Translation:
                    <CustomDropdown
                        options={Object.keys(versions)}
                        selectedOption={translation}
                        onChange={setTranslation}
                    />
                </label>
                <button onClick={() => setShowSettings(!showSettings)}>Settings</button>
                {showSettings && <SettingsMenu settings={settings} setSettings={setSettings} />}
            </div>
            <div id="QuranVerse" style={{ color: 'white', margin: '10px', textWrap: 'wrap', pointerEvents: 'all' }} dir="rtl">
                {formatData(data.quran.FinaleData)}
            </div>
            <button className="toggle-voice-card" onClick={() => setShowVoiceCard(!showVoiceCard)}>
                {showVoiceCard ? 'Hide Audio' : 'Show Audio'}
            </button>
            <div className={`voice-chat-card ${showVoiceCard ? 'show' : 'hide'}`}>
                <div className="voice-chat-card-header">
                    <img className="avatar" />
                    <div className="username" style={{ fontFamily: "surahs" }}>{data.quran.id}surah</div>
                    <div className="status"></div>
                    <div className="username" style={{ fontFamily: "surahs" }}>{selectedReciter.name}</div>
                </div>
                <div className="voice-chat-card-body">
                    <div className="audio-container">
                        <audio controls src={audioUrl} ref={audioRef}></audio>
                    </div>
                    <label>
                        <CustomDropdown
                            options={reciters.map(reciter => reciter.name)}
                            selectedOption={selectedReciter.name}
                            onChange={(name) => setSelectedReciter(reciters.find(reciter => reciter.name === name))}
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