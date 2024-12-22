
CREATE TABLE Admins (
    AdminID SERIAL PRIMARY KEY,
    Login VARCHAR(255) NOT NULL UNIQUE,
    Password VARCHAR(255) NOT NULL
);


CREATE TABLE Surveys (
    SurveyID SERIAL PRIMARY KEY,
    Title VARCHAR(255) NOT NULL,
    Description TEXT,
    CreatedAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CreatedBy INT REFERENCES Admins(AdminID) ON DELETE SET NULL
);


-- Создание типа ENUM для QuestionType
CREATE TYPE QuestionTypeEnum AS ENUM ('single_choice', 'multiple_choice', 'free_text');

CREATE TABLE Questions (
    QuestionID SERIAL PRIMARY KEY,
    SurveyID INT NOT NULL REFERENCES Surveys(SurveyID) ON DELETE CASCADE,
    QuestionText TEXT NOT NULL,
    QuestionType QuestionTypeEnum NOT NULL
);


CREATE TABLE AnswerOptions (
    OptionID SERIAL PRIMARY KEY,
    QuestionID INT NOT NULL REFERENCES Questions(QuestionID) ON DELETE CASCADE,
    OptionText TEXT NOT NULL
);


CREATE TABLE SurveyResults (
    ResultID SERIAL PRIMARY KEY,
    SurveyID INT NOT NULL REFERENCES Surveys(SurveyID) ON DELETE CASCADE,
    UserID VARCHAR(255) NOT NULL,
    CompletedAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE Answers (
    AnswerID SERIAL PRIMARY KEY,
    ResultID INT NOT NULL REFERENCES SurveyResults(ResultID) ON DELETE CASCADE,
    QuestionID INT NOT NULL REFERENCES Questions(QuestionID) ON DELETE CASCADE,
    AnswerText TEXT,
    SelectedOptions TEXT
);