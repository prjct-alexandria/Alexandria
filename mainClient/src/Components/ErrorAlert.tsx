type alertProps = {
  title: string;
  message: string;
};

export default function ErrorAlert(props: alertProps) {
  return (
    <div
      id="error-alert"
      className="alert alert-danger alert-dismissible fade show"
      role="alert"
    >
      <strong>{props.title}</strong> {props.message}
      <button
        type="button"
        className="btn-close"
        data-bs-dismiss="alert"
        aria-label="Close"
      ></button>
    </div>
  );
}
