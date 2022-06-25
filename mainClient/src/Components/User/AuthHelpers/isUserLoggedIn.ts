import removeUserFromLocalStorage from "./removeUserFromLocalStorage";

function getCookie(name: string) {
  let match = document.cookie.match(RegExp('(?:^|;\\s*)' + name + '=([^;]*)'));
  return match ? match[1] : null;
}

function authCookiePresent(): boolean {
  let present = getCookie("isAuthorized")
  if (present == null) {
    if (localStorage.getItem("isLoggedIn") === "true") {
      removeUserFromLocalStorage()

      //TODO: Show pop-up user has been logged out due to inactivity
    }
    return false
  }
  return true
}

export default function isUserLoggedIn(): boolean {
  return authCookiePresent() && localStorage.getItem("isLoggedIn") === "true";
}
