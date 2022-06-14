import * as React from "react";
import { useEffect, useState } from "react";
import LoadingSpinner from "../LoadingSpinner";
import { useParams } from "react-router-dom";
import Thread from "./Thread";
import CreateThread from "./CreateThread";

type ThreadListProps = {
  threadType: string;
  specificId: number;
};

type ThreadEntity = {
  articleId: number;
  id: number;
  specificId: number;
};

export default function ThreadList(props: ThreadListProps) {
  let baseUrl = "http://localhost:8080";
  let [threadListData, setData] = useState<ThreadEntity[]>();
  let [isLoaded, setLoaded] = useState<boolean>(false);
  let [error, setError] = useState<Error>();

  const params = useParams();

  useEffect(() => {
    let urlThreadList = "";
    if (props.threadType === "commit") {
      urlThreadList =
        baseUrl +
        "/articles/" +
        params.articleId +
        "/versions/" +
        params.versionId +
        "/history/" +
        params.historyId +
        "/threads";
    } else if (props.threadType === "request") {
      urlThreadList =
        baseUrl +
        "/articles/" +
        params.articleId +
        "/request/" +
        params.requestId +
        "/threads";
    }
    urlThreadList = "/threadList.json"; // Placeholder

    // get list of threads
    fetch(urlThreadList, {
      method: "GET",
      mode: "cors",
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
      },
      credentials: "include",
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
            <Thread
              key={i}
              id={thread.id}
              specificId={props.specificId}
              threadType={props.threadType}
            />
          ))}
      </div>
      <CreateThread
        id={undefined}
        specificId={props.specificId}
        threadType={props.threadType}
      />
    </div>
  );
}
