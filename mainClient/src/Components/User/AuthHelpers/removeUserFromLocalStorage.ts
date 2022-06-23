export default function removeUserFromLocalStorage() {
  // Remove login user details from local storage
  localStorage.setItem("loggedUserEmail", "");
  localStorage.setItem("loggedUsername", "");
  localStorage.setItem("isLoggedIn", "false");

  // Fire userAccountEvent to notify listeners that localstorage changes
  window.dispatchEvent(new Event("userAccountEvent"));
}
