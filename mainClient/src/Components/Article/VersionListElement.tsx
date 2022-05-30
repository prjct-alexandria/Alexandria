import * as React from "react";
import { Link, useParams } from "react-router-dom";

type VListProps = {
  version: {
    id: string;
    author: string;
    title: string;
    date_created: string;
    status: string;
  },
  mv: string | undefined;
};

export default function VersionListElement(props: VListProps) {
  let params = useParams();
  return (
    <div className='row row-no-gutters col-md-12 text-wrap'>
      <div className='col-md-1'>
          {props.version.id == props.mv && <span className="badge bg-success">Main</span>}
      </div>
      <div className='col-md-8'>
        <Link to={"/articles/" + params.articleId + "/versions/" + props.version.id}>
          {props.version.title}
        </Link>
      </div>
      <div className='col-md-1'>{props.version.author}</div>
      <div className='col-md-1'>{props.version.date_created}</div>
      <div className='col-md-1'>{props.version.status}</div>
    </div>
  );
}
