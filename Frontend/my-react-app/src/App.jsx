import { useState, useEffect } from 'react';
import { BrowserRouter as Router, Route, Routes, useParams, useNavigate } from 'react-router-dom';
import './App.css';
import SurahList from './components/SurahList';
import SurahCard from './components/SurahCard';
import Quran from './components/Quran';
import axios from 'axios';
import Dawah from './components/Dawah';

function App() {
  const [surahs, setSurahs] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchSurahs = async () => {
      try {
        const response = await axios.get('https://electric-mistakenly-rat.ngrok-free.app/api/surahs');
        setSurahs(response.data);
        setLoading(false);
      } catch (error) {
        console.error('Error fetching Surahs:', error);
        setLoading(false);
      }
    };

    fetchSurahs();
  }, []);

  return (
    <Router>
      <Routes>
        <Route exact path="/" element={<Home surahs={surahs} loading={loading} />} />
        <Route path="/:surahId/:ayahId?" element={<SurahView />} />
        <Route path="/dawah" element={<Dawah/>} />
      </Routes>
    </Router>
  );
}

function Home({ surahs, loading }) {
  const [surahId, setSurahId] = useState(null);
  const navigate = useNavigate();

  const handleSelectSurah = (id) => {
    setSurahId(id);
    const selectedSurah = surahs.find((surah) => surah.id === id);
    if (selectedSurah) {
      navigate(`/${selectedSurah.id}`, { replace: true });
    }
  };

  if (loading) {
    return <div>Loading...</div>;
  }

  return (
    <div>
      <SurahList onSelectSurah={handleSelectSurah} />
      {!surahId && (
        <div className="surah-cards-container">
          {surahs.map((surah) => (
            <SurahCard key={surah.id} surah={surah} onSelectSurah={handleSelectSurah} />
          ))}
        </div>
      )}
      {surahId && <Quran surahId={surahId} />}
    </div>
  );
}

function SurahView() {
  const { surahId, ayahId } = useParams();

  return <Quran surahId={parseInt(surahId, 10)} ayahId={ayahId ? parseInt(ayahId, 10) : null} />;
}

export default App;