import { Dispatch, useState } from "react";
import { useParams } from "react-router-dom";
import * as React from "react";
import FileUpload from "./FileUpload";
import NotificationAlert from "../NotificationAlert";

type EditorProps = {
  content: string;
  setContent: Dispatch<React.SetStateAction<string>>;
};

export default function InlineEditor(props: EditorProps) {
  let [editorData, setEditorData] = useState<string>(props.content);
  let [showSaveSuccessMessage, setSaveSuccess] = useState<boolean>(false);
  let [editError, setEditError] = useState<Error>();

  let params = useParams();

  const url =
    "http://localhost:8080/articles/" +
    params.articleId +
    "/versions/" +
    params.versionId;

  const saveChanges = () => {
    let editedFile: Blob = new Blob([editorData], { type: "application/json" });
    let formData = new FormData();
    formData.append("file", editedFile);

    fetch(url, {
      method: "POST",
      mode: "cors",
      credentials: "include",
      body: formData,
    })
      .then(
        async (response) => {
          if (response.ok) {
            // props.content = editorData;
            setSaveSuccess(true);
            setTimeout(() => setSaveSuccess(false), 3000);
            // dispatch
            props.setContent(editorData);
            return true;
          } else {
            // Set error with message returned from the server
            let responseJSON: {
              message: string;
            } = await response.json();

            let serverMessage: string = responseJSON.message;
            setEditError(new Error(serverMessage));
            return false;
          }
        },
        (error) => {
          setEditError(error);
          return false;
        }
      )
      .then((wasSuccessful) => {
        if (wasSuccessful) {
          window.dispatchEvent(new Event("changesSavedEvent"));
        }
      });
  };

  const discardChanges = () => {
    setEditorData(props.content);
  };

  // @ts-ignore
  return (
    <>
      {showSaveSuccessMessage && (
        <NotificationAlert
          errorType="success"
          title="Edits saved! "
          message={"Your changes have been successfully saved."}
        />
      )}
      {editError && (
        <NotificationAlert
          errorType="danger"
          title={"Error: "}
          message={"Something went wrong. " + editError}
        />
      )}
      <div className="inline-editor-article">
        <textarea
          className="col-xs-12 textarea inline-editor-textarea"
          value={editorData}
          onChange={(e) => {
            setEditorData(e.target.value);
          }}
        />
      </div>
      <nav className="nav d-grid gap-2 d-md-flex">
        <>
          <button
            type="button"
            className="btn  btn-light"
            data-bs-toggle="modal"
            data-bs-target="#uploadFile"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              width="20"
              height="20"
              fill="currentColor"
              className="bi bi-file-earmark-arrow-up"
              viewBox="0 0 16 16"
            >
              <path d="M8.5 11.5a.5.5 0 0 1-1 0V7.707L6.354 8.854a.5.5 0 1 1-.708-.708l2-2a.5.5 0 0 1 .708 0l2 2a.5.5 0 0 1-.708.708L8.5 7.707V11.5z" />
              <path d="M14 14V4.5L9.5 0H4a2 2 0 0 0-2 2v12a2 2 0 0 0 2 2h8a2 2 0 0 0 2-2zM9.5 3A1.5 1.5 0 0 0 11 4.5h2V14a1 1 0 0 1-1 1H4a1 1 0 0 1-1-1V2a1 1 0 0 1 1-1h5.5v2z" />
            </svg>
            Upload File (overwrites current)
          </button>
          <FileUpload />
        </>
        <div className="flex-fill"></div>
        <button
          className={"nav-item btn btn-success "}
          type="button"
          onClick={saveChanges}
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="16"
            height="16"
            fill="currentColor"
            className="bi bi-check-lg"
            viewBox="0 0 16 16"
          >
            <path d="M12.736 3.97a.733.733 0 0 1 1.047 0c.286.289.29.756.01 1.05L7.88 12.01a.733.733 0 0 1-1.065.02L3.217 8.384a.757.757 0 0 1 0-1.06.733.733 0 0 1 1.047 0l3.052 3.093 5.4-6.425a.247.247 0 0 1 .02-.022Z" />
          </svg>
          Save changes
        </button>
        <button
          className={"nav-item btn btn-light"}
          type="button"
          onClick={discardChanges}
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="16"
            height="16"
            fill="currentColor"
            className="bi bi-x"
            viewBox="0 0 16 16"
          >
            <path d="M4.646 4.646a.5.5 0 0 1 .708 0L8 7.293l2.646-2.647a.5.5 0 0 1 .708.708L8.707 8l2.647 2.646a.5.5 0 0 1-.708.708L8 8.707l-2.646 2.647a.5.5 0 0 1-.708-.708L7.293 8 4.646 5.354a.5.5 0 0 1 0-.708z" />
          </svg>
          Discard changes
        </button>
      </nav>
    </>
  );
}
