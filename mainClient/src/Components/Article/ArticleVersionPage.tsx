import * as React from "react";
import { useEffect, useState } from "react";
import { useParams, useSearchParams } from "react-router-dom";
import LoadingSpinner from "../LoadingSpinner";
import FileUpload from "./FileUpload";
import CreateMR from "./CreateMR";
import ThreadList from "./ThreadList";
import CreateArticleVersion from "./CreateArticleVersion";
import FileDownload from "./FileDownload";
import NotificationAlert from "../NotificationAlert";
import isUserLoggedIn from "../User/AuthHelpers/isUserLoggedIn";

type ArticleVersion = {
  owners: Array<string>;
  title: string;
  content: string;
};

export default function ArticleVersionPage() {
  let [versionData, setData] = useState<ArticleVersion>();
  let [isLoaded, setLoaded] = useState<boolean>(false);
  let [error, setError] = useState<Error>();
  let [isLoggedIn, setLoggedIn] = useState<boolean>(isUserLoggedIn());

  // Listen for userAccountEvent that fires when user in localstorage changes
  window.addEventListener("userAccountEvent", () => {
    setLoggedIn(isUserLoggedIn());
  });

  let params = useParams();

  let url =
    // "/article_version1.json"; // Placeholder
    "http://localhost:8080/articles/" +
    params.articleId +
    "/versions/" +
    params.versionId;

  // get the optional specific history param
  const [searchParams] = useSearchParams(); // used for the source and target
  let historyID = searchParams.get("history");
  const viewingOldVersion = historyID != null;
  if (viewingOldVersion) {
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
        } else {
          setLoaded(true);
          // Set error with message returned from the server
          let responseJSON: {
            message: string;
          } = await response.json();

          let serverMessage: string = responseJSON.message;
          setError(new Error(serverMessage));
        }
      },
      (error) => {
        setLoaded(true);
        setError(error);
      }
    );
  }, [url]);

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
          {!viewingOldVersion && isLoggedIn && (
            <li className="nav-item">
              <a className="nav-link">
                <button
                  type="button"
                  className="btn  btn-light"
                  data-bs-toggle="modal"
                  data-bs-target="#uploadFile"
                >
                  Upload File
                </button>
                <FileUpload />
              </a>
            </li>
          )}
          <li className="nav-item">
            <a className="nav-link">
              <FileDownload />
            </a>
          </li>
          {!viewingOldVersion && (
            <li className="nav-item">
              <a className="nav-link">
                <button
                  type="button"
                  className="btn btn-light"
                  data-bs-toggle="modal"
                  data-bs-target="#createNewVersion"
                >
                  Clone this version
                </button>
                <CreateArticleVersion />
              </a>
            </li>
          )}
          {!viewingOldVersion && isLoggedIn && (
            <li className="nav-item">
              <a className="nav-link">
                <button
                  type="button"
                  className="btn btn-light"
                  data-bs-toggle="modal"
                  data-bs-target="#createMR"
                >
                  Make Request
                </button>
                <CreateMR />
              </a>
            </li>
          )}
        </ul>

        {viewingOldVersion && (
          <p>
            <em>
              {
                "You are currently viewing a read-only version from the history, which might be outdated. Modifications are disabled."
              }
            </em>
          </p>
        )}

        <div className="row">
          <div className="row mb-2 mt-2">
            <div className="col-8 articleContent">
              <div style={{ whiteSpace: "pre-line" }}>
                {versionData && versionData.content}
              </div>
            </div>
            <div className="col-3">
              <ThreadList
                threadType={"commit"}
                specificId={parseInt(params.versionId as string)}
              />
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
