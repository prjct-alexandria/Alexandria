import * as React from "react";
import { Link } from "react-router-dom";
import { useEffect, useState } from "react";
import LoadingSpinner from "../LoadingSpinner";
import NotificationAlert from "../NotificationAlert";

type MRListProps = {
  MR: {
    requestID: number;
    articleID: number;
    sourceVersionID: number;
    sourceHistoryID: number;
    targetVersionID: number;
    targetHistoryID: number;
    state: string;
  };
};

type Version = {
  id: number;
  author: string;
  title: string;
  date_created: string;
  status: string;
};

export default function MRListElement(props: MRListProps) {
  let [sourceVersionData, setSourceVersionData] = useState<Version>();
  let [targetVersionData, setTargetVersionData] = useState<Version>();
  let [isLoaded, setLoaded] = useState<boolean>(false);
  let [error, setError] = useState<Error>();

  // use these URLs to get the name of the versions
  // const urlSource = "/version.json"; // Placeholder
  // const urlTarget = "/version.json"; // Placeholder
  const urlSource =
    "/articles/" + props.MR.articleID + "/versions/" + props.MR.sourceVersionID;
  const urlTarget =
    "/articles/" + props.MR.articleID + "/versions/" + props.MR.targetVersionID;

  useEffect(() => {
    fetch(urlSource, {
      method: "GET",
      mode: "cors",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
      },
    }).then(
      async (response) => {
        if (response.ok) {
          let sourceData: Version = await response.json();
          setSourceVersionData(sourceData);
        } else {
          // Set error with message returned from the server
          let responseJSON: {
            message: string;
          } = await response.json();

          let serverMessage: string = responseJSON.message;
          setError(new Error(serverMessage));
        }
        setLoaded(true);
      },
      (error) => {
        setLoaded(true);
        setError(error);
      }
    );
  }, []);

  useEffect(() => {
    fetch(urlTarget, {
      method: "GET",
      mode: "cors",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
      },
    }).then(
      async (response) => {
        if (response.ok) {
          let targetData: Version = await response.json();
          setSourceVersionData(targetData);
        } else {
          // Set error with message returned from the server
          let responseJSON: {
            message: string;
          } = await response.json();

          let serverMessage: string = responseJSON.message;
          setError(new Error(serverMessage));
        }
        setLoaded(true);
      },
      (error) => {
        setLoaded(true);
        setError(error);
      }
    );
  }, []);

  return (
    <div>
      {!isLoaded && <LoadingSpinner />}
      {error && (
        <NotificationAlert
          errorType="danger"
          title={"Error: "}
          message={"Something went wrong. " + error}
        />
      )}
      {sourceVersionData && targetVersionData && (
        <Link
          to={
            "/articles/" +
            props.MR.articleID +
            "/requests/" +
            props.MR.requestID
          }
          className="text-decoration-none"
        >
          <button className="row row-no-gutters col-md-12 m-1">
            <div className="col-md-2">
              {sourceVersionData && sourceVersionData.title}
            </div>
            <div className="col-md-2">{props.MR.sourceHistoryID}</div>
            <div className="col-md-2">
              {targetVersionData && targetVersionData.title}
            </div>
            <div className="col-md-2">{props.MR.targetHistoryID}</div>
            <div className="col-md-4">{props.MR.state}</div>
          </button>
        </Link>
      )}
    </div>
  );
}
