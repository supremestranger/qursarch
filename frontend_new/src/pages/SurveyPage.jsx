import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';

function SurveyPage() {
  const { id } = useParams();
  const [survey, setSurvey] = useState(null);

  useEffect(() => {
    const fetchSurvey = async () => {
      try {
        const response = await fetch(`http://localhost:3001/v1/surveys/${id}`);
        if (response.ok) {
          const data = await response.json();
          console.log(data);
          console.log(JSON.parse(data.Questions));
          setSurvey(data);
        } else {
          console.error('Failed to fetch survey:', await response.text());
        }
      } catch (error) {
        console.error('Error fetching survey:', error);
      }
    };

    fetchSurvey();
  }, [id]);

  if (!survey) {
    return <div className="p-5">Loading survey...</div>;
  }

  return (
    <div className="p-5">
      <h1 className="text-2xl font-bold mb-4">{survey.Title}</h1>
      <div>
        {JSON.parse(survey.Questions).map((q, index) => (
          <div key={index} className="mb-4">
            <p className="font-medium">{q.Title}</p>
            {q.Type === 'single_answer' && (
              q.Answers.map((answer, i) => (
                <label key={i} className="block">
                  <input type="radio" name={`question-${index}`} className="mr-2" />
                  {answer}
                </label>
              ))
            )}
            {q.Type === 'multiple_answers' && (
              q.Answers.map((answer, i) => (
                <label key={i} className="block">
                  <input type="checkbox" className="mr-2" />
                  {answer}
                </label>
              ))
            )}
            {q.Type === 's_text_answer' && <input type="text" className="border p-2 rounded w-full" />}
            {q.Type === 'l_text_answer' && <textarea className="border p-2 rounded w-full" />}
          </div>
        ))}
      </div>
    </div>
  );
}

export default SurveyPage;
