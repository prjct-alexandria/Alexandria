import * as React from "react";
import { useState } from "react";
import NotificationAlert from "../NotificationAlert";

export default function CreateArticle() {
  let [mainVersionTitle, setMainVersionTitle] = useState<string>("");
  let [mainVersionOwners, setMainVersionOwners] = useState<string>("");
  let [mainVersionTags, setMainVersionTags] = useState<string>("");
  let [error, setError] = useState<Error>();
  let [isTitleChanged, setIsTitleChanged] = useState<boolean>(false);

  // Variable and references to it to be removed when adding tags
  let areTagsImplemented = false;

  const onChangeTitle = (e: { target: { value: any } }) => {
    setMainVersionTitle(e.target.value);
    setIsTitleChanged(true);
  };

  const onChangeOwners = (e: { target: { value: any } }) => {
    setMainVersionOwners(e.target.value);
  };

  const onChangeTags = (e: { target: { value: any } }) => {
    setMainVersionTags(e.target.value);
  };

  // Send an HTTP POST request to /articles with new article's info
  const submitHandler = (e: { preventDefault: () => void }) => {
    // Prevent unwanted default browser behavior
    e.preventDefault();

    const url = "http://localhost:8080/articles";

    // Make list of strings from input string separated by ","
    let tagList: string[] = mainVersionTags.split(",");
    let ownerList: string[] = mainVersionOwners.split(",");

    // Remove extra spaces
    tagList = tagList.map((tag) => tag.trim());
    ownerList = ownerList.map((owner) => owner.trim());

    // Remove empty elements
    tagList = tagList.filter((tag) => tag != "");
    ownerList = ownerList.filter((owner) => owner != "");

    let loggedUser = localStorage.getItem("loggedUserEmail");
    ownerList[ownerList.length] =
      loggedUser === null || typeof loggedUser === "undefined"
        ? ""
        : loggedUser;

    // Construct request body
    const body = {
      title: mainVersionTitle,
      owners: ownerList,
    };

    // Send request with a JSON of the user data
    fetch(url, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      mode: "cors",
      body: JSON.stringify(body),
      credentials: "include"
    }).then(
      async (response) => {
        if (response.ok) {
          let responseJSON: {
            articleID: string;
            versionID: string;
          } = await response.json();

          const articleId: string = responseJSON.articleID;
          const versionId: string = responseJSON.versionID;

          if (typeof window !== "undefined") {
            window.location.href =
              "http://localhost:3000/articles/" +
              articleId +
              "/versions/" +
              versionId;
          }
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
        // Request returns an error
        setError(error);
      }
    );
  };

  return (
    <div
      className="modal fade create-article-modal"
      id="publishArticle"
      data-bs-backdrop="static"
      data-bs-keyboard="false"
      aria-labelledby="publishArticleLabel"
      aria-hidden="true"
    >
      {error && (
        <NotificationAlert
          errorType="danger"
          title={"Error: "}
          message={"Something went wrong. " + error}
        />
      )}
      <div className="modal-dialog">
        <div className="modal-content">
          <div className="modal-header">
            <h5 className="modal-title" id="publishArticleLabel">
              Create new article
            </h5>

            <button
              type="button"
              className="btn-close"
              data-bs-dismiss="modal"
              aria-label="Close"
            ></button>
          </div>
          <form onSubmit={submitHandler}>
            <div className="modal-body">
              <div>
                <h5 className="form-label">Title</h5>
                <input
                  name="title"
                  className="create-article-input"
                  onChange={onChangeTitle}
                />
                {mainVersionTitle.length === 0 && isTitleChanged && (
                  <span className="input-validation">
                    This field is mandatory.
                  </span>
                )}
              </div>
              {areTagsImplemented && (
                <div>
                  <h5 className="form-label">Tags</h5>
                  <span>Separate tags by ",".</span>
                  <input
                    name="tags"
                    className="create-article-input"
                    onChange={onChangeTags}
                  />
                  {mainVersionTags.length === 0 && (
                    <span className="input-validation">
                      This field is mandatory.
                    </span>
                  )}
                </div>
              )}
              <div>
                <h5 className="form-label">Other owners (optional)</h5>
                <span>Separate owner emails by ",".</span>
                <input
                  name="owners"
                  className="create-article-input"
                  onChange={onChangeOwners}
                />
              </div>
            </div>
            <div className="modal-footer">
              <button
                type="button"
                className="btn btn-secondary"
                data-bs-dismiss="modal"
              >
                Cancel
              </button>
              <button
                type="submit"
                className="btn btn-primary"
                disabled={mainVersionTitle.length === 0}
              >
                Submit
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
}
