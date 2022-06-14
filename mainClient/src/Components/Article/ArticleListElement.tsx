import * as React from "react";
import { Link } from "react-router-dom";

type ArticleListProps = {
  article: {
    articleId: string;
    mainVersionId: string;

    //Following attributes are from the main Version, but displayed as if they were from the article itself
    title: string;
    date_created: string;
    owners: string[];
    description: string;
  };
};

export default function ArticleListElement(props: ArticleListProps) {
  return (
    <div className="accordion-item">
      <div className="article-elem-header accordion-header">
        <h5 className="article-elem-title" data-testid={"title" + props.article.id}>
          <Link to={"/articles/" + props.article.articleId + "/versions/" + props.article.mainVersionId}>
            {props.article.title}
          </Link>
        </h5>
        <button
          className="article-elem-btn accordion-button collapsed"
          type="button"
          data-bs-toggle="collapse"
          data-bs-target={"#panelsStayOpen-collapse" + props.article.articleId}
          aria-expanded="false"
          aria-controls={"panelsStayOpen-collapse" + props.article.articleId}
        >
          <div className="d-flex justify-content-between flex-fill">
            <span className="p-2 flex-grow-1" data-testid={"author" + props.article.id}>By {props.article.owners.join(", ")}</span>
            <span className="p-2">See details</span>
          </div>
        </button>
      </div>
      <div
        id={"panelsStayOpen-collapse" + props.article.articleId}
        className="accordion-collapse collapse"
        aria-labelledby={"panelsStayOpen-heading" + props.article.articleId}
      >
        <div className="accordion-body" data-testid={"description" + props.article.id}>{props.article.description}</div>
      </div>
    </div>
  );
}