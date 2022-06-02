import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import * as React from "react";

export default function FileDownload() {
  let [isLoaded, setLoaded] = useState(false);
  let [error, setError] = useState(null);
  let [fileZip, setFileZip] = useState<Blob | MediaSource>();

  // Handler triggered on Download button click
  const downloadFileHandler = () => {
    // Make url for request to access ../files endpoint
    let params = useParams();
    const url =
      "http://localhost:8080/articles/" +
      params.articleId +
      "/versions/" +
      params.versionId +
      "/files";

    // GET request to get the files as a ZIP
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
        // Process the response as a BLOB (Binary large object)
        (result) => {
          setLoaded(true);
          setFileZip(result);

          const url = window.URL.createObjectURL(result);
          const a = document.createElement("a");
          a.style.display = "none";
          a.href = url;

          // the filename
          a.download = "source.zip";
          document.body.appendChild(a);
          a.click();
          window.URL.revokeObjectURL(url);
          alert("your file has downloaded!"); // replace by something with better UX
        },
        (error) => {
          setLoaded(true);
          setError(error);
        }
      );
  };

  return (
    <button className="btn btn-primary" onClick={downloadFileHandler}>
      Download source files
    </button>
  );
}
