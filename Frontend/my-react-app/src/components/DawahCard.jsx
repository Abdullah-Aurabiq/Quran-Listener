import React from 'react';
import './DawahCard.css';

const DawahCard = ({text}) => {
    return (
        <div className="dawah-card">
            <div className="dawah-arabic">
                {text}
            </div>
        </div>
    );
};

export default DawahCard;