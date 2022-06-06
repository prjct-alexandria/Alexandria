import * as React from "react";
import { useEffect, useState } from "react";
import ArticleListElement from "./ArticleListElement";
import LoadingSpinner from "../LoadingSpinner";
import {useParams} from "react-router-dom";
import Thread from "./Thread"

type Thread = {
    "articleId": number,
    "comment": Comment[]
    "id": number,
    "specificId": number
};

type ThreadListProps = {
    threadType: string
};

export default function ThreadList(props: ThreadListProps) {
    let [threadListData, setData] = useState<Thread[]>();
    let [isLoaded, setLoaded] = useState(false);
    let [error, setError] = useState(null);

    const params = useParams();

    useEffect(() => {
        // const urlThreadList = "http://localhost:8080/articles/4/thread/" + props.threadType;
        const urlThreadList = "/threadList.json"; // Placeholder

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
                        <Thread key={i} thread={thread}/>
                    ))}

            </div>
        </div>
    );
}
