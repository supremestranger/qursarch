import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';

function AuthPage() {
  const [login, setLogin] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const navigate = useNavigate();

  const handleAuth = async (isLogin) => {
    try {
      const response = await fetch(`http://localhost:3001/v1/accounts`, {
        method: `${isLogin ? 'GET' : 'POST'}`,
        body: isLogin ? null : JSON.stringify({ "login": login, "password": password }),
        credentials: "include"
      });

      if (response.ok) {
        const { token } = response.headers.getSetCookie();
        console.log(document.cookie)
        localStorage.setItem('token', token);
        navigate('/');
      } else {
        setError('Authentication failed. Check your credentials.');
      }
    } catch (error) {
      console.error('Auth error:', error);
      setError('Authentication error.');
    }
  };

  return (
    <div className="p-5">
      <h1 className="text-2xl font-bold mb-4">Register / Login</h1>
      {error && <p className="text-red-500 mb-3">{error}</p>}
      <input
        type="username"
        className="border p-2 rounded w-full mb-2"
        placeholder="Имя пользователя"
        value={login}
        onChange={(e) => setLogin(e.target.value)}
      />
      <input
        type="password"
        className="border p-2 rounded w-full mb-2"
        placeholder="Пароль"
        value={password}
        onChange={(e) => setPassword(e.target.value)}
      />
      <div>
        <button
          className="bg-blue-500 text-white px-4 py-2 rounded mr-2"
          onClick={() => handleAuth(true)}
        >
          Login
        </button>
        <button
          className="bg-green-500 text-white px-4 py-2 rounded"
          onClick={() => handleAuth(false)}
        >
          Register
        </button>
      </div>
    </div>
  );
}

export default AuthPage;
