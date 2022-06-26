import * as React from "react";
import { useEffect, useState } from "react";
import LoadingSpinner from "../LoadingSpinner";
import MRListElement from "./MRListElement";
import { useParams, useSearchParams } from "react-router-dom";
import NotificationAlert from "../NotificationAlert";
import configData from "../../config.json";
import { Request } from "./CompareView";

type MR = {
  request: Request;
  sourceTitle: string;
  targetTitle: string;
};

export default function MRList() {
  let params = useParams(); // used for the articleId
  const [searchParams] = useSearchParams(); // used for the source and target

  let [MRListData, setMRListData] = useState<MR[]>();
  let [isLoaded, setLoaded] = useState<boolean>(false);
  let [error, setError] = useState<Error>();

  const baseUrl =
    configData.back_end_url + "/articles/" + params.articleId + "/requests";

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
  },[]);

  const mrListMap = () => {
    return (
      <div>
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
        MRListData &&
        <div>
          <h2 className="text-center mb-5">See all requests</h2>
          <div className="row wrapper col-8 m-auto">
            <div className="col-md-4">
              <h5>From</h5>
            </div>
            <div className="col-md-4">
              <h5>To</h5>
            </div>
            <div className="col-md-4">
              <h5>Status</h5>
            </div>
            <div className="col-md-4">
              <h6>Source Version</h6>
            </div>
            <div className="col-md-4">
              <h6>Target Version</h6>
            </div>
            {/*<div className="col-md-2">*/}
            {/*  <h6>History:</h6>*/}
            {/*</div>*/}
            {mrListMap()}
          </div>
        </div>
      }
    </div>
  );
}
