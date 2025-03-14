import React from 'react';
import './SettingsMenu.css';

const SettingsMenu = ({ settings, setSettings }) => {
    const handleChange = (e) => {
        const { name, value, type, checked } = e.target;
        console.log(name, value, type, checked);
        setSettings((prevSettings) => ({
            ...prevSettings,
            [name]: type === 'checkbox' ? checked : value,
            
        }));
    };

    return (
        <div className="settings-menu">
            <label>
                Glow:
                <input
                    type="checkbox"
                    name="glow"
                    checked={settings.glow}
                    onChange={handleChange}
                />
            </label>
            <label>
                Font Size:
                <input
                    type="number"
                    name="fontSize"
                    value={settings.fontSize}
                    onChange={handleChange}
                />
            </label>
            <label>
                Translation Size:
                <input
                    type="number"
                    name="translationSize"
                    value={settings.translationSize}
                    onChange={handleChange}
                />
            </label>
        </div>
    );
};

export default SettingsMenu;