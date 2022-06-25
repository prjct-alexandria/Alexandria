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
    <div className="row row-no-gutters col-md-12 text-wrap mb-2 border-bottom border-dark">

      <div className="col-md-1">
        {props.version.versionID == props.mv && (
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
          data-testid={"title" + props.version.versionID}
        >
          {props.version.title}
        </Link>
      </div>
      <div className="col-md-6">
          <div className="btn-group dropend">
              <Link to={baseLink + "?relatedID=" + props.version.versionID}>
                  <button type="button" className="btn btn-secondary mb-1">
                      See all requests
                  </button>
              </Link>
              <button type="button" className="btn btn-secondary btn-sm dropdown-toggle dropdown-toggle-split mb-1"
                      data-bs-toggle="dropdown" aria-expanded="false">
                  <span className="visually-hidden">Toggle Dropright</span>
              </button>
              <ul className="dropdown-menu">
                  <li><a className="dropdown-item" href={baseLink + "?relatedID=" + props.version.versionID}>
                      See requests as source
                  </a></li>
                  <li><a className="dropdown-item" href={baseLink + "?targetID=" + props.version.versionID}>
                      See requests as target
                  </a></li>
              </ul>
          </div>
      </div>

      <div className="col-md-2" data-testid={"owners" + props.version.versionID}>{props.version.owners.join(", ")}</div>
      <div className="col-md-1" data-testid={"status" + props.version.versionID}>{props.version.status}</div>
    </div>
  );
}
