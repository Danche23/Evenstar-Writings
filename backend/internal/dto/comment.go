package dto

// 评论模块 DTO
//
// 建议包含的 DTO（字段参考 docs/openapi.yaml 的 comment 相关接口）：
//   - CreateCommentRequest 发表评论（article_id / parent_id / reply_to_id / content，content ≤ 400 字）
//   - CommentListRequest   评论列表查询（分页）
//   - CommentItemResponse  单条评论（含用户昵称、等级：一级/二级、回复对象「张三 → 李四」）
//   - CommentListResponse  评论列表响应
//   - AdminCommentListResponse 后台评论管理列表（含置顶、等级列）
