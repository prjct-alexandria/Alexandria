import React from 'react';
import '../../App.js';

export default class Article extends React.Component {

    constructor(props) {
        super(props);
        this.state = {
            error: null,
            isLoaded: false,
            items: []
        };
    }

    render() {
        return (
            <div className="col-xs-12">
                <div className="panel panel-default">
                    <div className="panel-heading">#{this.props.article.id}</div>
                    <div className="panel-body">
                        {this.props.article.title}
                    </div>
                    <div className="panel-footer">
                        {this.props.article.author}
                    </div>
                </div>
            </div>
        )
    }
}