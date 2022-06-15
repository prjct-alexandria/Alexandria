import * as React from "react";
import { useParams } from "react-router-dom";
import { useEffect, useState } from "react";
import PrismDiff from "./PrismDiff";
import LoadingSpinner from "../LoadingSpinner";
import ThreadList from "./ThreadList";
import NotificationAlert from "../NotificationAlert";
import isUserLoggedIn from "../User/AuthHelpers/isUserLoggedIn";

type Request = {
  sourceVersionID: number;
  sourceHistoryID: number;
  targetVersionID: number;
  targetHistoryID: number;
  status: string;
};

type ArticleVersion = {
  id: number;
  title: string;
  owners: string[];
  content: string;
};

function getRequest(
  url: string,
  setData: (r: Request) => void,
  setLoaded: (b: boolean) => void,
  setError: (e: Error | undefined) => void
) {
  fetch(url).then(
    async (response) => {
      if (response.ok) {
        setError(undefined);
        let requestData: Request = await response.json();
        setData(requestData);
        setLoaded(true);
      } else {
        // Set error with message returned from the server
        let responseJSON: {
          message: string;
        } = await response.json();

        let serverMessage: string = responseJSON.message;
        setError(new Error(serverMessage));
      }
    },
    (error) => {
      setError(error.message);
      setLoaded(true);
    }
  );
}

