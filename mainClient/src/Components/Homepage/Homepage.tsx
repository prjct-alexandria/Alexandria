import * as React from "react";
import HomepageHeader from "./HomepageHeader";

export default function Homepage() {
  return (
    <div>
      <HomepageHeader />
      <div className="container px-4 py-5" id="hanging-icons">
        <h2 className="pb-2 border-bottom">
          Scientific publishing, but make it...
        </h2>
        <div className="row g-4 py-5 row-cols-1 row-cols-lg-3">
          <div className="col d-flex align-items-start">
            <div className="icon-square bg-light text-dark flex-shrink-0 me-3"></div>
            <div>
              <h2>Open</h2>
              <p>
                Have access to collaborative scientific publishing, without
                extra technical or coding skills. Read articles without an
                account, or easily create one.
              </p>
            </div>
          </div>
          <div className="col d-flex align-items-start">
            <div className="icon-square bg-light text-dark flex-shrink-0 me-3"></div>
            <div>
              <h2>Transparent</h2>
              <p>
                See the data used in the papers, so you can verify that the
                numbers and figures are meaningful and reliable. View the
                contributions of different authors.
              </p>
            </div>
          </div>
          <div className="col d-flex align-items-start">
            <div className="icon-square bg-light text-dark flex-shrink-0 me-3"></div>
            <div>
              <h2>Collaborative</h2>
              <p>
                Upload your paper, collaborate with other researchers, make your
                own version of an article & edit it, merge your work together.
                Comment & peer review.
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
