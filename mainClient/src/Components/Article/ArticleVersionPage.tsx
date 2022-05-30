import * as React from "react";
import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import LoadingSpinner from "../LoadingSpinner";

type ArticleVersion = {
  owners: Array<string>;
  title: string;
  content: string;
};

export default function ArticleVersionPage() {
  let [versionData, setData] = useState<ArticleVersion>();
  let [isLoaded, setLoaded] = useState(false);
  let [error, setError] = useState(null);

  let params = useParams();
  const url =
    "http://localhost:8080/articles/" +
    params.articleId +
    "/versions/" +
    params.versionId;

  useEffect(() => {
    fetch(url, {
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
  }, [url]);

  return (
    <div className={"wrapper"}>
      {!isLoaded && <LoadingSpinner />}
      {error && <div>{`There is a problem fetching the data - ${error}`}</div>}
      {versionData != null && (
        <div className="article">
          <ul>
            {versionData.owners.map((owner, i) => (
              <li key={i}>{owner}</li>
            ))}
          </ul>
          <h1>{versionData.title}</h1>
          <div>{versionData.content}</div>
        </div>
      )}
    </div>
  );
}
