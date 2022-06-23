export default function setUserInLocalStorage(email: string, name: string) {
  // Set user login details in local storage
  localStorage.setItem("loggedUserEmail", email);
  localStorage.setItem("loggedUsername", name);
  localStorage.setItem("isLoggedIn", "true");

  // Fire userAccountEvent to notify listeners that localstorage changes
  window.dispatchEvent(new Event("userAccountEvent"));
}
