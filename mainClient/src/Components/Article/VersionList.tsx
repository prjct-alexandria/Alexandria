import * as React from "react";

import { useParams } from "react-router-dom";
import { useEffect, useState } from "react";
import VersionListElement from "./VersionListElement";
import LoadingSpinner from "../LoadingSpinner";
import NotificationAlert from "../NotificationAlert";

type Version = {
  articleID: string;
  versionID: string;
  title: string;
  owners: string[];
  status: string;
};

export default function VersionList() {
  let [error, setError] = useState<Error>();
  let params = useParams();
  // const urlVersions = "/versionList.json"; // Placeholder
  const urlVersions =
    "http://localhost:8080/articles/" + params.articleId + "/versions";

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
      credentials: "include",
    })
      .then((res) => res.json())
      .then(
        (response) => {
          if (response.ok) {
            setDataVersions(response);
            setErrorVersions(null);
            setLoadedVersions(true);
          } else {
            let serverMessage: string = response.message;
            setError(new Error(serverMessage));
          }
        },
        (error) => {
          setErrorVersions(error.message);
          setDataVersions(error);
          setLoadedVersions(true);
        }
      );
  }, [urlVersions]);

  // const urlMain = "/mainVersion.json"; // Placeholder
  const urlMain =
    "http://localhost:8080/articles/" + params.articleId + "/mainVersion";

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
      credentials: "include",
    })
      .then((res) => res.json())
      .then(
        (response) => {
          if (response.ok) {
            setDataMain(response);
            setErrorMain(null);
            setLoadedMain(true);
          } else {
            let serverMessage: string = response.message;
            setError(new Error(serverMessage));
          }
        },
        (error) => {
          setErrorMain(error.message);
          setDataMain(error);
          setLoadedMain(true);
        }
      );
  }, [urlVersions]);

  return (
    <div className="wrapper col-8 m-auto">
      {error && (
        <NotificationAlert
          errorType="danger"
          title={"Error: "}
          message={"Something went wrong. " + error.message}
        />
      )}
      <div>
        {!isLoadedVersions || (!isLoadedMain && <LoadingSpinner />)}
        {(errorVersions || errorMain) && (
          <NotificationAlert
            errorType="danger"
            title={"Error: "}
            message={"Something went wrong. " + errorVersions + errorMain}
          />
        )}
        {dataVersions != null &&
          dataMain != null &&
          dataVersions.map((version, i) => (
            <VersionListElement key={i} version={version} mv={dataMain} />
          ))}
      </div>
    </div>
  );
}
