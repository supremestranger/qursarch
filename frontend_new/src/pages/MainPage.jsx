import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';

function MainPage() {
  const [surveys, setSurveys] = useState([]);
  const navigate = useNavigate();

  useEffect(() => {
    async function fetchSurveys() {
      try {
        const response = await fetch('/api/surveys', {
          credentials: 'include', // Include cookies
        });
        if (response.ok) {
          const data = await response.json();
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
        {surveys.length > 0 ? (
          surveys.map((survey) => (
            <li
              key={survey.id}
              className="border-b py-3 hover:bg-gray-100 cursor-pointer"
              onClick={() => navigate(`/survey/${survey.id}`)}
            >
              {survey.title}
            </li>
          ))
        ) : (
          <p>No surveys found.</p>
        )}
      </ul>
    </div>
  );
}

export default MainPage;
