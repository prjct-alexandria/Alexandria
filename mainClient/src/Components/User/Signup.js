import React from "react";
import "../../App.js";
import SignupForm from "./SignupForm";

export default class Signup extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      error: null,
      user: {
        username: "",
        email: "",
        password: "",
        passwordConfirm: "",
      },
    };
    this.submitHandler = this.submitHandler.bind(this);
    this.changeHandler = this.changeHandler.bind(this);
  }

  changeHandler = (event) => {
    this.setState((state) => ({
      user: { ...state.user, [event.target.name]: event.target.value },
    }));
  };

  submitHandler = (event) => {
    event.preventDefault();

    // const url = "http://localhost:8080/register";
    const url = "http://localhost:3000";
    const user = this.state.user;
    const body = {
      name: user.username,
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
        // Uncomment when finished debugging
        // if (typeof window !== "undefined") {
        //   window.location.href = "http://localhost:3000/login";
        // }
      },
      (error) => {
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
      return (
        <div className="alert alert-danger" role="alert">
          Error: {error.message}. Please try again.
        </div>
      );
    } else {
      return (
        <div className="wrapper">
          <SignupForm
            state={this.state}
            changeHandler={this.changeHandler}
            submitHandler={this.submitHandler}
          />
        </div>
      );
    }
  }
}
