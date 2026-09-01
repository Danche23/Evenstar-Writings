package errors

import "net/http"

// HTTPStatus 业务错误码 → HTTP 状态码映射（REST 语义码）
// 约定：HTTP 状态码表达错误类别，body 里的 code 表达具体业务原因，前端拦截器两层处理。
func HTTPStatus(code int) int {
	switch {
	// 通用码：本身就是 HTTP 语义码
	case code == CodeBadRequest:
		return http.StatusBadRequest
	case code == CodeUnauthorized:
		return http.StatusUnauthorized
	case code == CodeForbidden:
		return http.StatusForbidden
	case code == CodeNotFound:
		return http.StatusNotFound
	case code == CodeMethodNotAllowed:
		return http.StatusMethodNotAllowed
	case code == CodeRequestTimeout:
		return http.StatusRequestTimeout
	case code == CodeConflict:
		return http.StatusConflict
	case code == CodeTooManyRequests:
		return http.StatusTooManyRequests

	// 服务器错误 5xx
	case code >= 500 && code < 600:
		return http.StatusInternalServerError

	// 业务错误 1xxx（用户类）
	case code == CodeUserNotFound:
		return http.StatusNotFound
	case code == CodeUserAlreadyExists:
		return http.StatusConflict
	case code == CodeInvalidCredentials:
		return http.StatusUnauthorized
	case code == CodeUserDisabled:
		return http.StatusForbidden
	case code == CodeInvalidToken, code == CodeTokenExpired:
		return http.StatusUnauthorized
	case code == CodeCaptchaRequired, code == CodeCaptchaVerifyFailed:
		return http.StatusBadRequest

	// 参数错误 2xxx
	case code >= 2000 && code < 3000:
		return http.StatusBadRequest

	// 资源错误 3xxx
	case code == CodeResourceNotFound:
		return http.StatusNotFound
	case code == CodeResourceAlreadyExists:
		return http.StatusConflict
	case code == CodeResourceLocked:
		return http.StatusConflict

	default:
		return http.StatusInternalServerError
	}
}
