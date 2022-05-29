import * as React from "react";
import { useEffect, useState } from "react";
import ArticleListElement from "./ArticleListElement";
import LoadingSpinner from "../LoadingSpinner";

type Article = {
  id: string;
  title: string;
  date_created: string;
  author: string;
  description: string;
};

export default function ArticleList() {
  let [articleListData, setData] = useState<Article[]>();
  let [isLoaded, setLoaded] = useState(false);
  let [error, setError] = useState(null);

  useEffect(() => {
    // const url = "http://localhost:8080/articles/";
    const url = "/articleList.json"; // Placeholder

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
    <div className={"accordion"}>
      {!isLoaded && <LoadingSpinner />}
      {error && <div>{`There is a problem fetching the data - ${error}`}</div>}
      {articleListData != null &&
        articleListData.map((article, i) => (
          <ArticleListElement key={i} article={article} />
        ))}
    </div>
  );
}
