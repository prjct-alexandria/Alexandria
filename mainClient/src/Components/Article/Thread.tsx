import * as React from "react";
import {useParams} from "react-router-dom";
import {useEffect, useState} from "react";
import LoadingSpinner from "../LoadingSpinner";
import CreateComment from "./CreateComment"

type ThreadProps = {
    thread: {
        "id": number,
    };
    "specificId": number
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

    const params = useParams();

    useEffect(() => {
        let urlCommentList = "";
        if (props.threadType === "commit") {
            urlCommentList = "http://localhost:8080/articles/" + params.articleId + "/versions/" + params.versionId +
            "/history/" + params.historyId + "/thread/" + props.thread.id + "/comments"
        } else if (props.threadType === "request") {
            urlCommentList = "http://localhost:8080/articles/" + params.articleId + "/requests/" + params.requestId +
            "/thread/" + props.thread.id + "/comments"
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
                        <CreateComment thread={props.thread} specificId={props.specificId} threadType={props.threadType}/>
                    </div>
                </div>
            }
        </div>

    );
}
