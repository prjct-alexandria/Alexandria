import * as React from "react";
import { Link } from "react-router-dom";

type ArticleListProps = {
  article: {
    id: string;
    title: string;
    date_created: string;
    author: string;
    description: string;
  };
};

export default function ArticleListElement(props: ArticleListProps) {
  return (
      <div className="accordion-item">
        <div className="article-elem-header accordion-header">
          <h5 className="article-elem-title">
            <Link to={"/articles/" + props.article.id + "/versions/1"}>
              {props.article.title}
            </Link>
          </h5>
          <button
              className="article-elem-btn accordion-button collapsed"
              type="button"
              data-bs-toggle="collapse"
              data-bs-target={"#panelsStayOpen-collapse" + props.article.id}
              aria-expanded="false"
              aria-controls={"panelsStayOpen-collapse" + props.article.id}
          >
            <div className="d-flex justify-content-between flex-fill">
              <span className="p-2 flex-grow-1">By {props.article.author}</span>
              <span className="p-2">See details</span>
            </div>
          </button>
        </div>
        <div
            id={"panelsStayOpen-collapse" + props.article.id}
            className="accordion-collapse collapse"
            aria-labelledby={"panelsStayOpen-heading" + props.article.id}
        >
          <div className="accordion-body">{props.article.description}</div>
        </div>
      </div>
  );
}
