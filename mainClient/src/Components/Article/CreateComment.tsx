import * as React from "react";
import {useParams} from "react-router-dom";
import {useEffect, useState} from "react";

type ThreadProps = {
    thread: undefined | {
        "id": number,
    };
    "specificId": number
    "threadType": string
};

export default function CreateComment(props: ThreadProps) {
    const params = useParams();
    let [error, setError] = useState(null);
    let [newCommentContent, setNewCommentContent] = useState<string>("");
    let [threadId, setThreadId] = useState((props.thread == undefined) ? undefined : props.thread.id)

    // post new comment in existing thread
    const addComment = (e: { target: { value: any } }) => {
        setNewCommentContent(e.target.value);
    };

    const submitHandler = (e: { preventDefault: () => void }) => {
        // Prevent unwanted default browser behavior
        e.preventDefault();

        // If the comment is not a reply on an existing thread, create a new thread
        if (threadId == undefined) {

            // the endpoint is depends on what type of thread it is
            let urlCreateThread = "http://localhost:8080/articles/" + params.articleId + "/thread/" +
                props.threadType + "/id/" + props.specificId;

            const bodyThread = {
                articleId: params.articleId,
                comment: null,
                specificId: props.specificId
            }

            fetch(urlCreateThread, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                mode: "cors",
                body: JSON.stringify(bodyThread)
            }).then(
                // Success
                async (response) => {
                    console.log("Success:", response);

                    if (response.ok) {
                        let responseJSON: {
                            threadId: string;
                        } = await response.json();

                        setThreadId(parseInt(responseJSON.threadId as string));
                    }
                },
                (error) => {
                    // Request returns an error
                    console.error("Error:", error);
                    setError(error);
                }
            );
        }

        // add comment to thread (either the one that is just created or one that already existed)
        let urlReplyComment = "";
        if (props.threadType == "commit" && props.specificId != -1 && threadId != undefined) {
            urlReplyComment = "http://localhost:8080/articles/" + params.articleId + "/versions/" + params.versionId +
                "/history/" + props.specificId + "/threads/" + threadId + "/comments";
        } else if (props.threadType == "request" && threadId != undefined) {
            urlReplyComment = "http://localhost:8080/articles/" + params.articleId + "/request/" +
                props.specificId + "/threads/" + threadId + "/comments";
        }

        // Construct request body
        const bodyComment = {
            content: newCommentContent,
        };

        // Send request with a JSON of the comment data
        fetch(urlReplyComment, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            mode: "cors",
            body: JSON.stringify(bodyComment),
        }).then(
            (response) => {
                let message: string =
                    response.status === 200
                        ? "Comment successfully placed"
                        : "Error: " + response.status + " " + response.statusText;
                console.log(message);
                // refresh page
                window.location.reload();
            },
            (error) => {
                console.error("Error: ", error);
                setError(error);
            }
        );
    };

    return (
        <div className="mb-2" style={{border: '#e9ecef 1px solid'}}>
            {error != null && <span>{error}</span>}
            <form className="text-end" onSubmit={submitHandler}>
                <div className="form-group accordion-body">
                    <textarea className="form-control" placeholder="write comment here..."
                              onChange={addComment}></textarea>
                </div>
                <button type="submit" className="btn btn-primary m-1" disabled={newCommentContent == ""}>Submit</button>
            </form>
        </div>

    );
}
