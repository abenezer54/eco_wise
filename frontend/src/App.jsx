import { useState, useEffect } from "react";
import "./App.css";
import axios from "axios";

function App() {
  const [data, setData] = useState(null);
  const [phoneNumber, setPhoneNumber] = useState("");
  const [language, setLanguage] = useState("english"); // Default language is English
  const [serverResponse, setServerResponse] = useState(null); // To store server response message

  // Fetch initial data from backend (if needed)
  useEffect(() => {
    axios
      .get("http://localhost:8080")
      .then((response) => {
        setData(response.data);
      })
      .catch((error) => {
        console.error("Error fetching data:", error);
      });
  }, []);

  const handlePhoneNumberChange = (event) => {
    setPhoneNumber(event.target.value);
  };

  const handleLanguageChange = (event) => {
    setLanguage(event.target.value);
  };

  const handleGetReward = () => {
    if (!phoneNumber) {
      setServerResponse(
        language === "english"
          ? "Please enter a valid phone number."
          : "እባኮትን ትክክለኛ ስልክ ቁጥር ያስገቡ።"
      );
      return;
    }

    const payload = { phone_number: phoneNumber.trim() };

    axios
      .post("http://localhost:8080/payment", payload, {
        headers: {
          "Content-Type": "application/json",
        },
      })
      .then((response) => {
        if (response.status === 200) {
          setServerResponse(
            response.data.message ||
              (language === "english"
                ? "Thank you for recycling! Your reward is being processed."
                : "እናመሰግናለን ለመቀመጥ። ሽልማትዎ በመሥራት ላይ ነው።")
          );
        } else {
          setServerResponse(
            language === "english"
              ? "Unexpected response from the server."
              : "ከመስመር ላይ ያልተጠበቀ ምላሽ።"
          );
        }
      })
      .catch((error) => {
        console.error("Error:", error);
        if (error.response) {
          setServerResponse(
            `Error: ${
              error.response.data.message ||
              (language === "english" ? "Request failed" : "ጥያቄ አልተሳካም።")
            }`
          );
        } else {
          setServerResponse(
            language === "english"
              ? "Error: Unable to connect to the server."
              : "እቅ፣ ከአገልግሎት አቅርቦት ጋር መገናኘት አልቻልኩም።"
          );
        }
      });
  };

  return (
    <div className="App">
      <div className="section welcome-screen">
        <h2>
          {language === "english"
            ? "Welcome to Eco-Wise Recycling Kiosk!"
            : "እንኳን ወደ Eco-Wise Recycling Kiosk  በደህና መጡ"}
        </h2>
        <p>
          {language === "english"
            ? "Join us in protecting the environment! Follow the steps below to recycle your plastic bottles and earn rewards."
            : "አካባቢን በመጠበቅ ላይ በጋራ እንቁም! የፕላስቲክ ጠርሙሶችዎን እንደገና ጥቅም ላይ ለማዋል እና ሽልማቶችን ለማግኘት ከታች ያሉትን ደረጃዎች ይከተሉ።"}
        </p>
        <select value={language} onChange={handleLanguageChange}>
          <option value="english">English</option>
          <option value="amharic">አማርኛ</option>
        </select>
      </div>

      <div className="section description">
        <h2>{language === "english" ? "About the Kiosk" : "ስለ ማስታወቂያ"}</h2>
        <p>
          {language === "english"
            ? "The Eco-Wise Smart Solutions Kiosk encourages recycling and sustainability. Deposit plastic bottles, earn rewards, and join the movement to reduce waste and pollution."
            : "አካባቢን በመጠበቅ ላይ! የፕላስቲክ ጠርሙሶችዎን እንደገና ጥቅም ላይ ለማዋል እና ሽልማቶችን ለማግኘት ከታች ያሉትን ደረጃዎች ይከተሉ።"}
        </p>
      </div>

      <div className="section phone-input">
        <h2>
          {language === "english"
            ? "Enter Your Phone Number"
            : "ስልክ ቁጥርዎን አስገቡ"}
        </h2>
        <p>
          {language === "english"
            ? "Provide your phone number to receive updates and rewards:"
            : "ማሻሻያዎችን እና ሽልማቶችን ለመቀበል ስልክ ቁጥርዎን ያቅርቡ።"}
        </p>
        <input
          type="text"
          placeholder={
            language === "english" ? "Enter phone number" : "ስልክ ቁጥር ያስገቡ"
          }
          value={phoneNumber}
          onChange={handlePhoneNumberChange}
        />
      </div>

      <div className="section get-reward">
        <button onClick={handleGetReward}>
          {language === "english" ? "Get Reward" : "ሽልማት ያግኙ"}
        </button>
      </div>

      {/* <div className="section server-response">
        <h2>{language === "english" ? "Server Response" : "የእርሶ ምላሽ"}</h2>
        <p>
          {language === "english"
            ? "Here's the data we received from the server:"
            : "ከመስመር የተቀበልነው መረጃ ይኸውና፡"}
        </p>
        <div className="response-container">
          <p>{serverResponse}</p>
        </div>
      </div> */}

      <div className="footer">
        <p>
          {language === "english"
            ? "Contact Support: 1-800-RECYCLE | Email: ekowise@gmail.com"
            : "እኛን ለመገናኘት፡ 1-800-RECYCLE | ኢሜል፡ ecowise123@gmail.com"}
        </p>
      </div>
    </div>
  );
}

export default App;
