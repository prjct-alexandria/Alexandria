import * as React from "react";
import { Link } from "react-router-dom";
import { Request } from "./CompareView";

type MRListProps = {
  MR: {
    request: Request;
    sourceTitle: string;
    targetTitle: string;
  };
};

export default function MRListElement(props: MRListProps) {
  return (
    <div>
      <Link
        to={
          "/articles/" +
          props.MR.request.articleID +
          "/requests/" +
          props.MR.request.requestID
        }
        className="text-decoration-none"
      >
        <button
          className="row row-no-gutters col-md-12 m-1"
          style={{ textAlign: "left" }}
        >
          <div className="col-md-4" data-testid={"sourceTitle" + props.MR.request.requestID}>{props.MR.sourceTitle}</div>
          {/*<div className="col-md-2">{props.MR.sourceHistoryID}</div>*/}
          <div className="col-md-4" data-testid={"targetTitle" + props.MR.request.requestID}>{props.MR.targetTitle}</div>
          {/*<div className="col-md-2">{props.MR.targetHistoryID}</div>*/}
          <div className="col-md-4" data-testid={"status" + props.MR.request.requestID}>{props.MR.request.status}</div>
        </button>
      </Link>
    </div>
  );
}
