import * as React from "react";
import { Link } from "react-router-dom";
import { useEffect, useState } from "react";
import LoadingSpinner from "../LoadingSpinner";
import NotificationAlert from "../NotificationAlert";
import configData from "../../config.json";

type MRListProps = {
  MR: {
    requestID: number;
    articleID: number;
    sourceVersionID: number;
    sourceHistoryID: number;
    targetVersionID: number;
    targetHistoryID: number;
    status: string;
    conflicted: boolean;
    sourceTitle: string;
    targetTitle: string;
  };
};

type Version = {
  id: number;
  author: string;
  title: string;
  date_created: string;
  status: string;
};

export default function MRListElement(props: MRListProps) {
  let [error, setError] = useState<Error>();

  return (
    <div>
      {error && (
        <NotificationAlert
          errorType="danger"
          title={"Error: "}
          message={"Something went wrong. " + error}
        />
      )}
      {(
        <Link
          to={
            "/articles/" +
            props.MR.articleID +
            "/requests/" +
            props.MR.requestID
          }
          className="text-decoration-none"
        >
          <button className="row row-no-gutters col-md-12 m-1">
            <div className="col-md-4">
              {props.MR.sourceTitle}
            </div>
            {/*<div className="col-md-2">{props.MR.sourceHistoryID}</div>*/}
            <div className="col-md-4">
              {props.MR.targetTitle}
            </div>
            {/*<div className="col-md-2">{props.MR.targetHistoryID}</div>*/}
            <div className="col-md-4">{props.MR.status}</div>
          </button>
        </Link>
      )}
    </div>
  );
}
