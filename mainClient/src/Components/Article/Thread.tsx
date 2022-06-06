import * as React from "react";
import {useParams} from "react-router-dom";
import {useEffect, useState} from "react";

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


    return (
        <div className="accordion-item mb-3" style={{border: '1px solid #e9ecef'}}>
            <button className="accordion-button collapsed"
                    type="button"
                    data-bs-toggle="collapse"
                    data-bs-target={"#panelsStayOpen-collapse" + props.thread.id}
                    aria-expanded="false"
                    aria-controls={"panelsStayOpen-collapse" + props.thread.id}>
                {commentData != null && commentData[0].content}
            </button>
            <div
                id={"panelsStayOpen-collapse" + props.thread.id}
                className="accordion-collapse collapse"
                aria-labelledby={"panelsStayOpen-heading" + props.thread.id}
            >
                    {commentData != null &&
                    commentData.map((comment, i) => (
                        i !== 0 && // don't show first element in the list
                        <div className="accordion-body" style={{border: '1px solid #e9ecef'}}>
                            {comment.content}
                        </div>
                    ))}
            </div>
        </div>
    );
}
