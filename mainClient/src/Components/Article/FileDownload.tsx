import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import * as React from "react";
import SuccessAlert from "../SuccessAlert";
import ErrorAlert from "../ErrorAlert";

export default function FileDownload() {
  let [error, setError] = useState(null);
  let [downloadSuccess, setDownloadSuccess] = useState<boolean>(false);
  let [filename, setFilename] = useState("");

  // Make url for request to access ../files endpoint. Debug with url "/source-file.txt"
  let params = useParams();

  const endpointUrl =
    "http://localhost:8080/articles/" +
    params.articleId +
    "/versions/" +
    params.versionId +
    "/files";

  // Handler triggered on Download button click
  const downloadFileHandler = (e: { preventDefault: () => void }) => {
    e.preventDefault();
    // GET request to get the files as a ZIP
    fetch(endpointUrl, {
      method: "GET",
      mode: "cors",
      headers: {
        "Content-Type": "application/x-zip-compressed",
        Accept: "application/x-zip-compressed",
      },
      credentials: "include",
    })
      .then((res) =>
          res.blob().then((blob) => {
              const filename = res.headers.get("content-disposition")!.split('"')[1]
              let tuple: [Blob, string] = [blob, filename];
              return tuple
          })
      )
      .then(
        // Process the response as a BLOB (Binary large object)
        (result) => {
            let filename = result[1]
            setFilename(filename)

          // Put the file in the DOM
          const windowUrl = window.URL.createObjectURL(result[0]);
          // Set filename
          // setFilename("source-file.zip");

          // Add a hidden <a> element to DOM, that downloads the file when clicking on it
          const a = document.createElement("a");
          a.style.display = "none";
          a.href = windowUrl;
          a.download = filename;
          document.body.appendChild(a);

          // Simulate clicking on <a> element to trigger download of file
          a.click();

          // Remove <a> from DOM
          window.URL.revokeObjectURL(windowUrl);

          // Set success in state to show success alert
          setDownloadSuccess(true);

          // After 3s, remove success from state to hide success alert
          setTimeout(() => setDownloadSuccess(false), 3000);
        },
        (error) => {
          setError(error);
        }
      );
  };

  return (
    <div>
      <form>
        <button className="btn btn-primary" onClick={downloadFileHandler}>
          Download source files
        </button>
      </form>
      {downloadSuccess && (
        <SuccessAlert
          title="Download successful! "
          message={
            "The source file " + filename + " has been successfully downloaded."
          }
        />
      )}
      {error && (
        <ErrorAlert
          title={"Error: "}
          message={"Something went wrong, please try again."}
        />
      )}
    </div>
  );
}
