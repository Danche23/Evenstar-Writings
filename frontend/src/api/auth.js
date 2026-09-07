import request from './request'

// 认证接口：与 docs/openapi.yaml 的 Auth 模块一致

export function sendCode(email, type, captchaVerifyParam) {
  return request.post('/api/auth/send-code', {
    email,
    type, // register | reset
    captcha_verify_param: captchaVerifyParam
  })
}

export function register(data) {
  // { username, password, email, code, nickname?, captcha_verify_param }
  return request.post('/api/auth/register', data)
}

export function resetPassword(data) {
  // { email, code, new_password, captcha_verify_param }
  return request.post('/api/auth/reset-password', data)
}

export function login(email, password, captchaVerifyParam) {
  const payload = { email, password }
  if (captchaVerifyParam) payload.captcha_verify_param = captchaVerifyParam
  return request.post('/api/auth/login', payload)
}
