type alertProps = {
  title: string;
  message: string;
};

export default function SuccessAlert(props: alertProps) {
  return (
    <div
      id="success-alert"
      className="alert alert-success alert-dismissible fade show"
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
