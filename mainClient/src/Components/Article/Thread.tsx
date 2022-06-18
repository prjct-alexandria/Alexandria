import * as React from "react";
import { useParams } from "react-router-dom";
import { useEffect, useState } from "react";
import LoadingSpinner from "../LoadingSpinner";
import CreateComment from "./CreateComment";
import configData from "../../config.json";
import NotificationAlert from "../NotificationAlert";
import isUserLoggedIn from "../User/AuthHelpers/isUserLoggedIn";

type ThreadProps = {
    "id": number,
    "specificId": string | undefined,
    "threadType": string,
    "comments": ThreadComment[],
};

type ThreadEntity = {
    "articleId": number,
    "id": number,
    "specificId": string | undefined,
    "comment": ThreadComment[],
}

type ThreadComment = {
    "id": number,
    "authorId": string,
    "threadId": number,
    "content": string,
    "creationDate": string,
}

export default function Thread(props: ThreadProps) {
    return (
        <div>
            {
                <div className="accordion-item mb-3" style={{border: '1px solid #e9ecef'}}>


                    <button className="accordion-button collapsed"
                        type="button"
                        data-bs-toggle="collapse"
                        data-bs-target={"#panelsStayOpen-collapse" + props.id}
                        aria-expanded="false"
                        aria-controls={"panelsStayOpen-collapse" + props.id}>
                        <div>
                            <div className="toast-header mb-2 p-0 row" style={{backgroundColor: 'transparent'}}>
                                <strong className="me-auto col-6">
                                    {props.comments[0].authorId}
                                </strong>
                                <small className="text-muted col-6">
                                    {props.comments[0].creationDate}
                                </small>
                            </div>
                            <div>
                                {props.comments[0].content}
                            </div>
                        </div>
                    </button>
                <div
                    id={"panelsStayOpen-collapse" + props.id}
                    className="accordion-collapse collapse"
                    aria-labelledby={"panelsStayOpen-heading" + props.id}
                >
                    {props.comments.map((comment, i) => (
                        i !== 0 && // don't show first element in the list
                        <div className="accordion-body" style={{border: '1px solid #e9ecef'}} key={i}>
                            <div className="toast-header mb-2 p-0">
                                <strong className="me-auto">
                                    {comment.authorId}
                                </strong>
                                <small className="text-muted">
                                    {comment.creationDate}
                                </small>
                            </div>
                            <div>
                            {comment.content}
                            </div>
                        </div>
                    ))}
                    <CreateComment id={props.id} specificId={props.specificId} threadType={props.threadType}/>
                </div>
            </div>
            }
        </div>

    );
}
