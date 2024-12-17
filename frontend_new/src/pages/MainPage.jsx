import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';

function MainPage() {
  const [surveys, setSurveys] = useState([]);
  const navigate = useNavigate();

  useEffect(() => {
    async function fetchSurveys() {
      try {
        const response = await fetch('http://localhost:3001/v1/surveys', {
        });
        if (response.ok) {
          const data = await response.json();
          console.log(data)
          setSurveys(data);
        } else {
          console.error('Failed to fetch surveys');
        }
      } catch (error) {
        console.error('Error fetching surveys:', error);
      }
    }
    fetchSurveys();
  }, []);

  return (
    <div className="p-5">
      {/* Header Section */}
      <div className="flex items-center justify-between mb-5">
        {/* Title */}
        <h1 className="text-2xl font-bold h-12 flex items-center">
          Опросы
        </h1>

        {/* Create Survey Button */}
        <button
          onClick={() => navigate('/create')}
          className="bg-green-500 text-white font-medium px-6 py-2 rounded-lg h-12 flex items-center"
        >
          Создать опрос
        </button>
      </div>

      {/* Surveys List */}
      <ul>
        {surveys != null && surveys.length > 0 ? (
          surveys.map((survey) => (
            <li
              key={survey.Id}
              className="border-b py-3 hover:bg-gray-100 cursor-pointer"
              onClick={() => navigate(`/survey/${survey.Id}`)}
            >
              {survey.Title}
            </li>
          ))
        ) : (
          <p>Опросов пока что нет.</p>
        )}
      </ul>
    </div>
  );
}

export default MainPage;
