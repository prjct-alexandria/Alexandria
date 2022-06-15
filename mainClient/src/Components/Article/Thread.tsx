import * as React from "react";
import { useParams } from "react-router-dom";
import { useEffect, useState } from "react";
import LoadingSpinner from "../LoadingSpinner";
import CreateComment from "./CreateComment";
import NotificationAlert from "../NotificationAlert";
import isUserLoggedIn from "../User/AuthHelpers/isUserLoggedIn";

type ThreadProps = {
  id: number;
  specificId: number;
  threadType: string;
};

type ThreadComment = {
  authorId: string;
  content: string;
  creationDate: string;
  commentId: number;
  threadId: number;
};

export default function Thread(props: ThreadProps) {
  let baseUrl = "http://localhost:8080";
  let [commentData, setData] = useState<ThreadComment[]>();
  let [isLoaded, setLoaded] = useState<boolean>(false);
  let [error, setError] = useState<Error>();

  let [isLoggedIn, setLoggedIn] = useState<boolean>(isUserLoggedIn());

  // Listen for userAccountEvent that fires when user in localstorage changes
  window.addEventListener("userAccountEvent", () => {
    setLoggedIn(isUserLoggedIn());
  });

  const params = useParams();

  useEffect(() => {
    let urlCommentList = "";
    if (props.threadType === "commit") {
      urlCommentList =
        baseUrl +
        "/articles/" +
        params.articleId +
        "/versions/" +
        params.versionId +
        "/history/" +
        params.historyId +
        "/thread/" +
        props.id +
        "/comments";
    } else if (props.threadType === "request") {
      urlCommentList =
        baseUrl +
        "/articles/" +
        params.articleId +
        "/requests/" +
        params.requestId +
        "/thread/" +
        props.id +
        "/comments";
    }
    urlCommentList = "/commentList1.json"; // Placeholder

    // get comments of specific thread
    fetch(urlCommentList, {
      method: "GET",
      mode: "cors",
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
      },
      credentials: "include",
    }).then(
      async (response) => {
        if (response.ok) {
          let data: ThreadComment[] = await response.json();
          setData(data);
        } else {
          // Set error with message returned from the server
          let responseJSON: {
            message: string;
          } = await response.json();

          let serverMessage: string = responseJSON.message;
          setError(new Error(serverMessage));
        }
        setLoaded(true);
      },
      (error) => {
        setLoaded(true);
        setError(error);
      }
    );
  }, []);

  return (
    <div>
      {!isLoaded && <LoadingSpinner />}
      {error && (
        <NotificationAlert
          errorType="danger"
          title={"Error: "}
          message={"Something went wrong. " + error}
        />
      )}
      {commentData != null && (
        <div
          className="accordion-item mb-3"
          style={{ border: "1px solid #e9ecef" }}
        >
          <button
            className="accordion-button collapsed"
            type="button"
            data-bs-toggle="collapse"
            data-bs-target={"#panelsStayOpen-collapse" + props.id}
            aria-expanded="false"
            aria-controls={"panelsStayOpen-collapse" + props.id}
          >
            {commentData[0].content}
          </button>
          <div
            id={"panelsStayOpen-collapse" + props.id}
            className="accordion-collapse collapse"
            aria-labelledby={"panelsStayOpen-heading" + props.id}
          >
            {commentData.map(
              (comment, i) =>
                i !== 0 && ( // don't show first element in the list
                  <div
                    className="accordion-body"
                    style={{ border: "1px solid #e9ecef" }}
                    key={i}
                  >
                    {comment.content}
                  </div>
                )
            )}
            {isLoggedIn && (
              <CreateComment
                id={props.id}
                specificId={props.specificId}
                threadType={props.threadType}
              />
            )}
          </div>
        </div>
      )}
    </div>
  );
}
