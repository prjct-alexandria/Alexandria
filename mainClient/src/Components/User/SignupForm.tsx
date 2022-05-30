import * as React from "react";
import { Link } from "react-router-dom";

type SignupFormProps = {
  username: string | undefined;
  email: string | undefined;
  password: string | undefined;
  confirmPassword: string | undefined;
  onChangeUsername: (e: any) => void;
  onChangeEmail: (e: any) => void;
  onChangePassword: (e: any) => void;
  onChangeConfirmPassword: (e: any) => void;
  submitHandler: (e: any) => void;
};

export default function SignupForm(props: SignupFormProps) {
  return (
    <div className="wrapper">
      <form onSubmit={props.submitHandler} className="col-6 m-auto">
        <div className="mb-3">
          <label className="form-label">Name</label>
          <input
            onChange={props.onChangeUsername}
            name="username"
            type="username"
            className="form-control"
          />
        </div>
        <div className="mb-3">
          <label className="form-label">Email address</label>
          <input
            onChange={props.onChangeEmail}
            type="email"
            className="form-control"
            name="email"
          />
        </div>
        <div className="mb-3">
          <label className="form-label">Password</label>
          <input
            onChange={props.onChangePassword}
            type="password"
            className="form-control"
            name="password"
          />
        </div>
        <div className="mb-3">
          <label className="form-label">Confirm Password</label>
          <input
            onChange={props.onChangeConfirmPassword}
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
          Already have an account?&nbsp;
          <Link to="/login">Log in here</Link>!
        </p>
      </form>
    </div>
  );
}
