import * as React from "react";
import { useState } from "react";
import { useParams } from "react-router-dom";

export default function FileUpload() {
  let [selectedFile, setSelectedFile] = useState<string | Blob>("file");
  let [error, setError] = useState(null);
  let [httpResponse, setHttpResponse] = useState<Response>();

  const onChangeFile = (e: any) => {
    setSelectedFile(e.target.files[0]);
  };

  let params = useParams();
  const url =
    "http://localhost:8080/articles/" +
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
        "Content-Type": "application/json",
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
        setHttpResponse(response);
      },
      (error) => {
        console.error("Error: ", error);
        setError(error);
      }
    );
  };

  return (
    <div className="col-6 align-content-center m-auto">
      {error && <div>{`Error: ${error}`}</div>}

      <h3>Upload a file</h3>
      <hr />
      <input
        className="file-upload-input"
        type="file"
        onChange={onChangeFile}
      />
      <button className="btn btn-primary" onClick={uploadFileHandler}>
        Upload
      </button>
    </div>
  );
}
