package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Category struct {
	ent.Schema
}

func (Category) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "categories"},
	}
}

func (Category) Fields() []ent.Field {
	return []ent.Field{
		field.Uint32("id"),
		field.String("name").MaxLen(64).Unique(),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

func (Category) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("articles", Article.Type).Ref("categories"),
	}
}
