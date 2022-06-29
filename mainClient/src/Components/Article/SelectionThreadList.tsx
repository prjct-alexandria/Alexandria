import * as React from "react";
import { useEffect, useState } from "react";
import LoadingSpinner from "../LoadingSpinner";
import { useParams } from "react-router-dom";
import NotificationAlert from "../NotificationAlert";
import isUserLoggedIn from "../User/AuthHelpers/isUserLoggedIn";
import SelectionThread from "./SelectionThread";
import backEndUrl from "../../urlUtils";

type ThreadListProps = {
    "threadType": string
    "specificId": string | undefined
};

type ThreadComment = {
    "id": number,
    "authorId": string,
    "threadId": number
    "content": string,
    "creationDate": string,
}

type ThreadEntity = {
    "articleId": number
    "comments": ThreadComment[]
    "id": number
    "specificId": string | undefined
    "selection": string
}


// A selectionThreadList is a list of threads that can be related to specific pieces of text (selections) in the document.
export default function SelectionThreadList(props: ThreadListProps) {
    let baseUrl= backEndUrl();
    let [threadListData, setData] = useState<ThreadEntity[]>();
    let [isLoaded, setLoaded] = useState<boolean>(false);
    let [error, setError] = useState<Error>();

    let [isLoggedIn, setLoggedIn] = useState<boolean>(isUserLoggedIn());

    // Listen for userAccountEvent that fires when user in localstorage changes
    window.addEventListener("userAccountEvent", () => {
        setLoggedIn(isUserLoggedIn());
    });

    const params = useParams();

    useEffect(() => {
        let urlThreadList = baseUrl + "/articles/" + params.articleId + "/versions/" + params.versionId +
            "/history/" + props.specificId + "/selectionThreads";

        // get list of threads
        fetch(urlThreadList, {
            method: "GET",
            mode: "cors",
            headers: {
                "Content-Type": "application/json",
                Accept: "application/json",
            },
            credentials: 'include',
        }).then(
            async (response) => {
                if (response.ok) {
                    let threadListData: ThreadEntity[] = await response.json();
                    if (threadListData != null) {
                        threadListData = threadListData.reverse()
                    }
                    setData(threadListData);
                } else {
                    // Set error with message returned from the server
                    let responseJSON: {
                        message: string;
                    } = await response.json();

                    let serverMessage: string = responseJSON.message;
                    setError(new Error(serverMessage));
                }
                setLoaded(true);
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
            {error && (
                <NotificationAlert
                    errorType="danger"
                    title={"Error: "}
                    message={"Something went wrong. " + error}
                />
            )}
            <div id="accordionPanelsStayOpenExample">
                <h5 className="mb-2">Comments on specific selections</h5>
                {threadListData != null &&
                    threadListData.map((thread, i) => (
                        <SelectionThread
                            key={i}
                            id={thread.id}
                            specificId={props.specificId}
                            threadType={props.threadType}
                            comments={thread.comments}
                            selection={thread.selection}
                        />
                    ))}
            </div>
        </div>
    );
}
