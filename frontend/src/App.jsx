import {useState}from 'react';
import logo from './assets/images/logo-universal.png';
import './App.css';
import {Greet} from "../wailsjs/go/main/App";
import {Log} from "../wailsjs/go/main/App";

function App() {
    const [resultText, setResultText] = useState("Please enter your coords below in format: latitude,longitude 👇");
    const [coords, setCoords] = useState('');
    const [greeted, setGreeted] = useState(false);
    const [weatherData, setWeatherData] = useState(null);

    const updateCoords = (e) => setCoords(e.target.value);
    const updateResultText = (result) => setResultText(result);



    function greet() {
        Greet(coords).then((response) => {
            try {
                const data = JSON.parse(response);
                setWeatherData(data);

                setResultText("Weather data fetched successfully!");
            }   catch (error) {
                    setResultText("Error parsing weather data!");
            }
        });
        setGreeted(true);
    }

    const getWeatherIconSrc = (symbolCode) => {
        const iconPath = `/assets/icons/${symbolCode}.svg`;
        return iconPath;
    };

    return (
        <div id="App">
            {!greeted && (
                <div>
                    <img src={logo} id="logo" alt="logo"/>
                    <div id="result" className="result">{resultText}</div>
                    <div id="input" className="input-box">
                        <input id="name" className="input" onChange={updateCoords}
                               autoComplete="off" name="input" type="text"/>
                        <button className="btn" onClick={greet}>Greet</button>
                    </div>
                </div>
            )}
            {greeted && weatherData && (
                <div>
                    <div class="main">
                        <div class="left">
                            <div class="upper">
                                <p class="day">Day</p>
                                <p class="date">date</p>
                                <p class="location">location</p>
                            </div>
                            <div class="lower">
                                <p>{weatherData.time}</p>
                                <img src={getWeatherIconSrc(weatherData.symbol_code)} alt={weatherData.symbol_code} className="weather-icon"/>
                            </div>
                        </div>
                        <div class="right">
                            <div class="details">details</div>
                            <div class="weekly">weekly</div>
                            <div class="change">change</div>
                        </div>
                    </div>
                </div>
            )}
        </div>
    )
}

export default App
