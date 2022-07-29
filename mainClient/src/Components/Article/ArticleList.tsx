import * as React from "react";
import { useEffect, useState } from "react";
import ArticleListElement from "./ArticleListElement";
import LoadingSpinner from "../LoadingSpinner";
import NotificationAlert from "../NotificationAlert";
import backEndUrl from "../../urlUtils";

type Article = {
  articleId: string;
  mainVersionId: string;
  // Following attributes are from the main Version,
  // but displayed as if they were from the article itself
  title: string;
  date_created: string;
  owners: string[];
  description: string;
};

export default function ArticleList() {
  let [articleList, setArticleList] = useState<Article[]>();
  let [isLoaded, setLoaded] = useState<boolean>(false);
  let [error, setError] = useState<Error>();

  useEffect(() => {
    const url = backEndUrl() + "/articles";

    fetch(url, {
      method: "GET",
      mode: "cors",
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
      },
      credentials: "include",
    }).then(
      async (response) => {
        if (response.ok) {
          setError(undefined);
          let articleList: Article[] = await response.json();
          setArticleList(articleList);
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
    <div className="wrapper col-8 m-auto">
      <div className={"accordion"}>
        {!isLoaded && <LoadingSpinner />}
        {error && (
          <NotificationAlert
            errorType="danger"
            title={"Error: "}
            message={"Something went wrong when getting the articles. " + error}
          />
        )}
        {articleList &&
          articleList.map((article, i) => (
            <ArticleListElement key={i} article={article} />
          ))}
      </div>
    </div>
  );
}
