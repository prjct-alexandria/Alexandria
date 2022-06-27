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
    <div className="row mb-2">
        <button
          className="text-start btn btn-secondary"
          onClick={(e) => {
            e.preventDefault();
            window.location.href = "/articles/" +
                props.MR.request.articleID +
                "/requests/" +
                props.MR.request.requestID;
          }}
        >
          <div className="row">
            <div className="col-md-4" data-testid={"sourceTitle" + props.MR.request.requestID}>{props.MR.sourceTitle}</div>
            {/*<div className="col-md-2">{props.MR.sourceHistoryID}</div>*/}
            <div className="col-md-4" data-testid={"targetTitle" + props.MR.request.requestID}>{props.MR.targetTitle}</div>
            {/*<div className="col-md-2">{props.MR.targetHistoryID}</div>*/}
            <div className="col-md-4" data-testid={"status" + props.MR.request.requestID}>{props.MR.request.status}</div>
          </div>
        </button>
    </div>
  );
}
