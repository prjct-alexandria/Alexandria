import * as React from "react";
import {useParams} from "react-router-dom";
import {useState} from "react";

type ThreadProps = {
    "id": undefined | number,
    "specificId": number
    "threadType": string
};

export default function CreateComment(props: ThreadProps) {
    let baseUrl = "http://localhost:8080";
    let [error, setError] = useState(null);
    let [newCommentContent, setNewCommentContent] = useState<string>("");
    let [threadId, setThreadId] = useState((props.id) ? undefined : props.id)

    const params = useParams();

    // post new comment in existing thread
    const addComment = (e: { target: { value: any } }) => {
        setNewCommentContent(e.target.value);
    };

    const submitHandler = (e: { preventDefault: () => void }) => {
        // Prevent unwanted default browser behavior
        e.preventDefault();

        // If the comment is not a reply on an existing thread, create a new thread
        if (!threadId) {
            // the endpoint is depends on what type of thread it is
            let urlCreateThread = baseUrl +  "/articles/" + params.articleId + "/thread/" +
                props.threadType + "/id/" + props.specificId;

            const bodyThread = {
                articleId: parseInt(params.articleId as string),
                comment: [{
                    "authorId": localStorage.getItem('loggedUserEmail'),
                    // immediately add comment when creating a new thread to make
                    // sure a thread cannot be without comments
                    "content": newCommentContent,
                    "creationDate": new Date()
                }],
                specificId: props.specificId
            }

            fetch(urlCreateThread, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                mode: "no-cors",
                body: JSON.stringify(bodyThread),
                credentials: 'include',
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
                    // refresh page, remove this for debugging
                    window.location.reload();
                },
                (error) => {
                    // Request returns an error
                    console.error("Error:", error);
                    setError(error);
                }
            );
        } else {
            // add comment to thread (either the one that is just created or one that already existed)
            let urlReplyComment = baseUrl + "/comments/thread/" + threadId

            // Construct request body
            const bodyComment = {
                "authorId": localStorage.getItem('loggedUserEmail'),
                "content": newCommentContent,
                "creationDate": new Date()
            };

            // Send request with a JSON of the comment data
            fetch(urlReplyComment, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                mode: "cors",
                body: JSON.stringify(bodyComment),
                credentials: 'include',
            }).then(
                (response) => {
                    let message: string =
                        response.status === 200
                            ? "Comment successfully placed"
                            : "Error: " + response.status + " " + response.statusText;
                    console.log(message);
                    // refresh page, remove this for debugging
                    window.location.reload();
                },
                (error) => {
                    console.error("Error: ", error);
                    setError(error);
                }
            );
        }

    };

    return (
        <div className="mb-2" style={{border: '#e9ecef 1px solid'}}>
            {error && <span>{error}</span>}
            <form className="text-end" onSubmit={submitHandler}>
                <div className="form-group accordion-body">
                    <textarea className="form-control" placeholder="write comment here..."
                              onChange={addComment}></textarea>
                </div>
                <button type="submit" className="btn btn-primary m-1" disabled={newCommentContent === ""}>Submit</button>
            </form>
        </div>

    );
}
