import * as React from "react";
import "jquery/dist/jquery.min.js";
import "bootstrap/dist/css/bootstrap.min.css";
import "bootstrap/dist/js/bootstrap.min.js";
import "./App.scss";
import Homepage from "./Components/Homepage/Homepage";
import Signup from "./Components/User/Signup";
import Login from "./Components/User/Login";
import FileUpload from "./Components/Article/FileUpload";
import ArticleList from "./Components/Article/ArticleList";
import ArticleVersionPage from "./Components/Article/ArticleVersionPage";
import NavigationBar from "./Components/NavigationBar";
import Footer from "./Components/Footer";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import VersionList from "./Components/Article/VersionList";
import CompareView from "./Components/Article/CompareView";
import MRList from "./Components/Article/MRList";

function App() {
  return (
    <>
      <Router>
        <div className="App d-flex flex-column min-vh-100">
          <NavigationBar />
          <div>
            <Routes>
              <Route path="/" element={<Homepage />}></Route>
              <Route path="/articles" element={<ArticleList />}></Route>
              <Route path="/upload" element={<FileUpload />}></Route>
              <Route
                path="/articles/:articleId/versions"
                element={<VersionList />}
              ></Route>
              <Route
                path="/articles/:articleId/versions/:versionId"
                element={<ArticleVersionPage />}
              ></Route>
              <Route path="/signup" element={<Signup />}></Route>
              <Route path="/login" element={<Login />}></Route>
              <Route
                path="/articles/:articleId/requests/:requestId"
                element={<CompareView />}
              ></Route>
              <Route
                path="/articles/:articleId/requests/*"
                element={<MRList />}
              ></Route>
            </Routes>
          </div>
          <Footer />
        </div>
      </Router>
    </>
  );
}

export default App;
