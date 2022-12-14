import * as React from "react";
import {useState} from "react";
import CreateComment from "./CreateComment";

type ThreadProps = {
    "id": undefined | number,
    "specificId": string | undefined,
    "threadType": string
};

export default function CreateThread(props: ThreadProps) {
    let [key, setKey] = useState(1);
    let [newThreadList, setNewThreadList] = useState([
        <CreateComment key={key} id={(props.id) ? undefined : props.id} specificId={props.specificId}
                       threadType={props.threadType} selection={undefined}/>
    ]);

    const addThread = () => {
        setKey(key+1);
        setNewThreadList(newThreadList.concat(
            <CreateComment key={key} id={(props.id) ? undefined : props.id} specificId={props.specificId}
                           threadType={props.threadType} selection={undefined}/>
        ));
    }

    return (
        <div className="text-center">
            {/*Leave out the first one. This one is used to initialize the array with the right type*/}
            {newThreadList.slice(1)}
            <button className="btn btn-primary m-2" type="submit" onClick={addThread}>+</button>
        </div>
    );
}