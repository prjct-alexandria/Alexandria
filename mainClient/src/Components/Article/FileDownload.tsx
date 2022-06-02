import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import * as React from "react";

export default function FileDownload() {
  let [isLoaded, setLoaded] = useState(false);
  let [error, setError] = useState(null);
  let [fileZip, setFileZip] = useState<string | Blob>("file");

  let params = useParams();
  const url =
    "http://localhost:8080/articles/" +
    params.articleId +
    "/versions/" +
    params.versionId +
    "/files";

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
          setFileZip(result);
        },
        (error) => {
          setLoaded(true);
          setError(error);
        }
      );
  }, [url]);

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
              Download source files
            </h5>

            <button
              type="button"
              className="btn-close"
              data-bs-dismiss="modal"
              aria-label="Close"
            ></button>
          </div>
          <div className="modal-body"></div>
          <div className="modal-footer"></div>
        </div>
      </div>
    </div>
  );
}
