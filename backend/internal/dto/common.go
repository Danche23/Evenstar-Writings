package dto

// 通用 DTO（跨模块复用）
//
// 约定（重要）：
// 1. DTO 只负责 API 请求/响应的数据传输 + 参数校验（binding 标签），
//    不写业务逻辑，也不直接复用 entity（避免把表结构直接暴露给前端）。
// 2. 分页相关 DTO 放这里：PageRequest / PageResponse。
//    注意：pkg/response 里已有 PageRequest/PageResponse（属于「统一响应外壳」的一部分），
//    后续做分页接口时建议统一在 dto/common.go 定义分页 DTO，避免两处重复维护。
//
// 建议包含的 DTO：
//   - PageRequest   分页请求参数（page / page_size）
//   - PageResponse  分页响应结构（list / total / page / page_size / total_page）
//   - IDRequest     通用「按 id 操作」的请求体（可选，视需要）
