import * as React from "react";
import { Link, useParams } from "react-router-dom";

type VListProps = {
  version: {
    versionID: string;
    title: string;
    owners: string[];
    status: string;
  };
};

export default function VersionListElement(props: VListProps) {
  let params = useParams();
  return (
    <div className="row row-no-gutters col-md-12 text-wrap">
      <div className="col-md-8 text-start">
        <Link
          to={"/articles/" + params.articleId + "/versions/" + props.version.versionID}
        >
          {props.version.title}
        </Link>
      </div>
      <div className="col-md-3">{props.version.owners.join(', ')}</div>
      <div className="col-md-1">{props.version.status}</div>
    </div>
  );
}
