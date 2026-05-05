# Caveats/workaround to consider

## go-graphql-client

Due to https://github.com/hasura/go-graphql-client/issues/158, queries that
contain nested slices should not reuse the same query object across calls
(e.g. for iteration over pagination).  They should rather define the query var
inside the loop.  (also see related
https://github.com/hasura/go-graphql-client/issues/152)
