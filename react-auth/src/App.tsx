import React, { useEffect, useState } from 'react';
import logo from './logo.svg';
import './App.css';
import Login from './pages/Login'
import { BrowserRouter, Navigate, Route, Routes } from 'react-router-dom';
import Home from './pages/Home'
import Register from './pages/Register'
import Error from './pages/Error'
import Footer from './components/footer';
import Navbar from './components/Navbar';
import Nav from './components/Navbar';


function App() {

  const [name, setName] = useState('')

  useEffect(() => {
    (
      async () => {
        const response = await fetch('http://localhost:8000/api/user', {
          headers: { 'Content-Type': 'application/json' },
          credentials: 'include'
        });

        const content = await response.json();

        if (content.message != "unauthenticated") {
          console.log(content)
          setName(content.name);
        }
      }
    )();
  });

  return (
    <div className="App">
      <BrowserRouter>
        <Nav name={name} setName={setName} />
        <div className="page-container">
          <div className="content-wrap">
            <main className="form-signin w-100 m-auto">
              <Routes>
                <Route path="/*" element={< Error />} />
                <Route path="/home" element={< Home name={name} />} />
                <Route path="/login" element={< Login setName={setName} />} />
                <Route path="/register" element={< Register />} />
                <Route path="/error" element={< Error />} />
              </Routes>
            </main>
          </div>
          < Footer />
        </div>
      </BrowserRouter>
    </div> 
  );
}

export default App;
