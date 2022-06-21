import * as React from "react";
import {useEffect, useState} from "react";
import CreateComment from "./CreateComment";

type SelectionThreadProps = {
    "id": undefined | number,
    "specificId": string | undefined,
    "threadType": string
    "posX" : number
    "posY": number
    "hidden": boolean
    "selection": string
};

export default function CreateSelectionThread(props: SelectionThreadProps) {
    let [key, setKey] = useState(1);

    let [newThreadList, setNewThreadList] = useState([
        <CreateComment key={key} id={(props.id) ? undefined : props.id} specificId={props.specificId}
                       threadType={props.threadType} selection={props.selection}/>
    ]);

    const addThread = () => {
        setKey(key+1);
        setNewThreadList(newThreadList.concat(
            <CreateComment key={key} id={(props.id) ? undefined : props.id} specificId={props.specificId}
                           threadType={props.threadType} selection={props.selection}/>
        ));
    }

    return (
        <div className="text-center">
            {/*Leave out the first one. This one is used to initialize the array with the right type*/}
            {newThreadList.slice(1)}
            <button className="btn btn-primary m-2" type="submit" onClick={addThread}
            style={{'position':"absolute", 'top':props.posY+'px', 'left':props.posX+'px'}}
                    hidden={props.hidden}
            >+ Add comment</button>
        </div>
    );
}