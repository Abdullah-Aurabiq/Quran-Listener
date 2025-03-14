import React, { useState, useEffect, useRef } from 'react';
import axios from 'axios';
import './SurahList.css';

const surahNamesEnglish = [
    "Al-Fatiha", "Al-Baqarah", "Aale Imran", "An-Nisa", "Al-Ma'idah", "Al-Anam", "Al-A'raf", "Al-Anfal", "At-Tawbah",
    "Yunus", "Hud", "Yusuf", "Ar-Ra'd", "Ibrahim", "Al-Hijr", "An-Nahl", "Al-Isra", "Al-Kahf", "Maryam", "Ta-Ha",
    "Al-Anbiya", "Al-Hajj", "Al-Mu'minun", "An-Nur", "Al-Furqan", "Ash-Shu'ara", "An-Naml", "Al-Qasas", "Al-Ankabut",
    "Ar-Rum", "Luqmaan", "As-Sajda", "Al-Ahzaab", "Saba (surah)", "Faatir", "Yaseen", "As-Saaffaat", "Saad", "Az-Zumar",
    "Ghafir", "Fussilat", "Ash_Shooraa", "Az-Zukhruf", "Ad Dukhaan", "Al Jaathiyah", "Al-Ahqaaf", "Muhammad", "Al-Fath",
    "Al-Hujuraat", "Qaaf", "Adh-Dhaariyaat", "At-Toor", "An-Najm", "Al-Qamar", "Ar-Rahman", "Al-Waqi'a", "Al-Hadeed",
    "Al-Mujadila", "Al-Hashr", "Al-Mumtahanah", "As-Saff", "Al-Jumu'ah", "Al-Munafiqoon", "At-Taghabun", "At-Talaq",
    "At-Tahreem", "Al-Mulk", "Al-Qalam", "Al-Haaqqa", "Al-Ma'aarij", "Nooh", "Al-Jinn", "Al-Muzzammil", "Al-Muddaththir",
    "Al-Qiyamah", "Al-Insaan", "Al-Mursalaat", "An-Naba", "An-Naazi'aat", "Abasa", "At-Takweer", "Al-Infitar",
    "Al-Mutaffifeen", "Al-Inshiqaaq", "Al-Burooj", "At-Taariq", "Al-A'laa", "Al-Ghaashiyah", "Al-Fajr", "Al-Balad",
    "Ash-Shams", "Al-Layl", "Ad-Dhuha", "Ash-Sharh (Al-Inshirah)", "At-Teen", "Al-'Alaq", "Al-Qadr", "Al-Bayyinahh",
    "Az-Zalzalah", "Al-'Aadiyaat", "Al-Qaari'ah", "At-Takaathur", "Al-'Asr", "Al-Humazah", "Al-Feel", "Quraysh",
    "Al-Maa'oon", "Al-Kawthar", "Al-Kaafiroon", "An-Nasr", "Al-Masad", "Al-Ikhlaas", "Al-Falaq", "Al-Naas"
];

const SurahList = ({ onSelectSurah }) => {
    const [surahNames, setSurahNames] = useState([]);
    const [searchTerm, setSearchTerm] = useState('');
    const [isDropdownOpen, setIsDropdownOpen] = useState(false);
    const inputRef = useRef(null);
    const dropdownRef = useRef(null);

    useEffect(() => {
        const fetchSurahNames = async () => {
            try {
                const response = await axios.get('https://www.mp3quran.net/api/_arabic_sura.php');
                const surahData = response.data.Suras_Name.map((surah, index) => ({
                    id: surah.id,
                    arabicName: surah.name,
                    englishName: surahNamesEnglish[index]
                }));
                setSurahNames(surahData);
            } catch (error) {
                console.error('Error fetching Surah names:', error);
            }
        };

        fetchSurahNames();
    }, []);

    const handleFocus = () => {
        setIsDropdownOpen(true);
    };

    const handleBlur = () => {
        setTimeout(() => {
            setIsDropdownOpen(false);
        }, 200); // Delay to allow click event to register
    };

    const filteredSurahNames = surahNames.filter((surah) =>
        surah.arabicName.toLowerCase().includes(searchTerm.toLowerCase()) ||
        surah.englishName.toLowerCase().includes(searchTerm.toLowerCase()) ||
        surah.id.toString().includes(searchTerm)
    );

    return (
        <div>
            <div id="poda">
                <div className="glow"></div>
                <div className="darkBorderBg"></div>
                <div className="darkBorderBg"></div>
                <div className="darkBorderBg"></div>
                <div className="white"></div>
                <div className="border"></div>
                <div id="main">
                    <input
                        type="text"
                        placeholder="Search Surah"
                        ref={inputRef}
                        onFocus={handleFocus}
                        onBlur={handleBlur}
                        value={searchTerm}
                        onChange={(e) => setSearchTerm(e.target.value)}
                        name="text"
                        className="input"
                    />
                    <div id="input-mask"></div>
                    <div id="pink-mask"></div>
                    <div className="filterBorder"></div>
                    <div id="filter-icon">
                        <svg
                            preserveAspectRatio="none"
                            height="27"
                            width="27"
                            viewBox="4.8 4.56 14.832 15.408"
                            fill="none"
                        >
                            <path
                                d="M8.16 6.65002H15.83C16.47 6.65002 16.99 7.17002 16.99 7.81002V9.09002C16.99 9.56002 16.7 10.14 16.41 10.43L13.91 12.64C13.56 12.93 13.33 13.51 13.33 13.98V16.48C13.33 16.83 13.1 17.29 12.81 17.47L12 17.98C11.24 18.45 10.2 17.92 10.2 16.99V13.91C10.2 13.5 9.97 12.98 9.73 12.69L7.52 10.36C7.23 10.08 7 9.55002 7 9.20002V7.87002C7 7.17002 7.52 6.65002 8.16 6.65002Z"
                                stroke="#d6d6e6"
                                strokeWidth="1"
                                strokeMiterlimit="10"
                                strokeLinecap="round"
                                strokeLinejoin="round"
                            ></path>
                        </svg>
                    </div>
                    <div id="search-icon">
                        <svg
                            xmlns="http://www.w3.org/2000/svg"
                            width="24"
                            viewBox="0 0 24 24"
                            strokeWidth="2"
                            strokeLinejoin="round"
                            strokeLinecap="round"
                            height="24"
                            fill="none"
                            className="feather feather-search"
                        >
                            <circle stroke="url(#search)" r="8" cy="11" cx="11"></circle>
                            <line
                                stroke="url(#searchl)"
                                y2="16.65"
                                y1="22"
                                x2="16.65"
                                x1="22"
                            ></line>
                            <defs>
                                <linearGradient gradientTransform="rotate(50)" id="search">
                                    <stop stopColor="#f8e7f8" offset="0%"></stop>
                                    <stop stopColor="#b6a9b7" offset="50%"></stop>
                                </linearGradient>
                                <linearGradient id="searchl">
                                    <stop stopColor="#b6a9b7" offset="0%"></stop>
                                    <stop stopColor="#837484" offset="50%"></stop>
                                </linearGradient>
                            </defs>
                        </svg>
                    </div>
                </div>
            </div>

            {isDropdownOpen && (
                <div className="dropdown" ref={dropdownRef}>
                    <div className="surah-list">
                        {filteredSurahNames.map((surah) => (
                            <button key={surah.id} className="surah-button" onClick={() => { onSelectSurah(surah.id); setIsDropdownOpen(false); }}>
                                {surah.id}. {surah.englishName} ({surah.arabicName})
                            </button>
                        ))}
                    </div>
                </div>
            )}
        </div>
    );
};

export default SurahList;