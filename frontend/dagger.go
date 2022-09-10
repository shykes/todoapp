package main

//go:generate dagger client gofile -o ./generated.go
import (
    "context"

    "dagger.io/dagger"
)

// NOTES:
//
// 1. Dagger generates a custom query builder called `Graph`
// 2. Presumably it would also generate input types for all fields in the graph, into a flat namespace.
//       This may not be practical, for now we are avoiding the use of input types
// 3. To keep things simple we assume that struct resolvers are "all or nothing": clients cannot differentiate
//      between an error resolving a struct, and an error resolving one of its fields.
//       Let's see how far we can go with this constraint.
// 4. Scalar fields can be included / excluded to avoid unnecessary network I/O
// 5. The API only supports a subset of queries, where scalars are only queried from one struct;
//     that struct is returned to the caller (intermediary fields are removed). We call this "flattening"
func init() {
    // - Graph is the (generated) client for the base graph we want to extend
    // - TodoApp is our graph extensions
    dagger.Extend(Graph, TodoApp))
}

func TodoApp(ctx context.Context, g Graph, source dagger.DirectoryID) *App {
    return &App{
        Source: source,
    }
}

type App {
    Source dagger.DirectoryID
}

func (app *App) Build(ctx context.Context, g Graph) (*Build, error) {
    q := g.Yarn().Run("build", app.source)
    resp, err := q.Execute(ctx, dagger.IncludeFields("Contents", "Logs"))
    if err != nil {
        return nil, err
    }
    return &Build{
        Contents: resp.Contents,
        Logs: resp.Logs,
    }, nil
}

type Build {
    Contents dagger.DirectoryID
    Logs string
}

func (b *Build) Deployment(ctx context.Context, g Graph, siteName string, token dagger.SecretID) (*Deployment, error) {
    q := g.Netlify(token).Site(siteName)).Deployment(b.contents)
    res, err := q.Execute(ctx, dagger.IncludeField("URL", "Logs"))
    if err != nil {
        return nil, err
    }
    return &Deployment{
        URL: res.URL,
        Logs: res.Logs,
    }, nil
}

type Deployment {
    URL string
    Logs string
}