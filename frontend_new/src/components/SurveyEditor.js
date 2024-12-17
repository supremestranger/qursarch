import React, { useState } from 'react';
import Question from './Question';

function SurveyEditor({ onSubmitSurvey }) {
  const [questions, setQuestions] = useState([]);

  const addQuestion = () => {
    setQuestions([
      ...questions,
      { id: Date.now(), text: 'Вопрос', type: 'text', options: [''] },
    ]);
  };

  const updateQuestion = (id, updatedData) => {
    setQuestions(questions.map((q) => (q.id === id ? { ...q, ...updatedData } : q)));
  };

  const deleteQuestion = (id) => {
    setQuestions(questions.filter((q) => q.id !== id));
  };

  const mapToServerSchema = () => {
    return {
      questions: questions.map((q) => {
        const mappedQuestion = {
          title: q.text,
          type: q.type === 'text' 
            ? (q.options.length > 20 ? 'l_text_answer' : 's_text_answer') 
            : q.type === 'radio' 
            ? 'single_answer'
            : 'multiple_answers',
          ...(q.type !== 'text' && { answers: q.options }),
        };
        return mappedQuestion;
      }),
    };
  };

  return (
    <div className="bg-white p-5 rounded shadow">
      {questions.map((question) => (
        <div key={question.id} className="mb-4 border border-solid border-b p-3 rounded">
          <Question
            question={question}
            onUpdate={updateQuestion}
            onDelete={deleteQuestion}
          />
        </div>
      ))}
      <button
        className="bg-blue-500 text-white px-4 py-2 rounded mr-2"
        onClick={addQuestion}
      >
        Добавить вопрос
      </button>
      <button
        className="bg-green-500 text-white px-8 absolute right-10 py-2 rounded"
        onClick={() => onSubmitSurvey(mapToServerSchema())}
      >
        Создать опрос
      </button>
    </div>
  );
}

export default SurveyEditor;
