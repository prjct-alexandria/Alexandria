import React from 'react';
import '../../App.js';

export default class FileUpload extends React.Component {

    constructor(props) {
        super(props);
        this.state = {
            error: null,
            selectedFile: null
        };
        this.uploadFileHandler = this.uploadFileHandler.bind(this);
        this.changeHandler = this.changeHandler.bind(this);
    }

    changeHandler = (event) => {
        this.setState({
            selectedFile: event.target.files[0]
        });
    }

    uploadFileHandler = () => {
        const formData = new FormData();
        formData.append('file', this.state.selectedFile);
        // TODO: server url in config file
        const url = 'http://localhost:8080/articles/1/versions/1';
        fetch(
            url,
            {
                method: 'POST',
                body: formData,
            }
        )
            .then(res => res.json())
            .then(
                (result) => {
                    console.log('Success:', result);
                    this.setState({
                        isLoaded: true,
                        items: result
                    });
                },
                (error) => {
                    console.error('Error:', error);
                    this.setState({
                        isLoaded: true,
                        error
                    });
                },
            );
    };

    render() {
        const {error} = this.state;

        if (error) {
            return <div>Error: {error.message}</div>;
        } else {
            return (
                <div className="col-6 align-content-center m-auto">
                    <h3>Upload a file</h3>
                    <hr/>

                    <input type="file" onChange={this.changeHandler}/>
                    <button onClick={this.uploadFileHandler}>Upload</button>
                </div>
            )
        }
    }
}