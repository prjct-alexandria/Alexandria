import { useParams } from "react-router-dom"
import {useEffect, useState} from "react";
import VersionListElement from "./VersionListElement";

export default function VersionList() {
    let params = useParams();
    const url = '/versionList.json'; // Placeholder
    // const url = '/articles/' + params.aid + "versions" // should be this when endpoint is there

    let [data, setData] = useState([]);
    let [loading, setLoading] = useState(true);
    let [error, setError] = useState(null);

    useEffect(() => {
        fetch(url
        )
            .then(res => res.json())
            .then(
                (result) => {
                    setData(result)
                    setError(null)
                    setLoading(false)
                },
                (error) => {
                    setError(error.message)
                    setData(null)
                    setLoading(false)
                },
            )
    }, []);

    return (
        <div>
            {loading && <div>A moment please...</div>}
            {error && (<div>{`There is a problem fetching the post data - ${error}`}</div>)}
            {data != null && data.map((version, i) =>
                <VersionListElement key={i} version={version} aId={params.aid}/>
            )}
        </div>


    )
}