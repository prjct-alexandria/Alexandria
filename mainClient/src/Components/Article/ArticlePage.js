import React from 'react';
import '../../App.js';

export default class ArticlePage extends React.Component {

    constructor(props) {
        super(props);
        this.state = {
            error: null,
            isLoaded: false,
            version: null
        };
    }

    componentDidMount() {
        const articleId = window.location.href.split("/")[4];
        const versionId = window.location.href.split("/")[6];
        const url = "http://localhost:8080/articles/" + articleId + "/versions/" + versionId;

        fetch(
            url, {
                headers: {
                    'Accept': 'application/json',
                }
            }
        )
            .then(res => res.json())
            .then(
                (result) => {
                    this.setState({
                        isLoaded: true,
                        version: result
                    });
                },
                (error) => {
                    this.setState({
                        isLoaded: true,
                        error
                    });
                },
            )

    }

    render() {
        const { error, isLoaded, version } = this.state;
        if(isLoaded) {
            return (
                <>
                    <div className="article">
                        <ul key={1}>
                            {version.owners.map((a, i) => (<li key={i}>{a}</li>))}
                        </ul>
                        <h1 key={2}>
                            {version.title}
                        </h1>
                        <div key={3}>
                            {version.content}
                        </div>
                    </div>
                </>
            );
        }
    };
}