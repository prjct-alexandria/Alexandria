import * as React from "react";
import { useEffect, useState } from "react";
import LoadingSpinner from "../LoadingSpinner";
import {useParams} from "react-router-dom";
import Thread from "./Thread"
import CreateThread from "./CreateThread";

type ThreadListProps = {
    "threadType": string
    "specificId": number
};


type ThreadComment = {
    "id": number,
    "authorId": string,
    "threadId": number
    "content": string,
    "creationDate": string,
}

type ThreadEntity = {
    "articleId":	number
    "comment": ThreadComment[]
    "id": number
    "specificId": number
}


export default function ThreadList(props: ThreadListProps) {
    let baseUrl = "http://localhost:8080";
    let [threadListData, setData] = useState<ThreadEntity[]>();
    let [isLoaded, setLoaded] = useState(false);
    let [error, setError] = useState(null);

    const params = useParams();

    useEffect(() => {
        let urlThreadList = "";
        if (props.threadType === "commit") {
            urlThreadList = baseUrl + "/articles/" + params.articleId + "/versions/" + params.versionId +
            // "/history/" + params.historyId + "/threads";
                "/history/" + 1 + "/threads";
        } else if (props.threadType === "request") {
            urlThreadList = baseUrl + "/articles/" + params.articleId + "/requests/" + params.requestId +
                "/threads";
        }
        // urlThreadList = "/threadList.json" // Placeholder

        // get list of threads
        fetch(urlThreadList, {
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
            {!isLoaded && <LoadingSpinner />}
            {error && <div>{`There is a problem fetching the data - ${error}`}</div>}
            <div id="accordionPanelsStayOpenExample">
                {threadListData != null &&
                    threadListData.map((thread, i) => (
                        <Thread key={i} id={thread.id} specificId={props.specificId} threadType={props.threadType} comments={thread.comment}/>
                    ))}
            </div>
            <CreateThread id={undefined} specificId={props.specificId} threadType={props.threadType}/>
        </div>
    );
}