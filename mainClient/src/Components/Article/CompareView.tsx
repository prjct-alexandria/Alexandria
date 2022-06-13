import * as React from "react"
import { useParams } from "react-router-dom"
import {useEffect, useState} from "react";
import PrismDiff from "./PrismDiff";
import LoadingSpinner from "../LoadingSpinner";
import ThreadList from "./ThreadList"

type RequestWithComparison = {
    request: Request;
    source: ArticleVersion;
    target: ArticleVersion;
    before: string;
    after: string;
}

type Request = {
    requestID: number;
    articleID: number;
    sourceVersionID: number;
    sourceHistoryID: number;
    targetVersionID: number;
    targetHistoryID: number;
    status: string;
    conflicted: boolean;
}

type ArticleVersion = {
    id: number;
    title: string;
    owners: string[];
    content: string;
};

export default function VersionList() {
    let params = useParams();

    //const urlRequest = '/request.json'
    const urlRequest = 'http://localhost:8080/articles/' + params.articleId  + "/requests/" + params.requestId;

    let [dataRequest, setDataRequest] = useState<RequestWithComparison>();
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



    // Disable button if it is already accepted or rejected
    const disableButton = () => {
        return (dataRequest !== undefined && (dataRequest.request.status === "accepted" || dataRequest.request.status === "rejected"))
    }

    // If already accepted, fill in the color of the button and disable. If request is pending, send HTTP request.
    const acceptButton = () => {
        let className = 'btn btn-outline-success';
        let disabledConflicted = dataRequest !== undefined && dataRequest.request.conflicted
        console.log(dataRequest !== undefined && dataRequest.request.status)
        if (dataRequest !== undefined && dataRequest.request.status === "accepted") {className = 'btn btn-success'}
        return (<button className={className}  disabled={disableButton() || disabledConflicted} onClick={handleAcceptClick}>Accept</button>)
    }

    // If already rejected, fill in the color of the button and disable. If request is pending, send HTTP request.
    const rejectButton = () => {
        let className = 'btn btn-outline-danger';
        if (dataRequest !== undefined && dataRequest.request.status === "rejected") {className = 'btn btn-danger'}
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
            <div className="row">
                <div>
                    <h1 style={{textAlign:"center", marginBottom:"30px"}}>Compare Changes</h1>
                    <div className='row justify-content-center'>
                        {/*Version names*/}
                        <div className='row col-8 mb-2'>
                            <div className='col-6'>
                                <h5>Changes of '{dataRequest !== undefined && dataRequest.source.title}'</h5>
                            </div>
                            <div className='col-6'>
                                <h5>Result: {dataRequest !== undefined && dataRequest.target.title}</h5>
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


                    {dataRequest !== undefined && dataRequest.request.conflicted &&
                        <div className={'alert alert-danger'} style={{whiteSpace:'pre-line'}}>
                            <p>{"Warning: There are conflicting changes in the two versions that this request would merge.\nThe conflicts are highlighted in the preview below between each set of  '<<<<', '====', and '>>>>' markers."}</p>
                        </div>
                    }


                    <div className={'row justify-content-center'}>
                        {/*Content of versions*/}
                        <div className={"wrapper col-8"}>
                            <div className='row overflow-scroll' style={{height:'500px',whiteSpace: 'pre-line', border: 'grey solid 3px'}}>
                                {/*Differences between before and after*/}
                                <div className='col-6'>
                                    {(dataRequest !== undefined) &&
                                        <PrismDiff
                                            sourceContent={dataRequest.before}
                                            targetContent={dataRequest.after}
                                        />
                                    }
                                </div>
                                {/*Result: so after only*/}
                                <div className='col-6'>
                                    {dataRequest !== undefined && dataRequest.after}
                                </div>
                            </div>
                        </div>
                        <div className="wrapper col-3">
                            <ThreadList threadType={"request"} specificId={parseInt(params.requestId as string)}/>
                        </div>
                    </div>
                </div>
            </div>
        );
    }

    return (
        <div>
            {!isLoadedRequest && <LoadingSpinner/>}
            {errorRequest && (<div>{`There is a problem fetching the post data - ${errorRequest}`}</div>)}
            {dataRequest != null && dataRequest.before != null && dataRequest.after != null && view()}
        </div>


    )
}