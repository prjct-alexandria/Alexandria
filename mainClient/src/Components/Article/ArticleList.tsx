import * as React from "react";
import { useEffect, useState } from "react";
import ArticleListElement from "./ArticleListElement";
import LoadingSpinner from "../LoadingSpinner";
import configData from "../../config.json";

type Article = {
  articleId: string;
  mainVersionId: string;

  //Following attributes are from the main Version, but displayed as if they were from the article itself
  title: string;
  date_created: string;
  owners: string[];
  description: string;
};

export default function ArticleList() {
  let [articleListData, setData] = useState<Article[]>();
  let [isLoaded, setLoaded] = useState(false);
  let [error, setError] = useState(null);

  useEffect(() => {
    const url = configData.back_end_url + "/articles";
    //const url = "/articleList.json"; // Placeholder

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
  }, []);

  return (
    <div className="wrapper col-8 m-auto">
      <div className={"accordion"}>
        {!isLoaded && <LoadingSpinner />}
        {error && <div>{`There is a problem fetching the data - ${error}`}</div>}
        {articleListData != null &&
          articleListData.map((article, i) => (
            <ArticleListElement key={i} article={article} />
          ))}
      </div>
    </div>
  );
}
