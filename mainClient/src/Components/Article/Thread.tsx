import * as React from "react";
import {useParams} from "react-router-dom";
import {useEffect, useState} from "react";
import LoadingSpinner from "../LoadingSpinner";

type ThreadProps = {
    thread: {
        "articleId": number,
        "id": number,
        "specificId": number
    };
    threadType: string
};

type ThreadComment = {
    "authorId": string,
    "content": string,
    "creationDate": string,
    "id": number,
    "threadId": number
}

export default function Thread(props: ThreadProps) {
    let [commentData, setData] = useState<ThreadComment[]>();
    let [isLoaded, setLoaded] = useState(false);
    let [error, setError] = useState(null);
    let [newCommentContent, setNewCommentContent] = useState<string>("");

    const params = useParams();

    useEffect(() => {
        // const urlCommentList = "http://localhost:8080/articles/" + params.articleId + "/thread/" +
        //     props.threadType + "/id/" + props.thread.specificId + "/comments";
        const urlCommentList = "/commentList1.json"; // Placeholder

        // get comments of specific thread
        fetch(urlCommentList, {
            method: "GET",
            mode: "cors",
            headers: {
                "Content-Type": "application/json",
                Accept: "application/json",
            },
            credentials: 'include',
        })
            .then((res) => res.json())
            .then(
                (result) => {
                    setLoaded(true);
                    setData(result);
                },
                (error) => {
                    setLoaded(true);
                    setError(error);
                }
            );
    }, []);


    const urlReplyComment = "http://localhost:8080/articles/" + params.articleId + "/versions/" + params.versionId + "/thread/" +
    props.thread.id + "/comments";
    // post new comment in existing thread

    const replyComment = (e: { target: { value: any } }) => {
        setNewCommentContent(e.target.value);
    };

    const submitHandler = (e: { preventDefault: () => void }) => {
        // Prevent unwanted default browser behavior
        e.preventDefault();

        // Construct request body
        const body = {
            content: newCommentContent,
        };

        // Send request with a JSON of the comment data
        fetch(urlReplyComment, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            mode: "cors",
            body: JSON.stringify(body),
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
        <div>
            {error != null && <span>{error}</span>}
            {!isLoaded && <LoadingSpinner/>}
            {
                commentData != null &&
                <div className="accordion-item mb-3" style={{border: '1px solid #e9ecef'}}>
                    <button className="accordion-button collapsed"
                            type="button"
                            data-bs-toggle="collapse"
                            data-bs-target={"#panelsStayOpen-collapse" + props.thread.id}
                            aria-expanded="false"
                            aria-controls={"panelsStayOpen-collapse" + props.thread.id}>
                        {commentData[0].content}
                    </button>
                    <div
                        id={"panelsStayOpen-collapse" + props.thread.id}
                        className="accordion-collapse collapse"
                        aria-labelledby={"panelsStayOpen-heading" + props.thread.id}
                    >
                        {commentData.map((comment, i) => (
                            i !== 0 && // don't show first element in the list
                            <div className="accordion-body" style={{border: '1px solid #e9ecef'}} key={i}>
                                {comment.content}
                            </div>
                        ))}
                        <form className="text-end" onSubmit={submitHandler}>
                            <div className="form-group accordion-body">
                                <textarea className="form-control" placeholder="write comment here..."
                                          onChange={replyComment}></textarea>
                            </div>
                            <button type="submit" className="btn btn-primary m-1">Submit</button>
                        </form>
                    </div>
                </div>
            }
        </div>

    );
}
