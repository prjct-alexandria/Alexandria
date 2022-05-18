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
        //const url = "../article.json"; // Placeholder
        const articleId = window.location.href.split("/")[4];
        const versionId = window.location.href.split("/")[6];
        const url = "localhost:8080/articles/" + articleId + "/version/" + versionId;

        fetch(
            url, {headers: {
                    'Content-Type': 'application/json',
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
                                    {data.authors.map((a, i) => (<li key={i}>{a}</li>))}
                                </ul>
                                <h1 key={2}>
                                    {data.title}
                                </h1>
                                <div key={3}>
                                    {data.content}
                                </div>
                            </div>
                        );
                    })}
                </div>
            </>
        );
    };
}