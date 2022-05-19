import React from "react";
import "jquery/dist/jquery.min.js";
import "bootstrap/dist/css/bootstrap.min.css";
import "bootstrap/dist/js/bootstrap.min.js";
import "./App.css";
import "./Components/Article/Article.js";
import Homepage from "./Components/Homepage/Homepage.js";
import Signup from "./Components/User/Signup.js";
import Login from "./Components/User/Login.js";
import FileUpload from "./Components/Repo/FileUpload.js";
import ArticleList from "./Components/Article/ArticleList";
import ArticlePage from "./Components/Article/ArticlePage";
import Header from "./Components/Header/Header.js";
import Footer from "./Components/Footer/Footer.js";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import VersionList from "./Components/Article/VersionList";

function App() {
  return (
    <Router>
      <div className="App d-flex flex-column min-vh-100">
        <Header />
        <Routes>
          <Route exact path="/" element={<Homepage />}></Route>
          <Route exact path="/articles" element={<ArticleList />}></Route>
          <Route exact path="/upload" element={<FileUpload />}></Route>
          <Route
            exact
            path="/articles/:aid/versions"
            element={<VersionList />}
          ></Route>
          <Route
            exact
            path="/articles/:articleId/versions/:versionId"
            element={<ArticlePage />}
          ></Route>
          <Route exact path="/signup" element={<Signup />}></Route>
          <Route exact path="/login" element={<Login />}></Route>
        </Routes>
        <Footer />
      </div>
    </Router>
  );
}

export default App;
