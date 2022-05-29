import * as React from "react";

import { useParams } from "react-router-dom";
import { useEffect, useState } from "react";
import VersionListElement from "./VersionListElement";
import LoadingSpinner from "../LoadingSpinner";

type Version = {
  id: string;
  author: string;
  title: string;
  date_created: string;
  status: string;
};

export default function VersionList() {
  let params = useParams();
  const url = "/versionList.json"; // Placeholder
  //const url = "http://localhost:8080/articles/" + params.articleId + "/versions";

  let [data, setData] = useState<Version[]>();
  let [isLoaded, setLoaded] = useState(false);
  let [error, setError] = useState(null);

  useEffect(() => {
    fetch(url, {
      method: "GET",
      mode: "cors",
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
      },
    })
      .then((res) => res.json())
      .then(
        (result) => {
          setData(result);
          setError(null);
          setLoaded(true);
        },
        (error) => {
          setError(error.message);
          setData(error);
          setLoaded(true);
        }
      );
  }, [url]);

  return (
    <div>
      {!isLoaded && <LoadingSpinner />}
      {error && <div>{`There is a problem fetching the data - ${error}`}</div>}
      {data != null &&
        data.map((version, i) => (
          <VersionListElement key={i} version={version} />
        ))}
    </div>
  );
}
