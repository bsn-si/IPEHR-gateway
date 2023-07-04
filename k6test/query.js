import encoding from 'k6/encoding';
import { uuidv4 } from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';
import { randomString } from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';
import chai, { describe, expect } from 'https://jslib.k6.io/k6chaijs/4.3.4.3/index.js';


// r.GET("/query/:qualified_query_name", a.Query.ExecStoredQuery)
// r.POST("/query/:qualified_query_name", a.Query.PostExecStoredQuery)
// r.GET("/query/aql", a.Query.ExecGetQuery)
// r.POST("/query/aql", a.Query.ExecPostQuery)


// query.GET("/definition/query/:qualified_query_name", a.Query.ListStored)

// query.PUT("/definition/query/:qualified_query_name", a.Query.Store)
// query.PUT("/definition/query/:qualified_query_name/:version", a.Query.StoreVersion)
// query.GET("/definition/query/:qualified_query_name/:version", a.Query.GetStoredByVersion)