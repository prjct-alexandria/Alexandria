import * as React from "react"
import { useParams } from "react-router-dom"
import {useEffect, useState} from "react";
import PrismDiff from "./PrismDiff";
import axios from 'axios';
import LoadingSpinner from "../LoadingSpinner";

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

    // const urlRequest = '/request.json'
    const urlRequest = 'http://localhost:8080/articles/' + params.articleId  + "/requests/" + params.requestId;

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
        // urlArticleSource = '/article_version1.json'; // Placeholder source version
        urlArticleSource = 'http://localhost:8080/articles/' + params.articleId + '/versions/' + params.versionId + '/history/' + dataRequest.sourceHistoryID;
        // urlArticleTarget = '/article_version2.json'; // Placeholder target version
        urlArticleTarget = 'http://localhost:8080/articles/' + params.articleId + '/versions/' + dataRequest.targetVersionID + '/history/' + dataRequest.targetHistoryID;
    }

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
        const acceptRequest = await axios.get('/articles/' + params.articleId + '/versions/' + params.versionId + '/requests/'
            + params.requestId + '/merge')
        window.location.reload()
    }

    // Send HTTP request and reload, so that the style (see "rejectButton") is updated.
    const handleRejectClick = async () => {
        const rejectRequest = await axios.get('/articles/' + params.articleId + '/versions/' + params.versionId + '/requests/'
            + params.requestId + '/reject')
        window.location.reload()
    }

    // Send HTTP request and reload, so that the style (see "deleteButton") is updated.
    const handleDeleteClick = async () => {
        const deleteRequest = await axios.delete('/articles/' + params.articleId + '/versions/' + params.versionId + '/requests/'
            + params.requestId + '/delete')
        window.location.reload()
    }

    const view = () => {
        return (
            <div>
                {/*Delete button*/}
                <div className='mt-3' style={{position:'absolute', right:'5%'}}>
                    {deleteButton()}
                </div>

                <h1>See changes</h1>

                {/*Accept and reject button*/}
                <div className='row justify-content-center'>
                    <div className='col-1' id='AcceptButton'>
                        {acceptButton()}
                    </div>
                    <div className='col-1'>
                        {rejectButton()}
                    </div>
                </div>


                <div className='row justify-content-center'>
                    {/*Version names*/}
                    <div className='col-4'>
                        Version: {dataSource !== undefined && dataSource.title}
                    </div>
                    <div className='col-4' >
                        Version: {dataTarget !== undefined && dataTarget.title}
                    </div>

                    {/*Content of versions*/}
                    <div className='col-8 '  >
                        <div className='row overflow-scroll' style={{height:'500px'}}>
                            {/*Source version, including changes that are made*/}
                            <div className='col-6' style={{border: 'black solid 3px'}}>
                                {(dataSource !== undefined && dataTarget !== undefined) &&
                                    <PrismDiff
                                        sourceContent={dataSource.content}
                                        targetContent={dataTarget.content}
                                    />
                                }
                            </div>
                            {/*Target version*/}
                            <div className='col-6' style={{border: 'black solid 3px'}}>
                                {dataTarget !== undefined && dataTarget.content}
                            </div>
                        </div>
                    </div>
                    {/*Space for threads regarding this request*/}
                    <div className='col-8 ' style={{border: 'black solid 1px', height: '100px'}}>
                        Threads
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