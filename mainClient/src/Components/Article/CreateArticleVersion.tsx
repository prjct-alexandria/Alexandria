import * as React from "react";
import { useState } from "react";
import { useParams } from "react-router-dom";

export default function CreateArticleVersion() {
    let [newVersionTitle, setNewVersionTitle] = useState<string>("");
    let [newVersionOwners, setNewVersionOwners] = useState<string>("");
    let [newVersionTags, setNewVersionTags] = useState<string>("");
    let [error, setError] = useState(null);
    let [isTitleChanged, setIsTitleChanged] = useState<boolean>(false);

    let params = useParams()

    // Variable and references to it to be removed when adding tags
    let areTagsImplemented = false;

    const onChangeTitle = (e: { target: { value: any } }) => {
        setNewVersionTitle(e.target.value);
        setIsTitleChanged(true);
    };

    const onChangeOwners = (e: { target: { value: any } }) => {
        setNewVersionOwners(e.target.value);
    };

    const onChangeTags = (e: { target: { value: any } }) => {
        setNewVersionTags(e.target.value);
    };

    // Send an HTTP POST request to /articles with new article's info
    const submitHandler = (e: { preventDefault: () => void }) => {
        // Prevent unwanted default browser behavior
        e.preventDefault();

        const url = "http://localhost:8080/articles/" +
        params.articleId +
        "/versions";


        // Make list of strings from input string separated by ","
        let tagList: string[] = newVersionTags.split(",");
        let ownerList: string[] = newVersionOwners.split(",");

        // Remove extra spaces
        tagList = tagList.map((tag) => tag.trim());
        ownerList = ownerList.map((owner) => owner.trim());

        // Remove empty elements
        tagList = tagList.filter((tag) => tag!== '')
        ownerList = ownerList.filter((owner) => owner!== '')

        let loggedUser = localStorage.getItem("loggedUserEmail");
        ownerList[ownerList.length] =
            loggedUser === null || typeof loggedUser === "undefined"
                ? ""
                : loggedUser;

        // Construct request body
        const body = {
            sourceVersionId: parseInt(params.versionId as string),
            title: newVersionTitle,
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
            id="createNewVersion"
            data-bs-backdrop="static"
            data-bs-keyboard="false"
            aria-labelledby="createNewVersionLabel"
            aria-hidden="true"
        >
            {error != null && <span>{error}</span>}
            <div className="modal-dialog">
                <div className="modal-content">
                    <div className="modal-header">
                        <h5 className="modal-title" id="createNewVersionLabel">
                            Create new version based off this one
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
                                {newVersionTitle.length === 0 && isTitleChanged && (
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
                                    {newVersionTags.length === 0 && (
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
                                disabled={newVersionTitle.length === 0}
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
