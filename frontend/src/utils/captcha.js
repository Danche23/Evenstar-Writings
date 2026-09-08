// 阿里云验证码 2.0 前端配置（从 .env.development.local 读取，非密钥）
// 与后端 configs/config.yaml 的 captcha.identity_prefix / captcha.scene_id 保持一致

export const CAPTCHA_PREFIX = import.meta.env.VITE_CAPTCHA_PREFIX || ''
export const CAPTCHA_SCENE_ID = import.meta.env.VITE_CAPTCHA_SCENE_ID || ''
export const CAPTCHA_REGION = 'cn'
