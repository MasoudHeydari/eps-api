package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// SearchQuery holds the schema definition for the SearchQueries entity.
type SearchQuery struct {
	ent.Schema
}

// Annotations of the SearchQuery.
func (SearchQuery) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "search_queries"},
	}
}

// Fields of the SearchQuery.
func (SearchQuery) Fields() []ent.Field {
	return []ent.Field{
		field.String("query").NotEmpty(),
		field.Int("loc_code"),
		field.String("language"),
		field.UUID("job_id", uuid.New()).Optional(),
		field.Bool("is_finished").Default(false),
		field.Time("created_at").SchemaType(TimeStampWithTZ).Default(time.Now),
	}
}

// Edges of the SearchQuery.
func (SearchQuery) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("serps", SERP.Type).Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}

// Indexes of the SearchQuery.
func (SearchQuery) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("query", "loc_code", "language", "job_id", "is_finished").
			Unique(),
	}
}
