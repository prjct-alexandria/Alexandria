type alertProps = {
  errorType: string; // success or danger
  title: string;
  message: string;
};

export default function NotificationAlert(props: alertProps) {
  return (
    <div
      id={props.errorType + "-alert"}
      className={
        "alert alert-" +
        props.errorType +
        " alert-dismissible notification-alert fade show "
      }
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
