package dto

// 上传模块 DTO
//
// 建议包含的 DTO（字段参考 docs/openapi.yaml 的 upload 相关接口）：
//   - UploadResponse 上传成功后返回（url / filename / size / mime）
//   - 注：上传是 multipart/form-data，请求参数多为表单字段（scene 等），
//         通常不需要单独的结构体，直接在 handler 里取表单字段即可。
