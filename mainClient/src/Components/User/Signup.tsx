import * as React from "react";
import { useState } from "react";
import $ from "jquery";
import SignupForm from "./SignupForm";
import configData from "../../config.json";
import NotificationAlert from "../NotificationAlert";

export default function Signup() {
  let [username, setUsername] = useState<string>();
  let [email, setEmail] = useState<string>();
  let [password, setPassword] = useState<string>();
  let [confirmPassword, setConfirmPassword] = useState<string>();
  let [error, setError] = useState<Error>();
  let [signupSuccess, setSignupSuccess] = useState<boolean>(false);

  const onChangeUsername = (e: { target: { value: any } }) => {
    setUsername(e.target.value);
  };

  const onChangeEmail = (e: { target: { value: any } }) => {
    setEmail(e.target.value);
  };

  const onChangePassword = (e: { target: { value: any } }) => {
    setPassword(e.target.value);
  };

  const onChangeConfirmPassword = (e: { target: { value: any } }) => {
    setConfirmPassword(e.target.value);
  };

  // Send an HTTP POST request to /register with user info
  const submitHandler = (e: { preventDefault: () => void }) => {
    // Prevent unwanted default browser behavior
    e.preventDefault();

    const url = configData.back_end_url +"/users";

    // Construct request body
    const body = {
      name: username,
      email: email,
      pwd: password,
    };

    // Send request with a JSON of the user data
    fetch(url, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      mode: "cors",
      body: JSON.stringify(body),
    }).then(
      (response) => {
        if (response.ok) {
          setError(undefined);
          setSignupSuccess(response.ok);
          // Use JQuery to "simulate" button presses,
          // which close the signup modal, then open the login
          $("#btn-close-signup-form").trigger("click");
          $("#btn-open-login-form").trigger("click");
        } else {
          setError(new Error(response.statusText));
        }
      },
      (error) => {
        // Request returns an error
        console.error("Error:", error);
        setError(error);
      }
    );
  };

  return (
    <div>
      {error && (
        <NotificationAlert
          errorType="danger"
          title={"Error: "}
          message={"Something went wrong. " + error}
        />
      )}
      {signupSuccess && (
        <NotificationAlert
          errorType="success"
          title="Account successfully created! "
          message={"You can now log into your account."}
        />
      )}
      {
        <SignupForm
          username={username}
          email={email}
          password={password}
          confirmPassword={confirmPassword}
          onChangeUsername={onChangeUsername}
          onChangeEmail={onChangeEmail}
          onChangePassword={onChangePassword}
          onChangeConfirmPassword={onChangeConfirmPassword}
          submitHandler={submitHandler}
        />
      }
    </div>
  );
}
