import * as React from "react";
import { useState } from "react";

export default function CreateArticle() {
  let [articleTitle, setArticleTitle] = useState<string>("");
  let [articleOwners, setArticleOwners] = useState<string>("");
  let [articleTags, setArticleTags] = useState<string>("");
  let [error, setError] = useState(null);
  let [isTitleChanged, setIsTitleChanged] = useState<boolean>(false);

  // Variable and references to it to be removed when adding tags
  let areTagsImplemented = false;

  const onChangeTitle = (e: { target: { value: any } }) => {
    setArticleTitle(e.target.value);
    setIsTitleChanged(true);
  };

  const onChangeOwners = (e: { target: { value: any } }) => {
    setArticleOwners(e.target.value);
  };

  const onChangeTags = (e: { target: { value: any } }) => {
    setArticleTags(e.target.value);
  };

  // Send an HTTP POST request to /register with user info
  const submitHandler = (e: { preventDefault: () => void }) => {
    // Prevent unwanted default browser behavior
    e.preventDefault();

    const url = "http://localhost:8080/articles";

    // Make list of strings from input string separated by ","
    const tagList: string[] = articleTags.split(",");
    let ownerList: string[] = [];

    // Remove extra spaces
    tagList.map((tag) => tag.trim());

    if (ownerList.length > 0) {
      ownerList = articleOwners.split(",");
      ownerList.map((owner) => owner.trim());
    }

    let loggedUser = localStorage.getItem("loggedUserEmail");
    ownerList[ownerList.length] = loggedUser == undefined ? "" : loggedUser;

    // Construct request body
    const body = {
      title: articleTitle,
      owners: ownerList,
    };

    // Send request with a JSON of the user data
    fetch(url, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      mode: "cors",
      body: JSON.stringify(body),
    }).then(
      // Success
      async (response) => {
        console.log("Success:", response);

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
        }
      },
      (error) => {
        // Request returns an error
        console.error("Error:", error);
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
      {error != null && <span>{error}</span>}
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
                {articleTitle.length == 0 && isTitleChanged && (
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
                  {articleTags.length == 0 && (
                    <span className="input-validation">
                      This field is mandatory.
                    </span>
                  )}
                </div>
              )}
              <div>
                <h5 className="form-label">Other owners *</h5>
                <span>Separate owner emails by ",".</span>
                <input
                  name="owners"
                  className="create-article-input"
                  onChange={onChangeOwners}
                />
              </div>
              <span>Fields marked with * are optional.</span>
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
                disabled={articleTitle.length == 0}
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
