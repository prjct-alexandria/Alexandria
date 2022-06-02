import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import * as React from "react";

export default function FileDownload() {
  let [isLoaded, setLoaded] = useState(false);
  let [error, setError] = useState(null);
  let [fileZip, setFileZip] = useState<Blob | MediaSource>();

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
      .then((res) => res.blob())
      .then(
        (result) => {
          setLoaded(true);
          setFileZip(result);

          const url = window.URL.createObjectURL(result);
          const a = document.createElement("a");
          a.style.display = "none";
          a.href = url;
          // the filename you want
          a.download = "todo-1.json";
          document.body.appendChild(a);
          a.click();
          window.URL.revokeObjectURL(url);
          alert("your file has downloaded!"); // or you know, something with better UX...
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
      aria-labelledby="downloadLabel"
      aria-hidden="true"
    >
      <div className="modal-dialog">
        <div className="modal-content">
          <div className="modal-header">
            <h5 className="modal-title" id="downloadLabel">
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
