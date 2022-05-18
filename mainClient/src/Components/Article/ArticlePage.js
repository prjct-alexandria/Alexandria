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
                            <div>
                                <ul key={key}>
                                    {data.authors.map(a => (<li>{a}</li>))}
                                </ul>
                                <h1 key={key}>
                                    {data.title}
                                </h1>
                                <div key={key}>
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