import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import MainPage from './pages/MainPage';
import CreateSurveyPage from './pages/CreateSurveyPage';
import SurveyPage from './pages/SurveyPage';
import AuthPage from './pages/AuthPage';

function App() {
  return (
    <Router>
      <div className="min-h-screen bg-gray-100">
        <Routes>
          <Route path="/" element={<MainPage />} />
          <Route path="/create" element={<CreateSurveyPage />} />
          <Route path="/survey/:id" element={<SurveyPage />} />
          <Route path="/auth" element={<AuthPage />} />
        </Routes>
      </div>
    </Router>
  );
}

export default App;
