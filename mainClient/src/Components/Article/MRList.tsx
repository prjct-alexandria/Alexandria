import * as React from "react";
import { useEffect, useState } from "react";
import LoadingSpinner from "../LoadingSpinner";
import MRListElement from "./MRListElement";
import { useParams, useSearchParams } from "react-router-dom";
import NotificationAlert from "../NotificationAlert";
import configData from "../../config.json";

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

  let [MRListData, setMRListData] = useState<MR[]>();
  let [isLoaded, setLoaded] = useState<boolean>(false);
  let [error, setError] = useState<Error>();

  const baseUrl = configData.back_end_url +"/articles/" + params.articleId + "/requests";

  let appendUrl = "?" + searchParams.toString();
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
    }).then(
      async (response) => {
        if (response.ok) {
          let MRList: MR[] = await response.json();
          setMRListData(MRList);
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
