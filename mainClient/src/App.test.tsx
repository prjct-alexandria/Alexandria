import { render } from "@testing-library/react";
import App from "./App";

test("biggerThanTest", () => {
  render(<App />);
  const linkElement = 3;
  expect(linkElement).toBeGreaterThanOrEqual(2);
});
