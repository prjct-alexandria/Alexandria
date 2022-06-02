import * as React from "react";
import { fireEvent, render } from "@testing-library/react";
import LoginForm from "../../Components/User/LoginForm";

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
    //Arrange: Asynchronously find element by given attribute
    const { findByTestId } = renderLoginForm();
    const loginForm = await findByTestId("login-form");

    // Assert: Expect form to have default values
    expect(loginForm).toHaveFormValues({
      email: "",
      password: "",
    });
  });

  test("Enter an email", async () => {
    //Arrange
    // Mock the onChangeEmail
    const onChangeEmail = jest.fn();

    // Asynchronously find element by given attribute
    const { findByTestId } = renderLoginForm({ onChangeEmail });
    const email = await findByTestId("email");

    // Act: Modify form state to have the given value in the email field
    fireEvent.change(email, { target: { value: "user@gmail.com" } });

    // Assert: Expect the onChangeEmail to have been called with the given value
    expect(onChangeEmail).toHaveBeenCalled();
    expect(email).toHaveValue("user@gmail.com");
  });

  test("Enter a password", async () => {
    // Arrange: Mock the onChangePassword
    const onChangePassword = jest.fn();

    // Asynchronously find element by given attribute
    const { findByTestId } = renderLoginForm({ onChangePassword });
    const password = await findByTestId("password");

    // Act: Modify form state to have the given value in the password field
    fireEvent.change(password, { target: { value: "password" } });

    // Assert: Expect the onChangePassword to have been called with the given value
    expect(onChangePassword).toHaveBeenCalled();
    expect(password).toHaveValue("password");
  });

  test("Submit form", async () => {
    // Arrange: Mock the submitHandler
    const submitHandler = jest.fn((e) => e.preventDefault());

    // Asynchronously find elements by given attribute
    const { findByTestId } = renderLoginForm({
      submitHandler,
    });

    const email = await findByTestId("email");
    const password = await findByTestId("password");
    const submit = await findByTestId("submit");

    //Act
    // Modify form state to have the given value in the email field
    fireEvent.change(email, { target: { value: "user@gmail.com" } });

    // Modify form state to have the given value in the password field
    fireEvent.change(password, { target: { value: "password" } });

    // Modify form state to click the submit button
    fireEvent.click(submit);

    // Assert: Expect the submitHandler to have been called with the given values
    expect(submitHandler).toHaveBeenCalled();
    expect(email).toHaveValue("user@gmail.com");
    expect(password).toHaveValue("password");
  });
});
