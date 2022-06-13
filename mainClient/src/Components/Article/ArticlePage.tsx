import * as React from "react";
import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import LoadingSpinner from "../LoadingSpinner";
import configData from "../../config.json";

type Article = {
  authors: Array<string>;
  title: string;
  content: string;
};

export default function ArticlePage() {
  let [articleData, setData] = useState<Article>();
  let [isLoaded, setLoaded] = useState(false);
  let [error, setError] = useState(null);

  let params = useParams();
  const url =
    configData.back_end_url +"/articles/" + params.aid + "versions/" + params.vId;

  useEffect(() => {
    fetch(url, {
      method: "GET",
      mode: "cors",
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
      },
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
    <div className={"wrapper"}>
      {!isLoaded && <LoadingSpinner />}
      {error && <div>{`There is a problem fetching the data - ${error}`}</div>}
      {articleData != null && (
        <div className="article">
          <ul>
            {articleData.authors.map((author, i) => (
              <li key={i}>{author}</li>
            ))}
          </ul>
          <h1>{articleData.title}</h1>
          <div>{articleData.content}</div>
        </div>
      )}
    </div>
  );
}
