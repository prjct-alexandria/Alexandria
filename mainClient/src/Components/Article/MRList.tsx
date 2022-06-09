import * as React from "react";
import { useEffect, useState } from "react";
import LoadingSpinner from "../LoadingSpinner";
import MRListElement from "./MRListElement";
import { useParams, useSearchParams } from "react-router-dom";
import NotificationAlert from "../NotificationAlert";

type MR = {
  requestID: number;
  articleID: number;
  sourceVersionID: number;
  sourceHistoryID: number;
  targetVersionID: number;
  targetHistoryID: number;
  state: string;
};

export default function MRList() {
  let params = useParams(); // used for the articleId
  const [searchParams] = useSearchParams(); // used for the source and target
  let sourceVersionID = searchParams.get("source");
  let targetVersionID = searchParams.get("target");

  let [MRListData, setMRListData] = useState<MR[]>();
  let [isLoaded, setLoaded] = useState(false);
  let [error, setError] = useState(null);

  const baseUrl =
    "http://localhost:8080/articles/" + params.articleId + "/requests";

  let appendUrl = "";
  if (sourceVersionID != null && targetVersionID != null) {
    appendUrl = "?source=" + sourceVersionID + "&target=" + targetVersionID;
  } else if (sourceVersionID != null) {
    appendUrl = "?source=" + sourceVersionID;
  } else if (targetVersionID != null) {
    appendUrl = "?target=" + targetVersionID;
  }

  // const url = "/requestList.json"; // Placeholder
  const url = baseUrl + appendUrl;

  useEffect(() => {
    fetch(url, {
      method: "GET",
      mode: "cors",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
      },
    })
      .then((res) => res.json())
      .then(
        (result) => {
          setLoaded(true);
          setMRListData(result);
        },
        (error) => {
          setLoaded(true);
          setError(error);
        }
      );
  }, []);

  const mrListMap = () => {
    return (
      <div className="wrapper col-8 m-auto">
        {!isLoaded && <LoadingSpinner />}
        {MRListData != null &&
          MRListData.map((mr, i) => <MRListElement key={i} MR={mr} />)}
      </div>
    );
  };

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
      {
        <div>
          <h2 className="text-center mb-5">See all requests</h2>
          <div className="row row-no-gutters col-md-12 m-1">
            <div className="col-md-4">
              <h5>From</h5>
            </div>
            <div className="col-md-4">
              <h5>to</h5>
            </div>
            <div className="col-md-4">
              <h5>State</h5>
            </div>
            <div className="col-md-2">
              <h6>Version:</h6>
            </div>
            <div className="col-md-2">
              <h6>History:</h6>
            </div>
            <div className="col-md-2">
              <h6>Version:</h6>
            </div>
            <div className="col-md-2">
              <h6>History:</h6>
            </div>
          </div>
          <div>{mrListMap()}</div>
        </div>
      }
    </div>
  );
}
