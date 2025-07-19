import React, { useState } from 'react';
import DawahCard from './DawahCard';

const Dawah = () => {
    const [text, setText] = useState('');
    const [currentIndex, setCurrentIndex] = useState(0);

    const textArray = text.split('\n').filter(line => line.trim() !== '');

    const handleNext = () => {
        if (currentIndex < textArray.length - 1) {
            setCurrentIndex(currentIndex + 1);
        }
    };

    const handleBack = () => {
        if (currentIndex > 0) {
            setCurrentIndex(currentIndex - 1);
        }
    };

    return (
        <>
            <textarea
                value={text}
                onChange={(e) => setText(e.target.value)}
                placeholder="Enter text (each line becomes a card)"
                style={{
                    width: '100vw',
                    padding: '1vh',
                    fontSize: '2vh',
                    marginBottom: '1vh',
                    boxSizing: 'border-box'
                }}
            />
            <button onClick={() => setText('')} style={{ margin: '1vh' }}>Clear</button>

            <div id="DVerse" style={{  height: '60vh', overflow: 'hidden', position: 'relative' }}>


                <div className="slider-wrapper">
                    <div
                        className="slider"
                        style={{
                            transform: `translateX(-${currentIndex * 100}vw)`,
                            width: `${textArray.length * 100}vw`
                        }}
                    >
                        {textArray.map((line, index) => (
                            <div className="slide" key={index}>
                                <DawahCard text={line} />
                            </div>
                        ))}
                    </div>
                </div>

                <button className="nav-button back" onClick={handleBack} disabled={currentIndex === 0}>
                    {'<'}
                </button>
                <button className="nav-button next" onClick={handleNext} disabled={currentIndex === textArray.length - 1}>
                    {/* Adding > directly gives an error  */}
                    {'>'}
                </button>
            </div>
        </>
    );
};

export default Dawah;
