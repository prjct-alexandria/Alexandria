import React from "react";
import "../../App.js";
import { Link } from "react-router-dom";

export default class LoginForm extends React.Component {
  render() {
    return (
      <div className="wrapper">
        <form
          onChange={this.props.changeHandler}
          onSubmit={this.props.submitHandler}
          className="col-6 m-auto"
        >
          <div className="mb-3">
            <label className="form-label">Email address</label>
            <input name="email" type="email" className="form-control" />
          </div>
          <div className="mb-3">
            <label className="form-label">Password</label>
            <input name="password" type="password" className="form-control" />
          </div>

          <button type="submit" className="btn btn-primary">
            Login
          </button>
          <p>
            {" "}
            Don't have an account yet? <br />
            <Link to="/signup">Sign up</Link>
          </p>
        </form>
      </div>
    );
  }
}
