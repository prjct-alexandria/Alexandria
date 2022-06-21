import * as React from "react";
import { useParams } from "react-router-dom";
import { useEffect, useState } from "react";
import LoadingSpinner from "../LoadingSpinner";
import CreateComment from "./CreateComment";
import configData from "../../config.json";
import NotificationAlert from "../NotificationAlert";
import isUserLoggedIn from "../User/AuthHelpers/isUserLoggedIn";

type SectionThreadProps = {
    "id": number,
    "specificId": string | undefined,
    "threadType": string,
    "comments": ThreadComment[],
    "section": string
};

type ThreadComment = {
    "id": number,
    "authorId": string,
    "threadId": number,
    "content": string,
    "creationDate": string,
}


// A SectionThread is a thread (list of comments) that can be related to a specific section in the document
export default function SectionThread(props: SectionThreadProps) {
    return (
        <div>
            <div>
                {"Related to text: \"" + props.section + '\"'}
            </div>
            {
                <div className="accordion-item mb-3" style={{border: '1px solid #e9ecef'}}>
                    <button className="accordion-button collapsed"
                    type="button"
                    data-bs-toggle="collapse"
                    data-bs-target={"#panelsStayOpen-collapse" + props.id}
                    aria-expanded="false"
                    aria-controls={"panelsStayOpen-collapse" + props.id}>
            {props.comments[0].content}
                    </button>
                        <div id={"panelsStayOpen-collapse" + props.id}
                             className="accordion-collapse collapse"
                             aria-labelledby={"panelsStayOpen-heading" + props.id}>
                            {props.comments.map((comment, i) => (
                            i !== 0 && // don't show first element in the list
                            <div className="accordion-body" style={{border: '1px solid #e9ecef'}} key={i}>
                                {comment.content}
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
