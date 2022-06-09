import * as React from "react";
import { Link } from "react-router-dom";
import CreateArticle from "./Article/CreateArticle";
import Login from "./User/Login";
import Signup from "./User/Signup";

export default function NavigationBar() {
  return (
    <div className="col-xs-12">
      <nav className="navbar navbar-expand-lg navbar-light bg-light">
        <a
          className="navbar-brand"
          href="/Users/Micutz/Software Project/alexandria/mainClient/public"
        >
          Project Alexandria
        </a>
        <button
          className="navbar-toggler"
          type="button"
          data-toggle="collapse"
          data-target="#navbarSupportedContent"
          aria-controls="navbarSupportedContent"
          aria-expanded="false"
          aria-label="Toggle navigation"
        >
          <span className="navbar-toggler-icon"></span>
        </button>

        <div className="collapse navbar-collapse" id="navbarSupportedContent">
          <ul className="navbar-nav mr-auto">
            <li className="nav-item active">
              <Link to="/" className="nav-link">
                Home
              </Link>
            </li>
            <li className="nav-item">
              <Link to="/articles" className="nav-link">
                Browse Articles
              </Link>
            </li>
            <li className="nav-item">
              <button
                type="button"
                className="nav-link btn"
                data-bs-toggle="modal"
                data-bs-target="#publishArticle"
              >
                Publish article
              </button>
              <CreateArticle />
            </li>
            <li className="nav-item">
              <button
                type="button"
                className="nav-link btn"
                data-bs-toggle="modal"
                data-bs-target="#signUp"
              >
                Sign up
              </button>
              <Signup />
            </li>
            <li className="nav-item">
              <button
                type="button"
                className="nav-link btn"
                data-bs-toggle="modal"
                data-bs-target="#login"
              >
                Log in
              </button>
              <Login />
            </li>
          </ul>
        </div>
      </nav>
    </div>
  );
}
