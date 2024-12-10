import React, { useState } from 'react';
import { Navigate } from 'react-router-dom';
import SurveyEditor from '../components/SurveyEditor';

function CreateSurveyPage() {
    const [title, setTitle] = useState('');
  const [error, setError] = useState('');

  const token = localStorage.getItem('token');
    if (!token) {
        return <Navigate to="/auth" />;
    }

  const submitSurvey = async (surveyData) => {
    

    if (!title.trim()) {
      setError('Survey title is required.');
      return;
    }

    const payload = {
      title,
      questions: surveyData.questions,
    };

    try {
      const response = await fetch('/api/surveys', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify(payload),
      });

      if (response.ok) {
        alert('Survey created successfully!');
      } else {
        const errorText = await response.text();
        setError(`Failed to create survey: ${errorText}`);
      }
    } catch (error) {
      setError('Error submitting survey.');
      console.error(error);
    }
  };

  return (
    <div className="p-5">
      <h1 className="text-2xl font-bold mb-4">Создать новый опрос</h1>
      {error && <p className="text-red-500 mb-3">{error}</p>}
      <input
        type="text"
        value={title}
        onChange={(e) => setTitle(e.target.value)}
        className="border p-2 rounded w-full mb-4"
        placeholder="Название Вашего опроса"
      />
      <SurveyEditor onSubmitSurvey={submitSurvey} />
    </div>
  );
}

export default CreateSurveyPage;
