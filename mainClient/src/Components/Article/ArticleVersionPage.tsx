import * as React from "react";
import { useEffect, useState } from "react";
import {useParams, useSearchParams} from "react-router-dom";
import LoadingSpinner from "../LoadingSpinner";
import FileUpload from "./FileUpload";
import CreateMR from "./CreateMR"
import ThreadList from "./ThreadList";
import CreateArticleVersion from "./CreateArticleVersion";
import FileDownload from "./FileDownload";

type ArticleVersion = {
  owners: Array<string>;
  title: string;
  content: string;
};

export default function ArticleVersionPage() {
  let [versionData, setData] = useState<ArticleVersion>();
  let [isLoaded, setLoaded] = useState(false);
  let [error, setError] = useState(null);

  let params = useParams();

  let url = //"/article_version1.json";
  "http://localhost:8080/articles/" +
    params.articleId +
    "/versions/" +
    params.versionId;

  // get the optional specific history param
  const [searchParams] = useSearchParams(); // used for the source and target
  let historyID = searchParams.get('history');
  const viewingOldVersion = historyID != null;
  if (viewingOldVersion) {
    url = url + '?historyID=' + historyID
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
    })
      .then((res) => res.json())
      .then(
        (result) => {
          setLoaded(true);
          setData(result);
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
        {error && <div>{`There is a problem fetching the data - ${error}`}</div>}
        {versionData != null && (
            <div className={"col-10"}>
              <h1>{versionData.title}</h1>
              <div>
                <h4>Owners:</h4>
                <ul>
                  {versionData.owners.map((owner, i) => (
                      <li key={i}>{owner}</li>
                  ))}
                </ul>
              </div>
              <div className={"row"}>
                {!viewingOldVersion && (
                    <div className="col-2">
                      <button
                          type="button"
                          className="btn btn-primary btn-lg"
                          data-bs-toggle="modal"
                          data-bs-target="#uploadFile"
                      >
                        Upload File
                      </button>
                      <FileUpload />
                    </div>
                )}
                <div className="col-2">
                  <FileDownload />
                </div>
                {!viewingOldVersion && (
                    <div className="col-2">
                      <button
                          type="button"
                          className="btn btn-primary btn-lg"
                          data-bs-toggle="modal"
                          data-bs-target="#createMR"
                      >
                        Make Request
                      </button>
                      <CreateMR />
                    </div>
                )}
                {!viewingOldVersion && (
                    <div className="col-2">
                      <button
                          type="button"
                          className="btn btn-primary btn-lg"
                          data-bs-toggle="modal"
                          data-bs-target="#createNewVersion"
                      >
                        Clone this version
                      </button>
                      <CreateArticleVersion />
                    </div>
                )}
              </div>
              {viewingOldVersion&&
                  <p><em>{"You are currently viewing a read-only version from the history, which might be outdated. Modifications are disabled."}</em></p>
              }
              <div className="row">
                <div className="row mb-2 mt-2">
                  <div className="col-8 articleContent">
                    <div style={{whiteSpace: "pre-line"}}>{versionData.content}</div>
                  </div>
                  <div className="col-3">
                    <ThreadList threadType={"commit"} specificId={parseInt(params.versionId as string)} />
                  </div>
                </div>
              </div>
            </div>
        )}
      </div>
  );
}
