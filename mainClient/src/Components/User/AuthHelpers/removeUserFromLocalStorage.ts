export default function removeUserFromLocalStorage() {
  // Remove login user details from local storage
  localStorage.removeItem("loggedUserEmail")
  localStorage.removeItem("loggedUsername")
  localStorage.removeItem("isLoggedIn")

  // Fire userAccountEvent to notify listeners that localstorage changes
  window.dispatchEvent(new Event("userAccountEvent"));
}