function getVersion(
  url: string,
  setData: (v: ArticleVersion) => void,
  setLoaded: (b: boolean) => void,
  setError: (e: Error | undefined) => void
) {
  fetch(url).then(
    async (response) => {
      if (response.ok) {
        setError(undefined);
        let articleSourceData: ArticleVersion = await response.json();
        setData(articleSourceData);
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
      setError(error.message);
      setLoaded(true);
    }
  );
}
export default function VersionList() {
  let [isLoggedIn, setLoggedIn] = useState<boolean>(isUserLoggedIn());

  // Listen for userAccountEvent that fires when user in localstorage changes
  window.addEventListener("userAccountEvent", () => {
    setLoggedIn(isUserLoggedIn());
  });

  let params = useParams();

  const urlRequest = "/request.json";
  //const urlRequest = 'http://localhost:8080/articles/' + params.articleId  + "/requests/" + params.requestId;

  let [dataRequest, setDataRequest] = useState<Request>();
  let [isLoadedRequest, setLoadedRequest] = useState<boolean>(false);
  let [errorRequest, setErrorRequest] = useState<Error>();

  // fetching information about the request: historyID of source version, versionID of target, historyID of target, state of the request
  useEffect(() => {
    getRequest(urlRequest, setDataRequest, setLoadedRequest, setErrorRequest);
  }, []);

  let urlArticleSource = "";
  let urlArticleTarget = "";

  if (dataRequest !== undefined) {
    urlArticleSource =
      "http://localhost:8080/articles/" +
      params.articleId +
      "/versions/" +
      params.versionId +
      "/history/" +
      dataRequest.sourceHistoryID;
    urlArticleTarget =
      "http://localhost:8080/articles/" +
      params.articleId +
      "/versions/" +
      dataRequest.targetVersionID +
      "/history/" +
      dataRequest.targetHistoryID;
  }
  urlArticleSource = "/article_version1.json"; // Placeholder source version
  urlArticleTarget = "/article_version2.json"; // Placeholder target version

  let [dataSource, setDataSource] = useState<ArticleVersion>();
  let [isLoadedSource, setLoadedSource] = useState<boolean>(false);
  let [errorSource, setErrorSource] = useState<Error>();

  let [dataTarget, setDataTarget] = useState<ArticleVersion>();
  let [isLoadedTarget, setLoadedTarget] = useState<boolean>(false);
  let [errorTarget, setErrorTarget] = useState<Error>();

  // fetching the actual articles
  useEffect(() => {
    getVersion(
      urlArticleSource,
      setDataSource,
      setLoadedSource,
      setErrorSource
    );
    getVersion(
      urlArticleTarget,
      setDataTarget,
      setLoadedTarget,
      setErrorTarget
    );
  }, []);

  // Disable button if it is already accepted or rejected
  const disableButton = () => {
    return (
      dataRequest !== undefined &&
      (dataRequest.status === "accepted" || dataRequest.status === "rejected")
    );
  };

  // If already accepted, fill in the color of the button and disable. If request is pending, send HTTP request.
  const acceptButton = () => {
    let className = "btn btn-outline-success";
    if (dataRequest !== undefined && dataRequest.status === "accepted") {
      className = "btn btn-success";
    }
    return (
      <button
        className={className}
        disabled={disableButton()}
        onClick={handleAcceptClick}
      >
        Accept
      </button>
    );
  };

  // If already reject, fill in the color of the button and disable. If request is pending, send HTTP request.
  const rejectButton = () => {
    let className = "btn btn-outline-danger";
    if (dataRequest !== undefined && dataRequest.status === "rejected") {
      className = "btn btn-danger";
    }
    return (
      <button
        className={className}
        disabled={disableButton()}
        onClick={handleRejectClick}
      >
        Reject
      </button>
    );
  };

  // If already accepted or rejected, do not show the button. If pending, send HTTP request.
  const deleteButton = () => {
    if (!disableButton()) {
      return (
        <button className={"btn btn-danger"} onClick={handleDeleteClick}>
          Delete
        </button>
      );
    }
  };

  let [acceptSuccess, setAcceptSuccess] = useState<boolean>(false);
  let [rejectSuccess, setRejectSuccess] = useState<boolean>(false);
  let [deleteSuccess, setDeleteSuccess] = useState<boolean>(false);
  let [error, setError] = useState<Error>();

  // Send HTTP request and reload, so that the style (see "acceptButton") is updated.
  const handleAcceptClick = async () => {
    const url =
      "http://localhost:8080/articles/" +
      params.articleId +
      "/requests/" +
      params.requestId +
      "/accept";
    fetch(url, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      mode: "cors",
      credentials: "include",
    }).then(
      async (response) => {
        if (response.ok) {
          setError(undefined);
          // Set success in state to show success alert
          setAcceptSuccess(true);

          // After 3s, remove success from state to hide success alert
          setTimeout(() => setAcceptSuccess(false), 3000);
        } else {
          let responseJSON: {
            message: string;
          } = await response.json();

          let serverMessage: string = responseJSON.message;
          setError(new Error(serverMessage));
        }
      },
      (error) => {
        setError(error);
      }
    );
    window.location.reload();
  };

  // Send HTTP request and reload, so that the style (see "rejectButton") is updated.
  const handleRejectClick = async () => {
    const url =
      "http://localhost:8080/articles/" +
      params.articleId +
      "/requests/" +
      params.requestId +
      "/reject";
    fetch(url, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      mode: "cors",
      credentials: "include",
    }).then(
      async (response) => {
        if (response.ok) {
          setError(undefined);
          // Set success in state to show success alert
          setRejectSuccess(true);

          // After 3s, remove success from state to hide success alert
          setTimeout(() => setRejectSuccess(false), 3000);
        } else {
          let responseJSON: {
            message: string;
          } = await response.json();

          let serverMessage: string = responseJSON.message;
          setError(new Error(serverMessage));
        }
      },
      (error) => {
        setError(error);
      }
    );
    window.location.reload();
  };

  // Send HTTP request and reload, so that the style (see "deleteButton") is updated.
  const handleDeleteClick = async () => {
    const url =
      "http://localhost:8080/articles/" +
      params.articleId +
      "/requests/" +
      params.requestId;
    fetch(url, {
      method: "DELETE",
      headers: { "Content-Type": "application/json" },
      mode: "cors",
      credentials: "include",
    }).then(
      async (response) => {
        if (response.ok) {
          setError(undefined);
          // Set success in state to show success alert
          setDeleteSuccess(true);

          // After 3s, remove success from state to hide success alert
          setTimeout(() => setDeleteSuccess(false), 3000);
        } else {
          let responseJSON: {
            message: string;
          } = await response.json();

          let serverMessage: string = responseJSON.message;
          setError(new Error(serverMessage));
        }
      },
      (error) => {
        setError(error);
      }
    );
  };

  const view = () => {
    return (
      <div className="row">
        <div>
          <h1 style={{ textAlign: "center", marginBottom: "30px" }}>
            Compare Versions
          </h1>
          <div className="row justify-content-center">
            {/*Version names*/}
            <div className="row col-8 mb-2">
              <div className="col-6">
                <h5>
                  Changes of '{dataTarget !== undefined && dataTarget.title}'
                </h5>
              </div>
              <div className="col-6">
                <h5>Result: {dataSource !== undefined && dataSource.title}</h5>
              </div>
            </div>

            {/*Accept, reject and delete button*/}
            <div className="col-1" id="AcceptButton">
              {acceptButton()}
            </div>
            <div className="col-1">{rejectButton()}</div>
            <div className="col-1">{deleteButton()}</div>
          </div>

          <div className="row justify-content-center">
            {/*Content of versions*/}
            <div className="wrapper col-8">
              <div
                className="row overflow-scroll"
                style={{
                  height: "500px",
                  whiteSpace: "pre-line",
                  border: "grey solid 3px",
                }}
              >
                {/*Source version, including changes that are made*/}
                <div className="col-6" style={{ whiteSpace: "pre-line" }}>
                  {dataSource !== undefined && dataTarget !== undefined && (
                    <PrismDiff
                      sourceContent={dataSource.content}
                      targetContent={dataTarget.content}
                    />
                  )}
                </div>
                {/*Target version*/}
                <div className="col-6" style={{ whiteSpace: "pre-line" }}>
                  {dataTarget !== undefined && dataTarget.content}
                </div>
              </div>
            </div>
            <div className="wrapper col-3">
              <ThreadList
                threadType={"request"}
                specificId={parseInt(params.requestId as string)}
              />
            </div>
          </div>
        </div>
      </div>
    );
  };

  return (
    <div>
      {error && (
        <NotificationAlert
          errorType="danger"
          title={"Error: "}
          message={"Something went wrong. " + error.message}
        />
      )}
      {acceptSuccess && (
        <NotificationAlert
          errorType="success"
          title="Merging approved! "
          message={"The merge request has been successfully accepted."}
        />
      )}
      {rejectSuccess && (
        <NotificationAlert
          errorType="success"
          title="Merging rejected! "
          message={"The merge request has been successfully rejected."}
        />
      )}
      {deleteSuccess && (
        <NotificationAlert
          errorType="success"
          title="Merge request deleted! "
          message={"The merge request has been successfully deleted."}
        />
      )}
      {!isLoadedRequest && !isLoadedSource && !isLoadedTarget && (
        <LoadingSpinner />
      )}
      {errorRequest && (
        <NotificationAlert
          errorType="danger"
          title={"Error: "}
          message={"Something went wrong. " + errorRequest}
        />
      )}
      {errorSource && (
        <NotificationAlert
          errorType="danger"
          title={"Error: "}
          message={"Something went wrong. " + errorSource}
        />
      )}
      {errorTarget && (
        <NotificationAlert
          errorType="danger"
          title={"Error: "}
          message={"Something went wrong. " + errorTarget}
        />
      )}
      {dataRequest && dataSource && dataTarget && view()}
    </div>
  );
}
