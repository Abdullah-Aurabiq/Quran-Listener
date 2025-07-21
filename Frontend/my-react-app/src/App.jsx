import { useState, useEffect } from 'react';
import { BrowserRouter as Router, Route, Routes, useParams, useNavigate } from 'react-router-dom';
import './App.css';
import SurahList from './components/SurahList';
import SurahCard from './components/SurahCard';
import Quran from './components/Quran';
import Dawah from './components/Dawah';
import axios from 'axios';

function App() {
  const [surahs, setSurahs] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchSurahs = async () => {
      try {
        const response = await axios.get('/surahs.json'); // ✅ from /public
        setSurahs(response.data);
      } catch (error) {
        console.error('Error fetching Surahs:', error);
      } finally {
        setLoading(false);
      }
    };
    fetchSurahs();
  }, []);

  return (
    <Router>
      <Routes>
        <Route path="/" element={<Home surahs={surahs} loading={loading} />} />
        <Route path="/:surahId/:ayahId?" element={<SurahView />} />
        <Route path="/dawah" element={<Dawah />} />
      </Routes>
    </Router>
  );
}

function Home({ surahs, loading }) {
  const navigate = useNavigate();

  const handleSelectSurah = (id) => {
    navigate(`/${id}`);
  };

  if (loading) return <div>Loading...</div>;

  return (
    <div>
      <SurahList onSelectSurah={handleSelectSurah} />
      <div className="surah-cards-container">
        {surahs.map((surah) => (
          <SurahCard key={surah.id} surah={surah} onSelectSurah={handleSelectSurah} />
        ))}
      </div>
    </div>
  );
}

function SurahView() {
  const { surahId, ayahId } = useParams();
  return (
    <Quran
      surahId={parseInt(surahId, 10)}
      ayahId={ayahId ? parseInt(ayahId, 10) : null}
    />
  );
}

export default App;
