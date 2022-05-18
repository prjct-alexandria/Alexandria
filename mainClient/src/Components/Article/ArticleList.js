import React from 'react';
import '../../App.js';
import ArticleListElement from "./ArticleListElement";

export default class ArticleList extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            error: null,
            isLoaded: false,
            items: []
        };
    }

    componentDidMount() {
        //"http://localhost:8080/helloWorldJson"
        const url = 'articleList.json'; // Placeholder
        // const url = '/articles' // should be this url, but used the one above for demonstration
        fetch(url
            )
            .then(res => res.json())
            .then(
                (result) => {
                    this.setState({
                        isLoaded: true,
                        items: result
                    });
                },
                // Note: it's important to handle errors here
                // instead of a catch() block so that we don't swallow
                // exceptions from actual bugs in components.
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
        if (error) {
            return <div>Error: {error.message}</div>;
        } else if (!isLoaded) {
            return <div>Loading...</div>;
        } else {
            return (
                items.map(function(article, i){
                    return (<ArticleListElement key={i} article={article} />);
                })
            );
        }
    }

}