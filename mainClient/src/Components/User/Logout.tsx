import * as React from "react";
import { useState } from "react";
import NotificationAlert from "../NotificationAlert";
import removeUserFromLocalStorage from "./AuthHelpers/removeUserFromLocalStorage";
import configData from "../../config.json";

export default function Logout() {
  let [error, setError] = useState<Error>();
  let [showSuccessMessage, setShowSuccessMessage] = useState<boolean>(false);

  const url = configData.back_end_url + "/logout";

  const LogoutHandler = () => {
    fetch(url, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      mode: "cors",
      credentials: "include",
    }).then(
      async (response) => {
        if (response.ok) {
          setError(undefined);
          // Set session details in localStorage
          removeUserFromLocalStorage();

          // Set success in state to show success alert
          setShowSuccessMessage(true);

          // After 3s, remove success from state to hide success alert
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
        // Request returns an error; set it in component's state
        setError(error);
      }
    );
  };

  return (
    <>
      {error && (
        <NotificationAlert
          errorType="danger"
          title={"Something went wrong when logging out. "}
          message={"" + error}
        />
      )}
      {showSuccessMessage && (
        <NotificationAlert
          errorType="success"
          title="Logout successful! "
          message={"You have logged out of your account."}
        />
      )}
      <button
        role="button"
        type="button"
        className="dropdown-item"
        onClick={LogoutHandler}
      >
        Log out
      </button>
    </>
  );
}
