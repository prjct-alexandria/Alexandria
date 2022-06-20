import { useState } from "react";
import { useParams } from "react-router-dom";
import * as React from "react";

type EditorProps = {
  content: string;
};

export default function InlineEditor(props: EditorProps) {
  let [editorData, setEditorData] = useState<string>(props.content);

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
    }).then((response) => {
      if (response.ok) {
        props.content = editorData;
      }
    });
  };

  const discardChanges = () => {
    setEditorData(props.content);
  };

  // @ts-ignore
  return (
    <>
      <nav className="nav d-grid gap-2 d-md-flex">
        <button className={"nav-item"} type="button" onClick={saveChanges}>
          Save changes
        </button>
        <button className={"nav-item"} type="button" onClick={discardChanges}>
          Discard changes
        </button>
      </nav>
      <div className="articleContent">
        <textarea
          style={{ whiteSpace: "pre-line" }}
          className="col-xs-12 textarea"
          value={editorData}
          onChange={(e) => {
            setEditorData(e.target.value);
          }}
        />
      </div>
    </>
  );
}
