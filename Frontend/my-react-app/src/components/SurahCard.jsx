import React from 'react';
import './SurahCard.css';

const SurahCard = ({ surah, onSelectSurah }) => {
  return (
    <div className="surah-card" onClick={() => onSelectSurah(surah.id)}>
      <div className="surah-card-details">
        <h2 className="text-title">{surah.englishName} ({surah.arabicName})</h2>
        <p className="text-body"><strong>Meaning:</strong> {surah.englishMeaning}</p>
        <p className="text-body"><strong>Total Verses:</strong> {surah.totalVerses}</p>
        {/* <p className="text-body">{surah.startingVerses}</p> */}
      </div>
    </div>
  );
};

export default SurahCard;