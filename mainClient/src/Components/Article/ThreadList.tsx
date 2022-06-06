import * as React from "react";
import { useEffect, useState } from "react";
import LoadingSpinner from "../LoadingSpinner";
import {useParams} from "react-router-dom";
import Thread from "./Thread"

type ThreadListProps = {
    threadType: string
    specificId: number
};

type ThreadComment = {
    "authorId": string,
    "content": string,
    "creationDate": string,
    "id": number,
    "threadId": number
}

type Thread = {
    "articleId": number,
    "comment": ThreadComment[]
    "id": number,
    "specificId": number
};

export default function ThreadList(props: ThreadListProps) {
    let [threadListData, setData] = useState<Thread[]>();
    let [isLoaded, setLoaded] = useState(false);
    let [error, setError] = useState(null);

    const params = useParams();

    useEffect(() => {
        // const urlThreadList = "http://localhost:8080/articles/" + params.articleId + "/thread/" + props.threadType
        //  + "/id/" + props.specificId;

        const urlThreadList = "/threadList.json"; // Placeholder
        // get list of threads of a certain article
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
            <div className="accordion" id="accordionPanelsStayOpenExample">
                {threadListData != null &&
                    threadListData.map((thread, i) => (
                        <Thread key={i} thread={thread} threadType={props.threadType}/>
                    ))}
            </div>
        </div>
    );
}
