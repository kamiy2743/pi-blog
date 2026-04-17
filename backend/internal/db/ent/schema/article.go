package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Article struct {
	ent.Schema
}

func (Article) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "articles"},
	}
}

func (Article) Fields() []ent.Field {
	return []ent.Field{
		field.Uint32("id"),
		field.String("title").MaxLen(255),
		field.Text("body"),
		field.Bool("is_published").Default(false),
		field.Time("publish_start_at").Optional().Nillable(),
		field.Time("publish_end_at").Optional().Nillable(),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

func (Article) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("categories", Category.Type),
		edge.To("comments", Comment.Type),
	}
}
