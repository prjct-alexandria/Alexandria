import React from 'react';
import '../../App.js';

export default class Homepage extends React.Component {
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
            <div>Homepage</div>
        )
    }
}