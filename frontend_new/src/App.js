import React, { useState } from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import MainPage from './pages/MainPage';
import CreateSurveyPage from './pages/CreateSurveyPage';
import SurveyPage from './pages/SurveyPage';
import AuthPage from './pages/AuthPage';

function App() {
  const [auth, setAuth] = useState(false);

  fetch(`http://localhost:3001/v1/accounts/check_auth`, { method: "POST", credentials: "include" })
    .then(function (response) { 
      if (response.ok) {
        setAuth(true);
        console.log("hello 123 123")
        return true;
      }else{
        setAuth(false);
        console.log("false")
        return false;
      }
    })
  return (
    <Router>
      <div className="min-h-screen bg-gray-100">
        <Routes>
          <Route path="/" element={<MainPage />} />
          <Route path="/create" element={auth ? <CreateSurveyPage /> : <Navigate to="/auth" />} />
          <Route path="/survey/:id" element={<SurveyPage />} />
          <Route path="/auth" element={<AuthPage />} />
        </Routes>
      </div>
    </Router>
  );
}

export default App;
