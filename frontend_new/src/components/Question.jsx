import React from 'react';

function Question({ question, onUpdate, onDelete }) {
  const updateField = (field, value) => {
    onUpdate(question.id, { [field]: value });
  };

  const updateOption = (index, value) => {
    const updatedOptions = question.options.map((opt, i) =>
      i === index ? value : opt
    );
    onUpdate(question.id, { options: updatedOptions });
  };

  const addOption = () => {
    onUpdate(question.id, { options: [...question.options, ''] });
  };

  const removeOption = (index) => {
    const updatedOptions = question.options.filter((_, i) => i !== index);
    onUpdate(question.id, { options: updatedOptions });
  };

  return (
    <div>
      <input
        type="text"
        value={question.text}
        onChange={(e) => updateField('text', e.target.value)}
        className="border p-2 rounded w-full mb-2"
        placeholder="Enter your question"
      />
      <select
        value={question.type}
        onChange={(e) => updateField('type', e.target.value)}
        className="border border-solid border-black p-2 rounded w-full mb-2"
      >
        <option value="text">Текстовый ответ</option>
        <option value="radio">Один ответ</option>
        <option value="checkbox">Несколько ответов</option>
      </select>

      {question.type !== 'text' && (
        <div>
          {question.options.map((opt, index) => (
            <div key={index} className="flex items-center p-4 mb-2">
              <input
                type="text"
                value={opt}
                onChange={(e) => updateOption(index, e.target.value)}
                className="border p-2 rounded w-full"
                placeholder={`Ответ ${index + 1}`}
              />
              <button
                type="button"
                className="ml-2 text-red-500"
                onClick={() => removeOption(index)}
              >
                Удалить
              </button>
            </div>
          ))}
          <button
            type="button"
            className="text-blue-500 mt-2"
            onClick={addOption}
          >
            Добавить ответ
          </button>
        </div>
      )}

      <button
        type="button"
        className="text-red-500 mt-2"
        onClick={() => onDelete(question.id)}
      >
        Удалить вопрос
      </button>
    </div>
  );
}

export default Question;
