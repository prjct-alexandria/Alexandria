import * as React from "react";

export default function LoadingSpinner() {
  return (
    <div className="d-flex justify-content-center" data-testid="loadingSpinner">
      <div className="spinner-border text-primary" role="status">
        <span className="visually-hidden">Loading...</span>
      </div>
    </div>
  );
}
