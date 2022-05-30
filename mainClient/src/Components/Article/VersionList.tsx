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
  const urlVersions = "/versionList.json"; // Placeholder
  //const url = "http://localhost:8080/articles/" + params.articleId + "/versions";

  let [dataVersions, setDataVersions] = useState<Version[]>();
  let [isLoadedVersions, setLoadedVersions] = useState(false);
  let [errorVersions, setErrorVersions] = useState(null);

  useEffect(() => {
    fetch(urlVersions, {
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
          setDataVersions(result);
          setErrorVersions(null);
          setLoadedVersions(true);
        },
        (error) => {
          setErrorVersions(error.message);
          setDataVersions(error);
          setLoadedVersions(true);
        }
      );
  }, [urlVersions]);

  const urlMain = "/mainVersion.json"; // Placeholder
  //const url = ...?

  let [dataMain, setDataMain] = useState<string>();
  let [isLoadedMain, setLoadedMain] = useState(false);
  let [errorMain, setErrorMain] = useState(null);

  useEffect(() => {
    fetch(urlMain, {
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
              setDataMain(result);
              setErrorMain(null);
              setLoadedMain(true);
            },
            (error) => {
              setErrorMain(error.message);
              setDataMain(error);
              setLoadedMain(true);
            }
        );
  }, [urlVersions]);

  console.log(dataMain)
  console.log(dataVersions)
  console.log(dataMain)
  return (
    <div>
      {!isLoadedVersions || !isLoadedMain && <LoadingSpinner />}
      {errorVersions || errorMain && <div>{`There is a problem fetching the data - ${errorVersions} ${errorMain}`}</div>}
      {(dataVersions != null && dataMain != null) &&
        dataVersions.map((version, i) => (
          <VersionListElement key={i} version={version} mv={dataMain} />
        ))}
    </div>
  );
}
