import React, { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';

function MainPage() {
  const [surveys, setSurveys] = useState([]);

  useEffect(() => {
    const fetchSurveys = async () => {
      try {
        const response = await fetch('/api/surveys');
        if (response.ok) {
          const data = await response.json();
          setSurveys(data);
        } else {
          console.error('Failed to fetch surveys:', await response.text());
        }
      } catch (error) {
        console.error('Error fetching surveys:', error);
      }
    };

    fetchSurveys();
  }, []);

  return (
    <div className="p-5">
      <h1 className="text-2xl font-bold mb-4">Опросы</h1>
      <ul className="space-y-2">
        {surveys.map((survey) => (
          <li key={survey.id} className="bg-white p-3 rounded shadow">
            <Link to={`/survey/${survey.id}`} className="text-blue-500">
              {survey.title}
            </Link>
          </li>
        ))}
      </ul>
      <div className="mt-5">
        <Link to="/create" className="bg-blue-500 text-white px-4 py-2 rounded">
          Создать новый опрос
        </Link>
      </div>
    </div>
  );
}

export default MainPage;
