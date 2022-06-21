import * as React from "react";
import { useParams } from "react-router-dom";
import { useEffect, useState } from "react";
import LoadingSpinner from "../LoadingSpinner";
import CreateComment from "./CreateComment";
import configData from "../../config.json";
import NotificationAlert from "../NotificationAlert";
import isUserLoggedIn from "../User/AuthHelpers/isUserLoggedIn";
import moment from 'moment';

type ThreadProps = {
    "id": number,
    "specificId": string | undefined,
    "threadType": string,
    "comments": ThreadComment[],
};

type ThreadComment = {
    "id": number,
    "authorId": string,
    "threadId": number,
    "content": string,
    "creationDate": string,
}

export default function Thread(props: ThreadProps) {
    const date = new Date(parseInt(props.comments[0].creationDate as string) * 1000);
    console.log(date)
    return (
        <div>
            {
                <div className="accordion-item mb-3 text-break" style={{border: '1px solid #e9ecef'}}>
                    <button className="accordion-button collapsed"
                        type="button"
                        data-bs-toggle="collapse"
                        data-bs-target={"#panelsStayOpen-collapse" + props.id}
                        aria-expanded="false"
                        aria-controls={"panelsStayOpen-collapse" + props.id}>
                        <div className="row">
                            <div className="toast-header mb-2" style={{backgroundColor: 'transparent'}}>
                                <strong className="me-auto p-1">
                                    {props.comments[0].authorId}
                                </strong>
                                <small className="text-muted">
                                    {moment(
                                        new Date(parseInt(props.comments[0].creationDate as string) * 1000)
                                    ).format('DD-MM-YYYY HH:mm')}
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
                                <strong className="me-auto p-1">
                                    {comment.authorId}
                                </strong>
                                <small className="text-muted float-end">
                                    {moment(
                                        new Date(parseInt(comment.creationDate as string) * 1000)
                                    ).format('DD-MM-YYYY HH:mm')}
                                </small>
                            </div>
                            <div>
                            {comment.content}
                            </div>
                        </div>
                    ))}
                    <CreateComment id={props.id} specificId={props.specificId}
                                   threadType={props.threadType} selection={undefined}/>
                </div>
            </div>
            }
        </div>

    );
}
