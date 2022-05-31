import * as React from "react";
import CreateArticle from "../Article/CreateArticle";
import Login from "../User/Login";
import Signup from "../User/Signup";

export default function HomepageHeader() {
  let loggedUser = localStorage.getItem("loggedUserEmail");

  return (
    <header>
      <div id="intro-example" className="p-5 text-center bg-image header-image">
        <div className="mask">
          <div className="d-flex justify-content-end align-items-end h-100">
            <div>
              <h1 className="mb-3">Alexandria</h1>
              <h5 className="mb-4">
                Collaborative and open-access scientific publishing.
              </h5>
              <a
                className="btn btn-outline-dark btn-lg m-2"
                href="/articles"
                role="button"
                rel="nofollow"
              >
                Browse articles
              </a>
              {loggedUser && (
                <div>
                  <button
                    type="button"
                    className="btn btn-outline-dark btn-lg m-2"
                    data-bs-toggle="modal"
                    data-bs-target="#publishArticle"
                  >
                    Publish article
                  </button>
                  <CreateArticle />
                </div>
              )}
              {!loggedUser && (
                <div>
                  <div className="btn-group" role="group">
                    <button
                      type="button"
                      className="btn btn-outline-dark btn-lg"
                      data-bs-toggle="modal"
                      data-bs-target="#login"
                    >
                      Log in
                    </button>

                    <button
                      type="button"
                      className="btn btn-outline-dark btn-lg"
                      data-bs-toggle="modal"
                      data-bs-target="#signUp"
                    >
                      Sign up
                    </button>
                  </div>
                  <Login />
                  <Signup />
                </div>
              )}
            </div>
          </div>
        </div>
      </div>
    </header>
  );
}
