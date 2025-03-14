import React, { useState, useRef, useEffect } from 'react';
import './VerseCard.css';
import { FaPlay, FaBook, FaEllipsisV, FaCopy, FaClipboard } from 'react-icons/fa';

const VerseCard = ({ arabic, english, verseNumber, currentWord }) => {
    const [showMenu, setShowMenu] = useState(false);
    const menuRef = useRef(null);

    const arabicNumbers = ['٠', '١', '٢', '٣', '٤', '٥', '٦', '٧', '٨', '٩'];
    const convertToArabicNumber = (num) => {
        return num.toString().split('').map(digit => arabicNumbers[parseInt(digit)]).join('');
    };

    const handleCopyText = () => {
        navigator.clipboard.writeText(arabic);
        setShowMenu(false);
        alert('Arabic text copied to clipboard');
    };

    const handleAdvancedCopy = () => {
        navigator.clipboard.writeText(`${arabic}\n${english}`);
        setShowMenu(false);
        alert('Arabic and translation text copied to clipboard');
    };

    const handleTafsir = () => {
        // Implement the logic to open the tafsir for the verse
        alert('Tafsir button clicked');
    };

    const handlePlayVerse = () => {
        // Implement the logic to play the verse
        alert('Play verse button clicked');
    };

    const handleClickOutside = (event) => {
        if (menuRef.current && !menuRef.current.contains(event.target)) {
            setShowMenu(false);
        }
    };

    useEffect(() => {
        document.addEventListener('mousedown', handleClickOutside);
        return () => {
            document.removeEventListener('mousedown', handleClickOutside);
        };
    }, []);

    const highlightWord = (text) => {
        console.log(`Highlighting word: ${currentWord}`);
        return text.split(' ').map((word, index) => (
            <span key={index} className={word === currentWord ? 'highlighted-word' : ''}>
                {word}{' '}
            </span>
        ));
    };

    return (
        <div className="verse-card">
            <div className="verse-arabic" dir="rtl">
                {highlightWord(arabic)} <span className="verse-end-icon">{convertToArabicNumber(verseNumber)}</span>
            </div>
            <div className="verse-english">
                {highlightWord(english)}
            </div>
            <div className="verse-menu">
                <button onClick={handlePlayVerse}><FaPlay /></button>
                <button onClick={handleTafsir}><FaBook /></button>
                <button onClick={() => setShowMenu(!showMenu)}><FaEllipsisV /></button>
                {showMenu && (
                    <div className="menu-dropdown" ref={menuRef}>
                        <button onClick={handleCopyText}>Copy Arabic Only<FaCopy /></button>
                        <button onClick={handleAdvancedCopy}>Advance Copy<FaClipboard /></button>
                    </div>
                )}
            </div>
        </div>
    );
};

export default VerseCard;