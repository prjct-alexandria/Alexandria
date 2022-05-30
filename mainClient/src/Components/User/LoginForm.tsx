import * as React from "react";
import { Link } from "react-router-dom";

type LoginFormProps = {
  email: string | undefined;
  password: string | undefined;
  onChangeEmail: (e: any) => void;
  onChangePassword: (e: any) => void;
  submitHandler: (e: any) => void;
};

export default function LoginForm(props: LoginFormProps) {
  return (
    <div className="wrapper">
      <form onSubmit={props.submitHandler} className="col-6 m-auto">
        <div className="mb-3">
          <label className="form-label">Email address</label>
          <input
            onChange={props.onChangeEmail}
            name="email"
            type="email"
            className="form-control"
          />
        </div>
        <div className="mb-3">
          <label className="form-label">Password</label>
          <input
            onChange={props.onChangePassword}
            name="password"
            type="password"
            className="form-control"
          />
        </div>

        <button type="submit" className="btn btn-primary">
          Login
        </button>
        <p>
          {" "}
          Don't have an account yet?&nbsp;
          <Link to="/signup">Sign up here</Link>!
        </p>
      </form>
    </div>
  );
}
