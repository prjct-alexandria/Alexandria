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
    <div
      className="modal fade"
      id="login"
      data-bs-backdrop="static"
      data-bs-keyboard="false"
      aria-labelledby="loginLabel"
      aria-hidden="true"
    >
      <div className="modal-dialog">
        <div className="modal-content">
          <div className="modal-header">
            <h5 className="modal-title" id="loginLabel">
              Login
            </h5>

            <button
              type="button"
              className="btn-close"
              data-bs-dismiss="modal"
              aria-label="Close"
            ></button>
          </div>
          <form onSubmit={props.submitHandler}>
            <div className="modal-body">
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
            </div>
            <div className="modal-footer">
              <button type="submit" className="btn btn-primary">
                Login
              </button>
              <div>
                {" "}
                Don't have an account yet?&nbsp;
                <a href="/" data-bs-toggle="modal" data-bs-target="#signUp">
                  Sign up here
                </a>
                !
              </div>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
}
