import React from "react";
import "../../App.js";
import LoginForm from "./LoginForm";

export default class Login extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      error: null,

      // Keep form data in component state
      user: {
        email: "",
        password: "",
      },
    };

    // Bind handler functions to component
    this.submitHandler = this.submitHandler.bind(this);
    this.changeHandler = this.changeHandler.bind(this);
  }

  // Update state on a change in one of the form fields
  // Keep previous state, only replace updated elements
  changeHandler = (event) => {
    this.setState((state) => ({
      user: { ...state.user, [event.target.name]: event.target.value },
    }));
  };

  // Send an HTTP POST request to /login with user info
  submitHandler = (event) => {
    // Prevent unwanted default browser behavior
    event.preventDefault();

    const url = "http://localhost:8080/login";
    const user = this.state.user;

    // Construct request body from state.user
    const body = {
      email: user.email,
      pwd: user.password,
    };

    // Send request with a JSON of the user data
    fetch(url, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      mode: "cors",
      body: JSON.stringify(body),
    }).then(
      // Success; set response in state
      (result) => {
        console.log("Success:", result);
        this.setState({
          isLoaded: true,
          items: result,
        });

        // Redirect to homepage; Comment this to debug the form submission
        if (typeof window !== "undefined") {
          window.location.href = "http://localhost:3000/";
        }
      },
      (error) => {
        // Request returns an error; set it in component's state
        console.error("Error:", error);
        this.setState({
          isLoaded: true,
          error,
        });
      }
    );
  };

  render() {
    const { error } = this.state;

    if (error) {
      // Render error message whenever state.error is set by HTTP response
      return (
        <div className="alert alert-danger" role="alert">
          Error: {error.message}. Please try again.
        </div>
      );
    } else {
      // Render LoginForm and pass state and handlers to its props
      return (
        <LoginForm
          state={this.state}
          changeHandler={this.changeHandler}
          submitHandler={this.submitHandler}
        />
      );
    }
  }
}
