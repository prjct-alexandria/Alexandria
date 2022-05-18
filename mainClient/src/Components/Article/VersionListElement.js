import React from 'react';
import '../../App.js';
import {Link} from "react-router-dom";

export default class VersionListElement extends React.Component {
    render() {
        return (
            <div className='row row-no-gutters col-md-12 text-wrap'>
                <div className='col-md-9 text-start'>
                    <Link to={"/articles/" + this.props.aId + "/versions/" + this.props.version.id}>
                        {this.props.version.title}
                    </Link>
                </div>
                <div className='col-md-1'>
                    {this.props.version.author}
                </div>
                <div className='col-md-1'>
                    {this.props.version.date_created}
                </div>
                <div className='col-md-1'>
                    {this.props.version.status}
                </div>
            </div>
        )
    }
}