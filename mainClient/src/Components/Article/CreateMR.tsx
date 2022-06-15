import * as React from "react";
import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import LoadingSpinner from "../LoadingSpinner";
import NotificationAlert from "../NotificationAlert";

type Version = {
  articleID: string;
  versionID: string;
  title: string;
  owners: string[];
  status: string;
};

export default function CreateMR() {
  let params = useParams();

  let [dataCurVersion, setDataCurVersion] = useState<Version>();
  let [isLoaded, setLoaded] = useState<boolean>(false);
  let [error, setError] = useState<Error>();
  let [createMRSuccess, setCreateMRSuccess] = useState<boolean>(false);

  let [dataVersions, setDataVersions] = useState<Version[]>();

  const [selectedVersion, setSelectedVersion] = useState("");

  // const urlCurrentArticle = "http://localhost:8080/articles/" + params.articleId + "/versions/" + params.versionId;
  const urlCurrentArticle = "/article_version1.json";
  // retrieving the list of versions
  useEffect(() => {
    fetch(urlCurrentArticle, {
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
          let versionData: Version = await response.json();
          setDataCurVersion(versionData);
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
  }, [urlCurrentArticle]);

  // const urlArticleVersions = "http://localhost:8080/articles/" + params.articleId + "/versions";
  const urlArticleVersions = "/versionList.json";
  useEffect(() => {
    fetch(urlArticleVersions, {
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
          let versionList: Version[] = await response.json();
          setDataVersions(versionList);
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
  }, [urlArticleVersions]);

  function handleSelectChange(event: {
    target: { value: React.SetStateAction<string> };
  }) {
    setSelectedVersion(event.target.value);
  }

  const urlSubmitMR =
    "http://localhost:8080/articles/" + params.articleId + "/requests";
  // post the new merge request
  const submitMR = () => {
    fetch(urlSubmitMR, {
      method: "POST",
      mode: "cors",
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
      },
      credentials: "include",
      body: JSON.stringify({
        sourceVersionID: parseInt(params.versionId as string),
        targetVersionID: parseInt(selectedVersion as string),
      }),
    }).then(
      (response) => {
        if (response.status === 200) {
          // Set success in state to show success alert
          setCreateMRSuccess(true);

          // After 3s, remove success from state to hide success alert
          setTimeout(() => setCreateMRSuccess(false), 3000);
        }
      },
      (error) => {
        console.error("Error: ", error);
        setError(error);
      }
    );
  };

  return (
    <>
      {createMRSuccess && (
        <NotificationAlert
          errorType="success"
          title="Success! "
          message={"Merge request successfully created."}
        />
      )}
      {error && (
        <NotificationAlert
          errorType="danger"
          title={"Error: "}
          message={"Something went wrong. " + error}
        />
      )}
      <div
        className="modal fade create-article-modal"
        id="createMR"
        data-bs-backdrop="static"
        data-bs-keyboard="false"
        aria-labelledby="createMRLabel"
        aria-hidden="true"
      >
        {!isLoaded && <LoadingSpinner />}
        {dataCurVersion && dataVersions && (
          <div className="modal-dialog">
            <div className="modal-content">
              <div className="modal-header">
                <h5 className="modal-title" id="createMRLabel">
                  Request to apply changes:
                </h5>
              </div>
              <div className="modal-body">
                <div className="row justify-content-center">
                  <div className="col-6">
                    <h6>From</h6>
                  </div>
                  <div className="col-6">
                    <h6>To</h6>
                  </div>
                  <div className="col-6">{dataCurVersion.title}</div>
                  <div className="col-6">
                    {/*list of other versions*/}
                    <select
                      className="form-select"
                      aria-label="Default select example"
                      value={selectedVersion}
                      onChange={handleSelectChange}
                    >
                      <option value="" disabled={true}></option>
                      {dataVersions.map((version, i) =>
                        dataCurVersion !== undefined &&
                        dataCurVersion.versionID !== version.versionID ? ( // do not show version if it is the same as the source version
                          <option value={version.versionID} key={i}>
                            {version.title}
                          </option>
                        ) : null
                      )}
                    </select>
                  </div>
                </div>
              </div>
              <div className="modal-footer">
                <button
                  type="button"
                  className="btn btn-secondary"
                  data-bs-dismiss="modal"
                >
                  Close
                </button>
                <button
                  type="button"
                  className="btn btn-primary"
                  onClick={submitMR}
                  disabled={selectedVersion == ""}
                >
                  Create request
                </button>
              </div>
            </div>
          </div>
        )}
      </div>
    </>
  );
}
