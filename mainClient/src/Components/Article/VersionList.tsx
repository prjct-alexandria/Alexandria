import * as React from "react";

import { useParams } from "react-router-dom";
import { useEffect, useState } from "react";
import VersionListElement from "./VersionListElement";
import LoadingSpinner from "../LoadingSpinner";
import NotificationAlert from "../NotificationAlert";
import backEndUrl from "../../urlUtils";

type Version = {
  articleID: string;
  versionID: string;
  title: string;
  owners: string[];
  status: string;
};

export default function VersionList() {
  let params = useParams();
  // const urlVersions = "/versionList.json"; // Placeholder
  const urlVersions =
    backEndUrl() + "/articles/" + params.articleId + "/versions";

  let [dataVersions, setDataVersions] = useState<Version[]>();
  let [isLoadedVersions, setLoadedVersions] = useState<boolean>(false);
  let [errorVersions, setErrorVersions] = useState<Error>();

  useEffect(() => {
    fetch(urlVersions, {
      method: "GET",
      mode: "cors",
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
      },
      credentials: "include",
    }).then(
      async (response) => {
        if (response.ok) {
          let versionList: Version[] = await response.json();
          setDataVersions(versionList);
        } else {
          // Set error with message returned from the server
          let responseJSON: {
            message: string;
          } = await response.json();

          let serverMessage: string = responseJSON.message;
          setErrorVersions(new Error(serverMessage));
        }
        setLoadedVersions(true);
      },
      (error) => {
        setLoadedVersions(true);
        setErrorVersions(error);
      }
    );
  }, []);

  // const urlMain = "/mainVersion.json"; // Placeholder
  const urlMain =
    backEndUrl() + "/articles/" + params.articleId + "/mainVersion";

  let [dataMain, setDataMain] = useState<string>();
  let [isLoadedMain, setLoadedMain] = useState<boolean>(false);
  let [errorMain, setErrorMain] = useState<Error>();

  useEffect(() => {
    fetch(urlMain, {
      method: "GET",
      mode: "cors",
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
      },
      credentials: "include",
    }).then(
      async (response) => {
        if (response.ok) {
          let main: string = await response.json();
          setDataMain(main);
        } else {
          // Set error with message returned from the server
          let responseJSON: {
            message: string;
          } = await response.json();

          let serverMessage: string = responseJSON.message;
          setErrorMain(new Error(serverMessage));
        }
        setLoadedMain(true);
      },
      (error) => {
        setLoadedMain(true);
        setErrorMain(error);
      }
    );
  }, []);

  return (
    <div className="wrapper col-8 m-auto">
      <div>
        {(!isLoadedVersions || !isLoadedMain) && <LoadingSpinner />}
        {(errorVersions || errorMain) && (
          <NotificationAlert
            errorType="danger"
            title={"Error: "}
            message={"Something went wrong. " + errorVersions + errorMain}
          />
        )}
        {dataVersions &&
          dataMain &&
          dataVersions
              .sort(
                  (a,b) => a.versionID > b.versionID ? 1 : -1
              )
              .map((version, i) => (
            <VersionListElement key={i} version={version} mv={dataMain} />
          ))}
      </div>
    </div>
  );
}
