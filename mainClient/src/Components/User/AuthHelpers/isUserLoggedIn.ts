import removeUserFromLocalStorage from "./removeUserFromLocalStorage";

// Source: https://stackoverflow.com/a/52406518/14209629
export function getCookie(name: string) {
  let match = document.cookie.match(RegExp('(?:^|;\\s*)' + name + '=([^;]*)'));
  return match ? match[1] : null;
}

export function authCookieCheck(): boolean {
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
  return authCookieCheck() && localStorage.getItem("isLoggedIn") === "true";
}
