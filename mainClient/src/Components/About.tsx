import LoadingSpinner from "./LoadingSpinner";
import ArticleListElement from "./Article/ArticleListElement";
import * as React from "react";

export default function NotificationAlert() {
    return (
        <div className="wrapper col-8 m-auto">
                <p>This application is created by five second year Computer Science and Engineering students from TU Delft,
                in the academic year 2021-2022:
                </p>
                <ul>
                    <li>Amy van der Meijden</li>
                    <li>Andreea Zlei</li>
                    <li>Emiel Witting</li>
                    <li>Jos Sloof</li>
                    <li>Mattheo de Wit</li>
                </ul>
        </div>
    );
}