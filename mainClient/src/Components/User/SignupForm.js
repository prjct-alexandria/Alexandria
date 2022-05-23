import React from "react";
import "../../App.js";
import { Link } from "react-router-dom";

export default class SignupForm extends React.Component {
  render() {
    return (
      <div className="wrapper">
        <form
          onChange={this.props.changeHandler}
          onSubmit={this.props.submitHandler}
          className="col-6 m-auto"
        >
          <div className="mb-3">
            <label className="form-label">Name</label>
            <input name="username" type="username" className="form-control" />
          </div>
          <div className="mb-3">
            <label className="form-label">Email address</label>
            <input type="email" className="form-control" name="email" />
          </div>
          <div className="mb-3">
            <label className="form-label">Password</label>
            <input type="password" className="form-control" name="password" />
          </div>
          <div className="mb-3">
            <label className="form-label">Confirm Password</label>
            <input
              type="password"
              className="form-control"
              name="passwordConfirm"
            />
          </div>

          <button type="submit" className="btn btn-primary">
            Submit
          </button>
          <p>
            {" "}
            Already have an account? <br />
            <Link to="/login">Log in</Link>
          </p>
        </form>
      </div>
    );
  }
}
