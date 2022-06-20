import * as React from "react";
import { Link, useParams } from "react-router-dom";

type VListProps = {
  version: {
    versionID: string;
    title: string;
    owners: string[];
    status: string;
  };
  mv: string | undefined;
};

export default function VersionListElement(props: VListProps) {
  let params = useParams();
  let baseLink = "/articles/" + params.articleId + "/requests";

  return (
    <div className="row row-no-gutters col-md-12 text-wrap">
      <div className="col-md-1">
        {props.version.versionID === props.mv && (
          <span className="badge bg-success">Main</span>
        )}
      </div>
      <div className="col-md-2">
        <Link
          to={
            "/articles/" +
            params.articleId +
            "/versions/" +
            props.version.versionID
          }
        >
          {props.version.title}
        </Link>
      </div>
      <div className="col-md-2">
        <Link to={baseLink + "?relatedID=" + props.version.versionID}>
          See all related requests
        </Link>
      </div>
      <div className="col-md-2">
        <Link to={baseLink + "?sourceID=" + props.version.versionID}>
          See requests as source
        </Link>
      </div>
      <div className="col-md-2">
        <Link to={baseLink + "?targetID=" + props.version.versionID}>
          See requests as target
        </Link>
      </div>

      <div className="col-md-2">{props.version.owners.join(", ")}</div>
      <div className="col-md-1">{props.version.status}</div>
    </div>
  );
}
