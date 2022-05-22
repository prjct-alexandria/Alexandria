import React from 'react';
import '../../App.js';

export default class ArticlePage extends React.Component {

    constructor(props) {
        super(props);
        this.state = {
            error: null,
            isLoaded: false,
            items: []
        };
    }

    componentDidMount() {
        const articleId = window.location.href.split("/")[4];
        const versionId = window.location.href.split("/")[6];
        const url = "http://localhost:8080/articles/" + articleId + "/versions/" + versionId;

        fetch(
            url, {
                headers: {
                    'Accept': 'application/json'
                }
            }
        )
            .then(res => res.json())
            .then(
                (result) => {
                    this.setState({
                        isLoaded: true,
                        items: result
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
        const { error, isLoaded, items } = this.state;
        return (
            <>
                <div className="article">
                    {items.map((data, key) => {
                        return (
                            <div key={key}>
                                <ul key={1}>
                                    {data.Authors.map((a, i) => (<li key={i}>{a}</li>))}
                                </ul>
                                <h1 key={2}>
                                    {data.Title}
                                </h1>
                                <div key={3}>
                                    {data.Content}
                                </div>
                            </div>
                        );
                    })}
                </div>
            </>
        );
    };
}