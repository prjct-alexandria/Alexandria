import * as React from "react";
import { useEffect, useState } from "react";
import { Link, useParams, useSearchParams } from "react-router-dom";
import LoadingSpinner from "../LoadingSpinner";
import CreateMR from "./CreateMR";
import ThreadList from "./ThreadList";
import CreateArticleVersion from "./CreateArticleVersion";
import FileDownload from "./FileDownload";
import configData from "../../config.json";
import NotificationAlert from "../NotificationAlert";
import isUserLoggedIn from "../User/AuthHelpers/isUserLoggedIn";
import InlineEditor from "./InlineEditor";
import getLoggedInEmail from "../User/AuthHelpers/getLoggedInEmail";

type ArticleVersion = {
  owners: Array<string>;
  title: string;
  content: string;
  latestHistoryID: string;
};

export default function ArticleVersionPage() {
  let [versionData, setData] = useState<ArticleVersion>();
  let [isLoaded, setLoaded] = useState<boolean>(false);
  let [error, setError] = useState<Error>();
  let [isLoggedIn, setLoggedIn] = useState<boolean>(isUserLoggedIn());
  let [isOutdated, setOutdated] = useState<boolean>();
  let [isUserAnOwner, setIsOwner] = useState<boolean>(false);
  let [editorContent, setContent] = useState<string>("");

  function isLoggedInUserTheOwner() {
    return versionData && versionData.owners.includes(getLoggedInEmail() || "");
  }

  // Listen for userAccountEvent that fires when user in localstorage changes
  window.addEventListener("userAccountEvent", () => {
    setLoggedIn(isUserLoggedIn());
  });

  window.addEventListener("changesSavedEvent", () => {
    // This is where you would do setData({content:editorContent})
    // if(versionData) setData((versionData =>{...versionData, "content":editorContent}));
    // Reloading the page isn't great for a single page app
    window.location.reload();
  });

  let params = useParams();

  let url = //"/article_version1.json";
    configData.back_end_url +
    "/articles/" +
    params.articleId +
    "/versions/" +
    params.versionId;

  // get the optional specific history param
  const [searchParams] = useSearchParams();
  const historyID = searchParams.get("history");
  if (historyID != null) {
    url = url + "?historyID=" + historyID;
  }

  useEffect(() => {
    fetch(url, {
      method: "GET",
      mode: "cors",
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
      },
      credentials: "include",
    }).then(
      async (response) => {
        if (response.ok) {
          setError(undefined);
          let VersionData: ArticleVersion = await response.json();
          setData(VersionData);
          setLoaded(true);
          setOutdated(historyID != null && VersionData.latestHistoryID != historyID)
          setIsOwner(isLoggedInUserTheOwner() || false);
        } else {
          // Set error with message returned from the server
          let responseJSON: {
            message: string;
          } = await response.json();

          let serverMessage: string = responseJSON.message;
          setError(new Error(serverMessage));
        }
        setLoaded(true);
      },
      (error) => {
        setLoaded(true);
        setError(error);
      }
    );
  }, []);

  return (
    <div className={"row justify-content-center wrapper"}>
      {!isLoaded && <LoadingSpinner />}
      {error && (
        <NotificationAlert
          errorType="danger"
          title={"Error: "}
          message={"Something went wrong. " + error}
        />
      )}
      <div className={"col-10"}>
        <h3 className={"mt-3"}>{versionData && versionData.title}</h3>
        <div>
          <ul>
            <li className="ownersLi">Owners: </li>
            {versionData &&
              versionData.owners.map((owner, i) => (
                <li className="ownersLi" key={i}>
                  {owner + ";"}
                </li>
              ))}
          </ul>
        </div>
        <hr />
        <ul className="nav justify-content-end d-grid gap-2 d-md-flex justify-content-md-end">
          <li className="nav-item">
            <a className="nav-link">
              <Link to={"/articles/" + params.articleId + "/versions"}>
                <button
                  type="button"
                  className="btn  btn-light"
                  data-bs-toggle="modal"
                  data-bs-target="#listVersions"
                >
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    width="20"
                    height="20"
                    fill="currentColor"
                    className="bi bi-list-ul"
                    viewBox="0 0 16 16"
                  >
                    <path
                      fillRule="evenodd"
                      d="M5 11.5a.5.5 0 0 1 .5-.5h9a.5.5 0 0 1 0 1h-9a.5.5 0 0 1-.5-.5zm0-4a.5.5 0 0 1 .5-.5h9a.5.5 0 0 1 0 1h-9a.5.5 0 0 1-.5-.5zm0-4a.5.5 0 0 1 .5-.5h9a.5.5 0 0 1 0 1h-9a.5.5 0 0 1-.5-.5zm-3 1a1 1 0 1 0 0-2 1 1 0 0 0 0 2zm0 4a1 1 0 1 0 0-2 1 1 0 0 0 0 2zm0 4a1 1 0 1 0 0-2 1 1 0 0 0 0 2z"
                    />
                  </svg>
                  Other versions
                </button>
              </Link>
            </a>
          </li>

          <li className="nav-item">
            <a className="nav-link">
              <FileDownload />
            </a>
          </li>
          {!isOutdated && isLoggedIn && (
            <li className="nav-item">
              <a className="nav-link">
                <button
                  type="button"
                  className="btn btn-light"
                  data-bs-toggle="modal"
                  data-bs-target="#createNewVersion"
                >
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    width="20"
                    height="20"
                    fill="currentColor"
                    className="bi bi-intersect"
                    viewBox="0 0 16 16"
                  >
                    <path d="M0 2a2 2 0 0 1 2-2h8a2 2 0 0 1 2 2v2h2a2 2 0 0 1 2 2v8a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2v-2H2a2 2 0 0 1-2-2V2zm5 10v2a1 1 0 0 0 1 1h8a1 1 0 0 0 1-1V6a1 1 0 0 0-1-1h-2v5a2 2 0 0 1-2 2H5zm6-8V2a1 1 0 0 0-1-1H2a1 1 0 0 0-1 1v8a1 1 0 0 0 1 1h2V6a2 2 0 0 1 2-2h5z" />
                  </svg>
                  Clone this version
                </button>
                <CreateArticleVersion />
              </a>
            </li>
          )}
          {!isOutdated && isLoggedIn && (
            <li className="nav-item">
              <a className="nav-link">
                <button
                  type="button"
                  className="btn btn-light"
                  data-bs-toggle="modal"
                  data-bs-target="#createMR"
                >
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    width="20"
                    height="20"
                    fill="currentColor"
                    className="bi bi-arrow-down-left-square"
                    viewBox="0 0 16 16"
                  >
                    <path
                      fillRule="evenodd"
                      d="M15 2a1 1 0 0 0-1-1H2a1 1 0 0 0-1 1v12a1 1 0 0 0 1 1h12a1 1 0 0 0 1-1V2zM0 2a2 2 0 0 1 2-2h12a2 2 0 0 1 2 2v12a2 2 0 0 1-2 2H2a2 2 0 0 1-2-2V2zm10.096 3.146a.5.5 0 1 1 .707.708L6.707 9.95h2.768a.5.5 0 1 1 0 1H5.5a.5.5 0 0 1-.5-.5V6.475a.5.5 0 1 1 1 0v2.768l4.096-4.097z"
                    />
                  </svg>
                  Request merging this version
                </button>
                <CreateMR />
              </a>
            </li>
          )}
        </ul>

        {isOutdated && (
          <p>
            <em>
              {
                "You are currently viewing a read-only version from the history, which might be outdated. Modifications are disabled."
              }
            </em>
          </p>
        )}

        <div className="row">
          <div className="col-9">
            <ul className="nav nav-tabs" id="articleTab" role="tablist">
              <li className="nav-item" role="presentation">
                <button
                  className="nav-link active"
                  id="raw-tab"
                  data-bs-toggle="tab"
                  data-bs-target="#raw-tab-pane"
                  type="button"
                  role="tab"
                  aria-controls="raw-tab-pane"
                  aria-selected="true"
                >
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    width="20"
                    height="20"
                    fill="currentColor"
                    className="bi bi-file-earmark-code"
                    viewBox="0 0 16 16"
                  >
                    <path d="M14 4.5V14a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V2a2 2 0 0 1 2-2h5.5L14 4.5zm-3 0A1.5 1.5 0 0 1 9.5 3V1H4a1 1 0 0 0-1 1v12a1 1 0 0 0 1 1h8a1 1 0 0 0 1-1V4.5h-2z" />
                    <path d="M8.646 6.646a.5.5 0 0 1 .708 0l2 2a.5.5 0 0 1 0 .708l-2 2a.5.5 0 0 1-.708-.708L10.293 9 8.646 7.354a.5.5 0 0 1 0-.708zm-1.292 0a.5.5 0 0 0-.708 0l-2 2a.5.5 0 0 0 0 .708l2 2a.5.5 0 0 0 .708-.708L5.707 9l1.647-1.646a.5.5 0 0 0 0-.708z" />
                  </svg>
                  Raw File mode
                </button>
              </li>
              <li className="nav-item" role="presentation">
                <button
                  className="nav-link "
                  id="rendered-tab"
                  data-bs-toggle="tab"
                  data-bs-target="#rendered-tab-pane"
                  type="button"
                  role="tab"
                  aria-controls="rendered-tab-pane"
                  aria-selected="true"
                >
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    width="20"
                    height="20"
                    fill="currentColor"
                    className="bi bi-eye"
                    viewBox="0 0 16 16"
                  >
                    <path d="M16 8s-3-5.5-8-5.5S0 8 0 8s3 5.5 8 5.5S16 8 16 8zM1.173 8a13.133 13.133 0 0 1 1.66-2.043C4.12 4.668 5.88 3.5 8 3.5c2.12 0 3.879 1.168 5.168 2.457A13.133 13.133 0 0 1 14.828 8c-.058.087-.122.183-.195.288-.335.48-.83 1.12-1.465 1.755C11.879 11.332 10.119 12.5 8 12.5c-2.12 0-3.879-1.168-5.168-2.457A13.134 13.134 0 0 1 1.172 8z" />
                    <path d="M8 5.5a2.5 2.5 0 1 0 0 5 2.5 2.5 0 0 0 0-5zM4.5 8a3.5 3.5 0 1 1 7 0 3.5 3.5 0 0 1-7 0z" />
                  </svg>
                  Rendered mode
                </button>
              </li>

              {!isOutdated && isLoggedIn && (
                <li className="nav-item" role="presentation">
                  <button
                    className="nav-link"
                    id="edit-tab"
                    data-bs-toggle="tab"
                    data-bs-target="#edit-tab-pane"
                    type="button"
                    role="tab"
                    aria-controls="edit-tab-pane"
                    aria-selected="true"
                  >
                    <svg
                      xmlns="http://www.w3.org/2000/svg"
                      width="20"
                      height="20"
                      fill="currentColor"
                      className="bi bi-pencil-square"
                      viewBox="0 0 16 16"
                    >
                      <path d="M15.502 1.94a.5.5 0 0 1 0 .706L14.459 3.69l-2-2L13.502.646a.5.5 0 0 1 .707 0l1.293 1.293zm-1.75 2.456-2-2L4.939 9.21a.5.5 0 0 0-.121.196l-.805 2.414a.25.25 0 0 0 .316.316l2.414-.805a.5.5 0 0 0 .196-.12l6.813-6.814z" />
                      <path
                        fillRule="evenodd"
                        d="M1 13.5A1.5 1.5 0 0 0 2.5 15h11a1.5 1.5 0 0 0 1.5-1.5v-6a.5.5 0 0 0-1 0v6a.5.5 0 0 1-.5.5h-11a.5.5 0 0 1-.5-.5v-11a.5.5 0 0 1 .5-.5H9a.5.5 0 0 0 0-1H2.5A1.5 1.5 0 0 0 1 2.5v11z"
                      />
                    </svg>
                    Edit mode (owners only)
                  </button>
                </li>
              )}
            </ul>
            <div className="tab-content" id="articleTabContent">
              <div
                className="tab-pane show active"
                id="raw-tab-pane"
                role="tabpanel"
                aria-labelledby="raw-tab"
                tabIndex={0}
              >
                <div className="raw-article">
                  {versionData && versionData.content}
                </div>
              </div>
              <div
                className="tab-pane fade fade"
                id="rendered-tab-pane"
                role="tabpanel"
                aria-labelledby="rendered-tab"
                tabIndex={0}
              >
                Not yet implemented
              </div>
              <div
                className="tab-pane fade"
                id="edit-tab-pane"
                role="tabpanel"
                aria-labelledby="edit-tab"
                tabIndex={0}
              >
                {versionData && (
                  <InlineEditor
                    content={versionData.content}
                    setContent={setContent}
                  />
                )}
              </div>
            </div>
          </div>

          <div className="col-3">
            {(versionData && !isOutdated && (
              <ThreadList
                threadType={"commit"}
                specificId={versionData.latestHistoryID}
              />
            )) ||
              (historyID && (
                <ThreadList threadType={"commit"} specificId={historyID} />
              ))}
          </div>
        </div>
      </div>
    </div>
  );
}
