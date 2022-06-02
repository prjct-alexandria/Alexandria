import * as React from "react";
import { fireEvent, render } from "@testing-library/react";
import LoginForm from "../../src/Components/User/LoginForm";

type LoginFormProps = {
  email: string | undefined;
  password: string | undefined;
  onChangeEmail: (e: any) => void;
  onChangePassword: (e: any) => void;
  submitHandler: (e: any) => void;
};

function renderLoginForm(props: Partial<LoginFormProps> = {}) {
  const defaultProps: LoginFormProps = {
    email: "",
    password: "",
    onChangePassword() {
      return;
    },
    onChangeEmail() {
      return;
    },
    submitHandler() {
      return;
    },
  };
  return render(<LoginForm {...defaultProps} {...props} />);
}

describe("<LoginForm />", () => {
  test("Display a blank login form", async () => {
    // Asynchroniously find element by given attribute
    const { findByTestId } = renderLoginForm();
    const loginForm = await findByTestId("login-form");

    // Expect form to have default values
    expect(loginForm).toHaveFormValues({
      email: "",
      password: "",
    });
  });

  test("Enter an email", async () => {
    // Mock the onChangeEmail
    const onChangeEmail = jest.fn();

    // Asynchroniously find element by given attribute
    const { findByTestId } = renderLoginForm({ onChangeEmail });
    const email = await findByTestId("email");

    // Modify form state to have the given value in the email field
    fireEvent.change(email, { target: { value: "user@gmail.com" } });

    // Expect the onChangeEmail to have been called with the given value
    expect(onChangeEmail).toHaveBeenCalledWith("email");
  });

  test("Enter a password", async () => {
    // Mock the onChangePassword
    const onChangePassword = jest.fn();

    // Asynchroniously find element by given attribute
    const { findByTestId } = renderLoginForm({ onChangePassword });
    const password = await findByTestId("password");

    // Modify form state to have the given value in the password field
    fireEvent.change(password, { target: { value: "password" } });

    // Expect the onChangePassword to have been called with the given value
    expect(onChangePassword).toHaveBeenCalledWith("password");
  });
});

test("Submit form", async () => {
  // Mock the submitHandler
  const submitHandler = jest.fn();

  // Asynchroniously find elements by given attribute
  const { findByTestId } = renderLoginForm({
    submitHandler,
  });

  const email = await findByTestId("email");
  const password = await findByTestId("password");
  const submit = await findByTestId("submit");

  // Modify form state to have the given value in the email field
  fireEvent.change(email, { target: { value: "user@gmail.com" } });

  // Modify form state to have the given value in the password field
  fireEvent.change(password, { target: { value: "password" } });

  // Modify form state to click the submit button
  fireEvent.click(submit);

  // Expect the submitHandler to have been called with the given values
  expect(submitHandler).toHaveBeenCalledWith("user@gmail.com", "password");
});
