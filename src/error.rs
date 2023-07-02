use axum::{http::StatusCode, response::IntoResponse};

#[derive(Debug)]
pub enum Error {
    ServiceDown,
    InternalServer,
    Timeout,
    Duplicated,
    BadRequest,
    Unauthorized,
    Forbidden,
    NotFound,
    WithMessage(StatusCode, String),
}

impl std::fmt::Display for Error {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        f.write_str("{self:?}")
    }
}

impl std::error::Error for Error {}

impl From<tonic::Status> for Error {
    fn from(value: tonic::Status) -> Self {
        use Error::*;

        match value.code() {
            tonic::Code::Ok => unreachable!(),
            tonic::Code::Cancelled => ServiceDown,
            tonic::Code::Unknown => ServiceDown,
            tonic::Code::InvalidArgument => BadRequest,
            tonic::Code::DeadlineExceeded => Timeout,
            tonic::Code::NotFound => NotFound,
            tonic::Code::AlreadyExists => Duplicated,
            tonic::Code::PermissionDenied => Forbidden,
            tonic::Code::ResourceExhausted => BadRequest,
            tonic::Code::FailedPrecondition => BadRequest,
            tonic::Code::Aborted => ServiceDown,
            tonic::Code::OutOfRange => BadRequest,
            tonic::Code::Unimplemented => BadRequest,
            tonic::Code::Internal => InternalServer,
            tonic::Code::Unavailable => ServiceDown,
            tonic::Code::DataLoss => InternalServer,
            tonic::Code::Unauthenticated => Unauthorized,
        }
    }
}

impl From<axum::extract::multipart::MultipartError> for Error {
    fn from(value: axum::extract::multipart::MultipartError) -> Self {
        Error::WithMessage(value.status(), value.body_text())
    }
}

impl IntoResponse for Error {
    fn into_response(self) -> axum::response::Response {
        match self {
            Error::ServiceDown => (StatusCode::SERVICE_UNAVAILABLE, "Service down").into_response(),
            Error::InternalServer => {
                (StatusCode::INTERNAL_SERVER_ERROR, "Internal server error").into_response()
            }
            Error::Timeout => (StatusCode::REQUEST_TIMEOUT, "Request time out").into_response(),
            Error::Duplicated => (StatusCode::CONFLICT, "Duplicated").into_response(),
            Error::BadRequest => (StatusCode::BAD_REQUEST, "Bad request").into_response(),
            Error::Unauthorized => (StatusCode::UNAUTHORIZED, "Unauthorized").into_response(),
            Error::Forbidden => (StatusCode::FORBIDDEN, "Forbidden").into_response(),
            Error::NotFound => (StatusCode::NOT_FOUND, "Not found").into_response(),
            Error::WithMessage(code, text) => (code, text).into_response(),
        }
    }
}
