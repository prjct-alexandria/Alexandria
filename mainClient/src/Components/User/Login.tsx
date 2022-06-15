import * as React from "react";
import { useState } from "react";
import LoginForm from "./LoginForm";
import configData from "../../config.json";
import NotificationAlert from "../NotificationAlert";
import $ from "jquery";
import setUserInLocalStorage from "./AuthHelpers/setUserInLocalStorage";

export default function Login() {
  let [email, setEmail] = useState<string>("");
  let [password, setPassword] = useState<string>("");
  let [showSuccessMessage, setShowSuccessMessage] = useState<boolean>(false);
  let [error, setError] = useState<Error>();
  let [username, setUsername] = useState<string>("");

  const submitHandler = (e: { preventDefault: () => void }) => {
    // Prevent unwanted default browser behavior
    e.preventDefault();

    const url = configData.back_end_url +"/login";

    // Construct request body
    const body = {
      email: email,
      pwd: password,
    };

    // Send request with a JSON of the user data
    fetch(url, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      mode: "cors",
      credentials: "include",
      body: JSON.stringify(body),
    }).then(
      async (response) => {
        if (response.ok) {
          setError(undefined);
          let responseJSON: {
            name: string;
            email: string;
          } = await response.json();

          // Set session details in localStorage
          setUserInLocalStorage(email, responseJSON.name);

          // Close login form
          $("#btn-close-login-form").trigger("click");

          // Show success alert for 3 seconds
          setUsername(responseJSON.name);
          setShowSuccessMessage(true);
          setTimeout(() => setShowSuccessMessage(false), 3000);
        } else {
          // Set error with message returned from the server
          let responseJSON: {
            message: string;
          } = await response.json();

          let serverMessage: string = responseJSON.message;
          setError(new Error(serverMessage));
        }
      },
      (error) => {
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
      {showSuccessMessage && (
        <NotificationAlert
          errorType="success"
          title={"Welcome, " + username + "!"}
          message={"You are now logged in."}
        />
      )}
      {
        <LoginForm
          email={email}
          password={password}
          onChangeEmail={(e) => {
            setEmail(e.target.value);
          }}
          onChangePassword={(e) => {
            setPassword(e.target.value);
          }}
          submitHandler={submitHandler}
        />
      }
    </div>
  );
}
