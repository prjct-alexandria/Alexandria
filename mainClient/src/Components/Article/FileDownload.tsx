import { useState } from "react";
import { useParams } from "react-router-dom";
import * as React from "react";
import NotificationAlert from "../NotificationAlert";
import backEndUrl from "../../urlUtils";

export default function FileDownload() {
  let [error, setError] = useState<Error>();
  let [downloadSuccess, setDownloadSuccess] = useState<boolean>(false);
  let [filename, setFilename] = useState<string>("");

  // Make url for request to access ../files endpoint. Debug with url "/source-file.txt"
  let params = useParams();

  const endpointUrl =
    backEndUrl() +
    "/articles/" +
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
    }).then(
      // Process the response as a BLOB (Binary large object)
      async (response) => {
        if (response.ok) {
          let blob = await response.blob();
          // Put the file in the DOM
          const windowUrl = window.URL.createObjectURL(blob);
          // Set filename
          let filename = response.headers
            .get("content-disposition")!
            .split('"')[1];
          setFilename(filename);

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

          // Set success in state to show success alert for 3 seconds
          setDownloadSuccess(true);
          setTimeout(() => setDownloadSuccess(false), 3000);
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
    <div>
      <form>
        <button className="btn  btn-light" onClick={downloadFileHandler}>
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="16"
            height="16"
            fill="currentColor"
            className="bi bi-download"
            viewBox="0 0 16 16"
          >
            <path d="M.5 9.9a.5.5 0 0 1 .5.5v2.5a1 1 0 0 0 1 1h12a1 1 0 0 0 1-1v-2.5a.5.5 0 0 1 1 0v2.5a2 2 0 0 1-2 2H2a2 2 0 0 1-2-2v-2.5a.5.5 0 0 1 .5-.5z" />
            <path d="M7.646 11.854a.5.5 0 0 0 .708 0l3-3a.5.5 0 0 0-.708-.708L8.5 10.293V1.5a.5.5 0 0 0-1 0v8.793L5.354 8.146a.5.5 0 1 0-.708.708l3 3z" />
          </svg>
          Download source files
        </button>
      </form>
      {downloadSuccess && (
        <NotificationAlert
          errorType="success"
          title="Download successful! "
          message={
            "The source file " + filename + " has been successfully downloaded."
          }
        />
      )}
      {error && (
        <NotificationAlert
          errorType="danger"
          title={"Error: "}
          message={"Something went wrong, please try again."}
        />
      )}
    </div>
  );
}
