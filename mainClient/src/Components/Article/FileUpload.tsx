import * as React from "react";
import { useState } from "react";
import { useParams } from "react-router-dom";
import configData from "../../config.json";

export default function FileUpload() {
  let [selectedFile, setSelectedFile] = useState<string | Blob>("file");
  let [error, setError] = useState(null);

  const onChangeFile = (e: any) => {
    setSelectedFile(e.target.files[0]);
  };

  let params = useParams();
  const url =
    configData.back_end_url +"/articles/" +
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
        credentials: 'include',
      body: formData,
    }).then(
      (response) => {
        let message: string =
          response.status === 200
            ? "File successfully uploaded"
            : "Error: " + response.status + " " + response.statusText;
        console.log(message);

        // refresh page
        window.location.reload();
      },
      (error) => {
        console.error("Error: ", error);
        setError(error);
      }
    );
  };

  return (
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
            {error && <div>{`Error: ${error}`}</div>}

            <h5>Upload a file</h5>
            <hr />
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
  );
}
