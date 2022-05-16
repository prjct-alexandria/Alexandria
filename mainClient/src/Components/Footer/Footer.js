import React from 'react';
import '../../App.js';

export default class Footer extends React.Component {

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
            <footer className="mt-auto col-xs-12 d-flex flex-wrap justify-content-between align-items-center py-3 my-4 border-top">
                <span>Footer</span>
            </footer>
        )
    }
}