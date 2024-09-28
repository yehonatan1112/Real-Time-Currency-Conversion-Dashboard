// eslint-disable-next-line no-unused-vars
import React, { useState, useEffect } from 'react';

function CurrencyConverter() {
    const [rates, setRates] = useState({});
    const [fromCurrency, setFromCurrency] = useState('USD');
    const [toCurrency, setToCurrency] = useState('EUR');
    const [amount, setAmount] = useState(1);
    const [convertedAmount, setConvertedAmount] = useState(null);
    const [error, setError] = useState(null);

    useEffect(() => {
        const socket = new WebSocket('ws://localhost:8080/live-updates');

        socket.onmessage = (event) => {
            const updatedRates = JSON.parse(event.data);
            setRates(updatedRates);
            setError(null);  // Clear any previous errors
        };

        socket.onerror = () => {
            setError("WebSocket connection error.");
        };

        return () => socket.close();  // Cleanup WebSocket on component unmount
    }, []);

    const handleConvert = () => {
        // Convert input to uppercase to match rates keys
        const from = fromCurrency.toUpperCase();
        const to = toCurrency.toUpperCase();

        if (rates[from] && rates[to]) {
            const rate = rates[to] / rates[from];
            setConvertedAmount((amount * rate).toFixed(2));
            setError(null);  // Clear any error
        } else {
            setError("Invalid currency code(s). Please check your input.");
        }
    };

    return (
        <div>
            <h1>Currency Converter</h1>
            {error && <p style={{ color: 'red' }}>{error}</p>} {/* Display error if any */}
            <div>
                <label>Amount: </label>
                <input type="number" value={amount} onChange={e => setAmount(e.target.value)} />
            </div>
            <div>
                <label>From: </label>
                <input type="text" value={fromCurrency} onChange={e => setFromCurrency(e.target.value)} />
            </div>
            <div>
                <label>To: </label>
                <input type="text" value={toCurrency} onChange={e => setToCurrency(e.target.value)} />
            </div>
            <button onClick={handleConvert}>Convert</button>
            {convertedAmount && <h2>Converted Amount: {convertedAmount} {toCurrency.toUpperCase()}</h2>}
        </div>
    );
}

export default CurrencyConverter;
