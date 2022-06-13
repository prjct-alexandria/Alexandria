import * as React from "react";
import { useState } from "react";
import LoginForm from "./LoginForm";
import configData from "../../config.json";

export default function Login() {
  let [email, setEmail] = useState<string>("");
  let [password, setPassword] = useState<string>();
  let [error, setError] = useState(null);

  const onChangeEmail = (e: { target: { value: any } }) => {
    setEmail(e.target.value);
  };

  const onChangePassword = (e: { target: { value: any } }) => {
    setPassword(e.target.value);
  };

  const submitHandler = (e: { preventDefault: () => void }) => {
    // Prevent unwanted default browser behavior
    e.preventDefault();

    const url= configData.back_end_url +"/login";

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
      // Success; set response in state
      (response) => {
        console.log("Success:", response);

        if (response.ok) {
          localStorage.setItem("loggedUserEmail", email);

          // Redirect to homepage; Comment this to debug the form submission
          if (typeof window !== "undefined") {
            window.location.href = configData.front_end_url;
          } else {
            console.log("Error: Undefined window");
          }
        }
      },
      (error) => {
        // Request returns an error; set it in component's state
        console.error("Error:", error);
        setError(error);
      }
    );
  };

  return (
    <div>
      {error && <div>{`There is a problem - ${error}`}</div>}
      {
        <LoginForm
          email={email}
          password={password}
          onChangeEmail={onChangeEmail}
          onChangePassword={onChangePassword}
          submitHandler={submitHandler}
        />
      }
    </div>
  );
}
