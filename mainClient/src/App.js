import React from 'react';
import 'jquery/dist/jquery.min.js';
import 'bootstrap/dist/css/bootstrap.min.css';
import 'bootstrap/dist/js/bootstrap.min.js';
import './App.css';
import './Components/Article/Article.js';
import Homepage from './Components/Homepage/Homepage.js';
import FileUpload from './Components/Repo/FileUpload.js';
import ArticleList from "./Components/Article/ArticleList";
import ArticlePage from "./Components/Article/ArticlePage";
import Header from './Components/Header/Header.js';
import Footer from './Components/Footer/Footer.js';
import {
    BrowserRouter as Router,
    Routes,
    Route
} from 'react-router-dom';

function App() {
  return (
      <Router>
          <div className="App d-flex flex-column min-vh-100">
              <Header/>
              <Routes>
                  <Route exact path='/' element={<Homepage/>}></Route>
                  <Route exact path='/papers' element={<ArticleList/>}></Route>
                  <Route exact path='/upload' element={<FileUpload/>}></Route>
                  <Route path="/papers/:id" element={<ArticlePage/>}></Route>
              </Routes>
              <Footer/>
          </div>
      </Router>
  );
}

export default App;