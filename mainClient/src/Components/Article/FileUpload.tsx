import * as React from "react";
import { useState } from "react";
import { useParams } from "react-router-dom";
import NotificationAlert from "../NotificationAlert";
import backEndUrl from "../../urlUtils";

export default function FileUpload() {
  let [selectedFile, setSelectedFile] = useState<string | Blob>("file");
  let [error, setError] = useState<Error>();
  let [uploadSuccess, setUploadSuccess] = useState<boolean>(false);

  const onChangeFile = (e: any) => {
    setSelectedFile(e.target.files[0]);
  };

  let params = useParams();
  const url =
    backEndUrl() +"/articles/" +
    params.articleId +
    "/versions/" +
    params.versionId;

  // Send an HTTP POST request with file to /articles/#id/versions/#id
  const uploadFileHandler = () => {
    const formData = new FormData();
    formData.append("file", selectedFile);

    fetch(url, {
      method: "POST",
      mode: "cors",
      headers: {
        Accept: "application/json",
      },
      credentials: "include",
      body: formData,
    }).then(
      async (response) => {
        // refresh page
        window.location.reload();
        if (response.ok) {
          // Set success in state to show success alert for 3 seconds
          setUploadSuccess(true);
          setTimeout(() => setUploadSuccess(false), 3000);
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
        setError(error);
      }
    );
  };

  return (
    <>
      {uploadSuccess && (
        <NotificationAlert
          errorType="success"
          title="Upload successful! "
          message={"The file has been successfully uploaded."}
        />
      )}
      {error && (
        <NotificationAlert
          errorType="danger"
          title={"Error: "}
          message={"Something went wrong, please try again."}
        />
      )}
      <div
        className="modal fade"
        id="uploadFile"
        data-bs-backdrop="static"
        data-bs-keyboard="false"
        aria-labelledby="uploadLabel"
        aria-hidden="true"
      >
        <div className="modal-dialog">
          <div className="modal-content">
            <div className="modal-header">
              <h5 className="modal-title" id="uploadLabel">
                Upload file to article
              </h5>

              <button
                type="button"
                className="btn-close"
                data-bs-dismiss="modal"
                aria-label="Close"
              ></button>
            </div>
            <div className="modal-body">
              {error && (
                <NotificationAlert
                  errorType="danger"
                  title={"Error: "}
                  message={"Something went wrong. " + error}
                />
              )}

              <input
                className="file-upload-input"
                type="file"
                onChange={onChangeFile}
              />
            </div>
            <div className="modal-footer">
              <button className="btn btn-primary" onClick={uploadFileHandler}>
                Upload
              </button>
            </div>
          </div>
        </div>
      </div>
    </>
  );
}
