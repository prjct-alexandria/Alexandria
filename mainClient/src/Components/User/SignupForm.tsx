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
    <div
      className="modal fade"
      id="signUp"
      data-bs-backdrop="static"
      data-bs-keyboard="false"
      aria-labelledby="signUpLabel"
      aria-hidden="true"
    >
      <div className="modal-dialog">
        <div className="modal-content">
          <div className="modal-header">
            <h5 className="modal-title" id="signUpLabel">
              Sign up
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
            </div>
            <div className="modal-footer">
              <button type="submit" className="btn btn-primary">
                Sign up
              </button>
              <div>
                {" "}
                Already have an account?&nbsp;
                <a href="/" data-bs-toggle="modal" data-bs-target="#login">
                  Log in here
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
