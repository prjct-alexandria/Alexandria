import * as React from "react";
import CreateArticle from "../Article/CreateArticle";
import Login from "../User/Login";
import Signup from "../User/Signup";
import { useState } from "react";
import isUserLoggedIn from "../User/AuthHelpers/isUserLoggedIn";

export default function HomepageHeader() {
  let [isLoggedIn, setLoggedIn] = useState<boolean>(isUserLoggedIn());

  // Listen for userAccountEvent that fires when user in localstorage changes
  window.addEventListener("userAccountEvent", () => {
    setLoggedIn(isUserLoggedIn());
  });

  return (
    <header>
      <div id="intro-example" className="p-5 text-center bg-image header-image">
        <div className="mask">
          <div className="d-flex justify-content-end align-items-end h-100"></div>
        </div>
      </div>
      <div className="header-buttons-wrapper d-flex justify-content-center align-items-center">
        <a
          className="btn btn-light btn-lg header-browse-articles"
          href="/articles"
          role="button"
          rel="nofollow"
        >
          Browse articles
        </a>
        {isLoggedIn && (
          <button
            type="button"
            className="btn btn-light btn-lg"
            data-bs-toggle="modal"
            data-bs-target="#publishArticle"
          >
            Publish article
          </button>
        )}
        {!isLoggedIn && (
          <div className="btn-group" role="group">
            <button
              type="button"
              className="btn btn-light btn-lg"
              data-bs-toggle="modal"
              data-bs-target="#login"
              id="btn-open-login-form"
            >
              Log in
            </button>

            <button
              type="button"
              className="btn btn-light btn-lg"
              data-bs-toggle="modal"
              data-bs-target="#signUp"
            >
              Sign up
            </button>
          </div>
        )}
      </div>
      <Login />
      <Signup />
      <CreateArticle />
    </header>
  );
}
