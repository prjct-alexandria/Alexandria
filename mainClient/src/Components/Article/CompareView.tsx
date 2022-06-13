import * as React from "react"
import { useParams } from "react-router-dom"
import {useEffect, useState} from "react";
import PrismDiff from "./PrismDiff";
import axios from 'axios';
import LoadingSpinner from "../LoadingSpinner";
import ThreadList from "./ThreadList"

type Request = {
    sourceVersionID: number;
    sourceHistoryID: number;
    targetVersionID: number;
    targetHistoryID: number;
    status: string;
}

type ArticleVersion = {
    id: number;
    title: string;
    owners: string[];
    content: string;
};

export default function VersionList() {
    let params = useParams();

    const urlRequest = '/request.json'
    //const urlRequest = 'http://localhost:8080/articles/' + params.articleId  + "/requests/" + params.requestId;

    let [dataRequest, setDataRequest] = useState<Request>();
    let [isLoadedRequest, setLoadedRequest] = useState(false);
    let [errorRequest, setErrorRequest] = useState(null);

    // fetching information about the request: historyID of source version, versionID of target, historyID of target, state of the request
    useEffect(() => {
        fetch(urlRequest
        )
            .then(res => res.json())
            .then(
                (result) => {
                    setDataRequest(result)
                    setLoadedRequest(true)
                },
                (error) => {
                    setErrorRequest(error.message)
                    setLoadedRequest(true)
                },
            )
    }, []);

    let urlArticleSource = "";
    let urlArticleTarget = "";

    if (dataRequest !== undefined) {
        urlArticleSource = 'http://localhost:8080/articles/' + params.articleId + '/versions/' + params.versionId + '/history/' + dataRequest.sourceHistoryID;
        urlArticleTarget = 'http://localhost:8080/articles/' + params.articleId + '/versions/' + dataRequest.targetVersionID + '/history/' + dataRequest.targetHistoryID;
    }
    urlArticleSource = '/article_version1.json'; // Placeholder source version
    urlArticleTarget = '/article_version2.json'; // Placeholder target version

    let [dataSource, setDataSource] = useState<ArticleVersion>();
    let [isLoadedSource, setLoadedSource] = useState(false);
    let [errorSource, setErrorSource] = useState(null);

    let [dataTarget, setDataTarget] = useState<ArticleVersion>();
    let [isLoadedTarget, setLoadedTarget] = useState(false);
    let [errorTarget, setErrorTarget] = useState(null);

    // fetching the actual articles
    useEffect(() => {
        fetch(urlArticleSource, {
            method: "GET",
            mode: "cors",
            headers: {
                "Content-Type": "application/json",
                Accept: "application/json",
            },
        })
            .then(res => res.json())
            .then(
                (result) => {
                    setDataSource(result)
                    setLoadedSource(true)
                },
                (error) => {
                    setErrorSource(error.message)
                    setLoadedSource(true)
                },
            )
        fetch(urlArticleTarget, {
            method: "GET",
            mode: "cors",
            headers: {
                "Content-Type": "application/json",
                Accept: "application/json",
            },
        })
            .then(res => res.json())
            .then(
                (result) => {
                    setDataTarget(result)
                    setLoadedTarget(true)
                },
                (error) => {
                    setErrorTarget(error.message)
                    setLoadedTarget(true)
                },
            )
    }, []);

    // Disable button if it is already accepted or rejected
    const disableButton = () => {
        return (dataRequest !== undefined && (dataRequest.status === "accepted" || dataRequest.status === "rejected"))
    }

    // If already accepted, fill in the color of the button and disable. If request is pending, send HTTP request.
    const acceptButton = () => {
        let className = 'btn btn-outline-success';
        if (dataRequest !== undefined && dataRequest.status === "accepted") {className = 'btn btn-success'}
        return (<button className={className}  disabled={disableButton()} onClick={handleAcceptClick}>Accept</button>)
    }

    // If already reject, fill in the color of the button and disable. If request is pending, send HTTP request.
    const rejectButton = () => {
        let className = 'btn btn-outline-danger';
        if (dataRequest !== undefined && dataRequest.status === "rejected") {className = 'btn btn-danger'}
        return (<button className={className}  disabled={disableButton()} onClick={handleRejectClick}>Reject</button>)
    }

    // If already accepted or rejected, do not show the button. If pending, send HTTP request.
    const deleteButton = () => {
        if (!disableButton()) {
            return (<button className={'btn btn-danger'} onClick={handleDeleteClick}>Delete</button>)
        }
    }

    // Send HTTP request and reload, so that the style (see "acceptButton") is updated.
    const handleAcceptClick = async () => {
        const url = 'http://localhost:8080/articles/' + params.articleId + '/requests/' + params.requestId + '/accept'
        fetch(url, {
            method: "PUT",
            headers: { "Content-Type": "application/json" },
            mode: "cors",
            credentials: "include"
        }).then(
            // Success
            (response) => {
                console.log(response)
            },
            (error) => {
                // Request returns an error
                console.error("Error:", error);
            }
        );
        window.location.reload()
    }

    // Send HTTP request and reload, so that the style (see "rejectButton") is updated.
    const handleRejectClick = async () => {
        const url = 'http://localhost:8080/articles/' + params.articleId + '/requests/' + params.requestId + '/reject'
        fetch(url, {
            method: "PUT",
            headers: { "Content-Type": "application/json" },
            mode: "cors",
            credentials: "include"
        }).then(
            // Success
            (response) => {
                console.log(response)
            },
            (error) => {
                // Request returns an error
                console.error("Error:", error);
            }
        );
        window.location.reload()
    }

    // Send HTTP request and reload, so that the style (see "deleteButton") is updated.
    const handleDeleteClick = async () => {
        const url = 'http://localhost:8080/articles/' + params.articleId + '/requests/' + params.requestId
        fetch(url, {
            method: "DELETE",
            headers: { "Content-Type": "application/json" },
            mode: "cors",
            credentials: "include"
        }).then(
            // Success
            (response) => {
                console.log(response)
            },
            (error) => {
                // Request returns an error
                console.error("Error:", error);
            }
        );
    }

    const view = () => {
        return (
            <div className='row'>
                <div>
                    <h1 style={{textAlign:"center", marginBottom:"30px"}}>Compare Versions</h1>
                    <div className='row justify-content-center'>
                        {/*Version names*/}
                        <div className='row col-8 mb-2'>
                            <div className='col-6'>
                                <h5>Changes of '{dataTarget !== undefined && dataTarget.title}'</h5>
                            </div>
                            <div className='col-6'>
                                <h5>Result: {dataSource !== undefined && dataSource.title}</h5>
                            </div>
                        </div>

                        {/*Accept, reject and delete button*/}
                        <div className='col-1' id='AcceptButton'>
                            {acceptButton()}
                        </div>
                        <div className='col-1'>
                            {rejectButton()}
                        </div>
                        <div className='col-1'>
                            {deleteButton()}
                        </div>
                    </div>

                    <div className='row justify-content-center'>
                        {/*Content of versions*/}
                        <div className="wrapper col-8">
                            <div className='row overflow-scroll' style={{height:'500px',whiteSpace: 'pre-line', border: 'grey solid 3px'}}>
                                {/*Source version, including changes that are made*/}
                                <div className='col-6'>
                                    {(dataSource !== undefined && dataTarget !== undefined) &&
                                        <PrismDiff
                                            sourceContent={dataSource.content}
                                            targetContent={dataTarget.content}
                                        />
                                    }
                                </div>
                                {/*Target version*/}
                                <div className='col-6'>
                                    {dataTarget !== undefined && dataTarget.content}
                                </div>
                            </div>
                        </div>
                        <div className="wrapper col-3">
                            <ThreadList threadType={"request"} specificId={parseInt(params.requestId as string)}/>
                        </div>
                    </div>
                </div>
            </div>
        )
    }

    return (
        <div>
            {!isLoadedRequest && !isLoadedSource && !isLoadedTarget && <LoadingSpinner />}
            {errorRequest && (<div>{`There is a problem fetching the post data - ${errorRequest}`}</div>)}
            {errorSource && (<div>{`There is a problem fetching the post data - ${errorSource}`}</div>)}
            {errorTarget && (<div>{`There is a problem fetching the post data - ${errorTarget}`}</div>)}
            {dataRequest != null && dataSource != null && dataTarget != null && view()}
        </div>


    )
}