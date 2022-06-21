import * as React from "react";
import { Link } from "react-router-dom";
import CreateArticle from "./Article/CreateArticle";
import Login from "./User/Login";
import Signup from "./User/Signup";
import Logout from "./User/Logout";
import isUserLoggedIn from "./User/AuthHelpers/isUserLoggedIn";
import getLoggedInUsername from "./User/AuthHelpers/getLoggedInUsername";
import { useState } from "react";

export default function NavigationBar() {
  let [isLoggedIn, setLoggedIn] = useState<boolean>(isUserLoggedIn());
  let [loggedInUsername, setUsername] = useState<string>(
    getLoggedInUsername() || ""
  );

  // Listen for userAccountEvent that fires when user in localstorage changes
  window.addEventListener("userAccountEvent", () => {
    setLoggedIn(isUserLoggedIn());
    setUsername(getLoggedInUsername() || "");
  });

  return (
    <div className="col-xs-12">
      <nav className="navbar navbar-expand-lg fixed-top navbar-light bg-light">
        <div className="container-fluid">
          <Link to="/" className="navbar-brand">
            Project Alexandria
          </Link>
          <button
            className="navbar-toggler"
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
              <div className="d-flex flex-fill">
                <li className="nav-item d-flex active">
                  <Link to="/" className="nav-link">
                    Home
                  </Link>
                </li>
                <li className="nav-item d-flex">
                  <Link to="/articles" className="nav-link">
                    Browse Articles
                  </Link>
                </li>
              </div>
              <div className="d-flex">
                {!isLoggedIn && (
                  <li className="nav-item dropdown navbar-user-buttons">
                    <a
                      className="nav-link dropdown-toggle"
                      id="navbarDropdown"
                      role="button"
                      data-bs-toggle="dropdown"
                      aria-expanded="false"
                    >
                      Login
                    </a>

                    <ul
                      className="dropdown-menu"
                      aria-labelledby="navbarDropdown"
                    >
                      <li>
                        <button
                          role="button"
                          className="dropdown-item"
                          data-bs-toggle="modal"
                          data-bs-target="#login"
                        >
                          Login
                        </button>
                      </li>
                      <li>
                        <button
                          role="button"
                          className="dropdown-item"
                          data-bs-toggle="modal"
                          data-bs-target="#signUp"
                        >
                          Sign up
                        </button>
                      </li>
                    </ul>
                  </li>
                )}

                {isLoggedIn && (
                  <li className="nav-item dropdown navbar-user-buttons">
                    <svg
                      xmlns="http://www.w3.org/2000/svg"
                      width="16"
                      height="16"
                      fill="currentColor"
                      className="inline-block bi bi-person-circle"
                      viewBox="0 0 16 16"
                    >
                      <path d="M11 6a3 3 0 1 1-6 0 3 3 0 0 1 6 0z" />
                      <path
                        fillRule="evenodd"
                        d="M0 8a8 8 0 1 1 16 0A8 8 0 0 1 0 8zm8-7a7 7 0 0 0-5.468 11.37C3.242 11.226 4.805 10 8 10s4.757 1.225 5.468 2.37A7 7 0 0 0 8 1z"
                      />
                    </svg>
                    <a
                      className="nav-link dropdown-toggle inline-block "
                      href="#"
                      id="navbarDropdown"
                      role="button"
                      data-bs-toggle="dropdown"
                      aria-expanded="false"
                    >
                      {loggedInUsername}
                    </a>
                    <ul
                      className="dropdown-menu"
                      aria-labelledby="navbarDropdown"
                    >
                      <li>
                        <button
                          role="button"
                          className="nav-link btn dropdown-item"
                          data-bs-toggle="modal"
                          data-bs-target="#publishArticle"
                        >
                          Publish article
                        </button>
                      </li>
                      <li>
                        <hr className="dropdown-divider" />
                      </li>
                      <li>
                        <Logout />
                      </li>
                    </ul>
                  </li>
                )}
              </div>
            </ul>
          </div>
        </div>
      </nav>
      <Login />
      <Signup />
      <CreateArticle />
    </div>
  );
}
