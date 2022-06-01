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
  let baseLink = '/articles/' + params.aid + '/requests'

  return (
    <div className="row row-no-gutters col-md-12 text-wrap">
      <div className="col-md-4 text-start">
        <Link
          to={"/articles/" + params.articleId + "/versions/" + props.version.versionID}
        >
          {props.version.title}
        </Link>
      </div>
      <div className="col-md-2"><Link to={baseLink}>See all requests</Link></div>
      <div className="col-md-2"><Link to={baseLink + '?source=' + props.version.versionID}>See requests as source</Link></div>
      <div className="col-md-2"><Link to={baseLink + '?target=' + props.version.versionID}>See requests as target</Link></div>

      <div className="col-md-1">{props.version.owners.join(', ')}</div>
      <div className="col-md-1">{props.version.status}</div>
    </div>
  );
}
