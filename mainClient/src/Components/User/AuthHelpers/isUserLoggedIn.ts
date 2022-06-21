export default function isUserLoggedIn(): boolean {
  return localStorage.getItem("isLoggedIn") === "true";
}
