import logo from './logo.svg';
import React from 'react';
import './App.css';

class Joke extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      liked: ""
    }
    this.like = this.like.bind(this);
  }

  like() {
    // ... we'll add this block later
  }

  render() {
    return (
        <div className="col-xs-4">
          <div className="panel panel-default">
            <div className="panel-heading">#{this.props.joke.id} <span className="pull-right">{this.state.liked}</span></div>
            <div className="panel-body">
              {this.props.joke.joke}
            </div>
            <div className="panel-footer">
              {this.props.joke.likes} Likes &nbsp;
              <a onClick={this.like} className="btn btn-default">
                <span className="glyphicon glyphicon-thumbs-up"></span>
              </a>
            </div>
          </div>
        </div>
    )
  }
}

class HelloWorld extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      error: null,
      isLoaded: false,
      items: []
    };
  }

  componentDidMount() {
    fetch("http://localhost:8080/helloWorldJson")
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
          this.state.items.map(function(joke, i){
              return (<Joke key={i} joke={joke} />);
            })
      );
    }
  }

}

function App() {
  return (
    <div className="App">
      <HelloWorld/>
    </div>
  );
}

export default App;
