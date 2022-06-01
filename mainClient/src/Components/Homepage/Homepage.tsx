import * as React from "react";
import HomepageHeader from "./HomepageHeader";

export default function Homepage() {
  return (
    <div>
      <HomepageHeader />
      <div className="container px-4 py-5" id="hanging-icons">
        <h2 className="pb-2 border-bottom">Hanging icons</h2>
        <div className="row g-4 py-5 row-cols-1 row-cols-lg-3">
          <div className="col d-flex align-items-start">
            <div className="icon-square bg-light text-dark flex-shrink-0 me-3"></div>
            <div>
              <h2>Featured title</h2>
              <p>
                Paragraph of text beneath the heading to explain the heading.
                We'll add onto it with another sentence and probably just keep
                going until we run out of words.
              </p>
              <a href="#" className="btn btn-dark">
                Button
              </a>
            </div>
          </div>
          <div className="col d-flex align-items-start">
            <div className="icon-square bg-light text-dark flex-shrink-0 me-3"></div>
            <div>
              <h2>Featured title</h2>
              <p>
                Paragraph of text beneath the heading to explain the heading.
                We'll add onto it with another sentence and probably just keep
                going until we run out of words.
              </p>
              <a href="#" className="btn btn-dark">
                Button
              </a>
            </div>
          </div>
          <div className="col d-flex align-items-start">
            <div className="icon-square bg-light text-dark flex-shrink-0 me-3"></div>
            <div>
              <h2>Featured title</h2>
              <p>
                Paragraph of text beneath the heading to explain the heading.
                We'll add onto it with another sentence and probably just keep
                going until we run out of words.
              </p>
              <a href="#" className="btn btn-dark">
                Button
              </a>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
