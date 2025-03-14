import React, { useState } from 'react';
import './CustomDropdown.css';

const CustomDropdown = ({ options, selectedOption, onChange }) => {
    const [isOpen, setIsOpen] = useState(false);

    const handleSelect = (value) => {
        onChange(value);
        setIsOpen(false);
    };

    return (
        <div className="select" onClick={() => setIsOpen(!isOpen)}>
            Translation:
            <div className="selected">
                {selectedOption}
                <svg className="arrow" viewBox="0 0 24 24">
                    <path d="M7 10l5 5 5-5z" />
                </svg>
            </div>
            {isOpen && (
                <div className="options">
                    {options.map((option, index) => (
                        <div key={index} className="option" onClick={() => handleSelect(option)}>
                            {option}
                        </div>
                    ))}
                </div>
            )}
        </div>
    );
};

export default CustomDropdown;