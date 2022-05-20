import React from "react";
import "../../App.js";
import LoginForm from "./LoginForm";

export default class Login extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      error: null,
      user: {
        email: "",
        password: "",
      },
    };
    this.submitHandler = this.submitHandler.bind(this);
    this.changeHandler = this.changeHandler.bind(this);
  }

  submitHandler = (event) => {
    event.preventDefault();

    const url = "http://localhost:8080/login";
    const user = this.state.user;
    const body = {
      email: user.email,
      pwd: user.password,
    };

    fetch(url, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      mode: "cors",
      body: JSON.stringify(body),
    }).then(
      (result) => {
        console.log("Success:", result);
        this.setState({
          isLoaded: true,
          items: result,
        });
        if (typeof window !== "undefined") {
          window.location.href = "http://localhost:3000/";
        }
      },
      (error) => {
        console.error("Error:", error);
        this.setState({
          isLoaded: true,
          error,
        });
      }
    );
  };

  changeHandler = (event) => {
    this.setState((state) => ({
      user: { ...state.user, [event.target.name]: event.target.value },
    }));
  };

  render() {
    const { error } = this.state;

    if (error) {
      return (
        <div className="alert alert-danger" role="alert">
          Error: {error.message}. Please try again.
        </div>
      );
    } else {
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
