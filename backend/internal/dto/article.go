package dto

// 文章模块 DTO
//
// 建议包含的 DTO（字段参考 docs/openapi.yaml 的 article 相关接口）：
//   - CreateArticleRequest  创建文章（title / summary / content / cover / status / 分类 / 标签）
//   - UpdateArticleRequest  更新文章
//   - ArticleListRequest    文章列表查询（分页 + 分类/标签/关键字筛选）
//   - ArticleListResponse   文章列表响应
//   - ArticleDetailResponse 文章详情响应（含作者、分类、标签、浏览量）
//   - AdminArticleListResponse 后台文章列表（含草稿，可展示 status）
