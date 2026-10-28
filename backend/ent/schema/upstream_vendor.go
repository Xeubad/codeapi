package schema

import (
	"github.com/Wei-Shaw/sub2api/ent/schema/mixins"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// UpstreamVendor 上游厂家：同一厂家的多把 key（条目）共享账号级并发预算。
// 条目归属通过 accounts.extra.upstream_vendor_id 关联，不建边以保持账号表零改动。
type UpstreamVendor struct {
	ent.Schema
}

func (UpstreamVendor) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "upstream_vendors"},
	}
}

func (UpstreamVendor) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.TimeMixin{},
		mixins.SoftDeleteMixin{},
	}
}

func (UpstreamVendor) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			MaxLen(100).
			NotEmpty(),
		field.Int("concurrency_budget").
			Default(0).
			Min(0).
			Comment("账号级并发预算；0 表示不启用共享限制"),
		field.String("notes").
			Optional().
			Nillable().
			SchemaType(map[string]string{dialect.Postgres: "text"}),
		// base_urls 是厂家的候选接入地址；条目保存时会逐个测速并只固化最快的一个，
		// 运行时不再切换（避免请求途中换线路）。
		field.JSON("base_urls", []string{}).
			Optional().
			SchemaType(map[string]string{dialect.Postgres: "jsonb"}),
	}
}

func (UpstreamVendor) Indexes() []ent.Index {
	return []ent.Index{
		// 部分唯一索引：软删除后可重用名称（与 groups 同型）。
		index.Fields("name").
			Unique().
			StorageKey("idx_upstream_vendors_name_active").
			Annotations(entsql.IndexWhere("deleted_at IS NULL")),
		index.Fields("deleted_at"),
	}
}
