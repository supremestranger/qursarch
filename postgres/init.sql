CREATE TABLE Accounts (
    ID SERIAL NOT NULL PRIMARY KEY,
    Username VARCHAR,
    Password VARCHAR
);

CREATE TABLE Surveys (
    ID SERIAL NOT NULL PRIMARY KEY,
    Questions VARCHAR,
    Creator SERIAL references Accounts(ID)
);

INSERT INTO Accounts (Username, Password) values ('Admin', '123123');