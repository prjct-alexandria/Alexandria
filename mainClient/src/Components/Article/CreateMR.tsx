import * as React from "react";
import {useCallback, useEffect, useState} from "react";
import {useParams} from "react-router-dom";
import VersionListElement from "./VersionListElement";
import {isDisabled} from "@testing-library/user-event/dist/utils";

type Version = {
    articleID: string;
    versionID: string;
    title: string;
    owners: string[];
    status: string;
};

export default function CreateMR() {
    let params = useParams();

    // const urlCurrentArticle = "http://localhost:8080/articles/" + params.articleId + "/versions/" + params.versionId;
    const urlCurrentArticle = "/article_version1.json"

    let [dataCurVersion, setDataCurVersion] = useState<Version>();
    let [isLoaded, setLoaded] = useState(false);
    let [error, setError] = useState(null);

    useEffect(() => {
        fetch(urlCurrentArticle, {
            method: "GET",
            mode: "cors",
            headers: {
                "Content-Type": "application/json",
                Accept: "application/json",
            },
            credentials: 'include',
        })
            .then((res) => res.json())
            .then(
                (result) => {
                    setDataCurVersion(result);
                    setLoaded(true);
                },
                (error) => {
                    setError(error.message);
                    setLoaded(true);
                }
            );
    }, [urlCurrentArticle]);

    // const urlArticleVersions = "http://localhost:8080/articles/" + params.articleId + "/versions";
    const urlArticleVersions = "/versionList.json"

    let [dataVersions, setDataVersions] = useState<Version[]>();

    useEffect(() => {
        fetch(urlArticleVersions, {
            method: "GET",
            mode: "cors",
            headers: {
                "Content-Type": "application/json",
                Accept: "application/json",
            },
            credentials: 'include',
        })
            .then((res) => res.json())
            .then(
                (result) => {
                    setDataVersions(result);
                    setLoaded(true);
                },
                (error) => {
                    setError(error.message);
                    setLoaded(true);
                }
            );
    }, [urlArticleVersions]);

    const [selectedClient,setSelectedClient] = useState(""); //default value

    function handleSelectChange(event: { target: { value: React.SetStateAction<string>; }; }) {
        setSelectedClient(event.target.value);
    }

    //post

    return (
        <div
            className="modal fade create-article-modal"
            id="createMR"
            data-bs-backdrop="static"
            data-bs-keyboard="false"
            aria-labelledby="publishArticleLabel"
            aria-hidden="true"
        >
            {/*{error != null && <span>{error}</span>}*/}
            <div className="modal-dialog">
                <div className="modal-content">
                    <div className="modal-header">
                        <h5 className="modal-title" id="publishArticleLabel">
                            Create request to apply changes
                        </h5>
                    </div>
                    <div className="modal-body">
                        <div className='row justify-content-center'>
                            <div className='col-6'>
                                <h6>From</h6>
                            </div>
                            <div className='col-6'>
                                <h6>To</h6>
                            </div>
                            <div className='col-6'>
                                {dataCurVersion !== undefined && dataCurVersion.title}
                            </div>
                            <div className='col-6'>
                                <select className="form-select" aria-label="Default select example"
                                        value={selectedClient} onChange={handleSelectChange}>
                                    {dataVersions != null &&
                                        dataVersions.map((version, i) => (
                                            <option value={version.versionID}>{version.title}</option>
                                    ))}
                                </select>
                                <p>You selected {selectedClient}</p>
                            </div>
                        </div>
                    </div>
                    <div className="modal-footer">
                        <button type="button" className="btn btn-secondary" data-bs-dismiss="modal">Close</button>
                        <button type="button" className="btn btn-primary">Save changes</button>
                    </div>
                </div>
            </div>
        </div>
    );
}
