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
                                <div class="day">{weatherData.week_day}</div>
                                <div class="date">{weatherData.date}</div>
                                <div class="location">{weatherData.city}, {weatherData.country_code}</div>
                            </div>
                            <div class="lower">
                                    <div class="symbol">
                                        <img src={getWeatherIconSrc(weatherData.symbol_code)} alt={weatherData.symbol_code} class="weather-icon"/>
                                    </div>
                                    <div class="temperature">
                                        {weatherData.temperature}°C
                                    </div>
                                    <div class="weather">
                                        {weatherData.symbol_code_nice}
                                    </div>
                            </div>
                        </div>
                        <div class="right">
                            <div class="details">
                                <div class="line">
                                    <div class="title">PRESSURE</div>
                                    <div class="value">{weatherData.air_pressure} hPa</div>
                                </div>
                                <div class="line">
                                    <div class="title">HUMIDITY</div>
                                    <div class="value">{weatherData.air_humidity} %</div>
                                </div>
                                <div class="line">
                                    <div class="title">WIND</div>
                                    <div class="value">{weatherData.wind_speed} m/s</div>
                                </div>
                            </div>
                            <div class="weekly">
                                <ul class="weekly-list">
                                    <li class="active">
                                        <img src={getWeatherIconSrc(weatherData.symbol_code)} alt={weatherData.symbol_code} class="day-icon"/>
                                        <span class="day">{weatherData.first_day}</span>
                                        <span class="day-temp">{weatherData.temperature}°C</span>
                                    </li>
                                    <li>
                                        <img src={getWeatherIconSrc(weatherData.second_symbol)} alt={weatherData.symbol_code} class="day-icon"/>
                                        <span class="day">{weatherData.second_day}</span>
                                        <span class="day-temp">{weatherData.second_temp}°C</span>
                                    </li>
                                    <li>
                                        <img src={getWeatherIconSrc(weatherData.third_symbol)} alt={weatherData.symbol_code} class="day-icon"/>
                                        <span class="day">{weatherData.third_day}</span>
                                        <span class="day-temp">{weatherData.third_temp}°C</span>
                                    </li>
                                    <li>
                                        <img src={getWeatherIconSrc(weatherData.fourth_symbol)} alt={weatherData.symbol_code} class="day-icon"/>
                                        <span class="day">{weatherData.fourth_day}</span>
                                        <span class="day-temp">{weatherData.fourth_temp}°C</span>
                                    </li>
                                </ul>
                            </div>
                            <div class="change">change</div>
                        </div>
                    </div>
                </div>
            )}
        </div>
    )
}

export default App
