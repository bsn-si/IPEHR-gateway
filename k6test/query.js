import encoding from 'k6/encoding';
import { uuidv4 } from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';
import { randomString } from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';
import chai, { describe, expect } from 'https://jslib.k6.io/k6chaijs/4.3.4.3/index.js';


// r.POST("/query/aql", a.Query.ExecPostQuery)
export function post_exec_query(ctx, aqlQueryString) {
    const payload = JSON.stringify({});
    const response = ctx.session.post(`/query/aql`, aqlQueryString, payload);

    expect(response.status, "exec query").to.equal(200);
    expect(response).to.have.validJsonBody();

    const resp = JSON.parse(response.body);

    return resp;
}

// r.POST("/query/:qualified_query_name", a.Query.PostExecStoredQuery)
export function post_exec_stored_query(ctx, qualified_query_name) {
    const payload = JSON.stringify({});
    const response = ctx.session.post(`/query/${qualified_query_name}`, payload);

    expect(response.status, "exec stored query").to.equal(200);
    expect(response).to.have.validJsonBody();

    const resp = JSON.parse(response.body);

    return resp;
}

// r.GET("/query/aql", a.Query.ExecGetQuery)
export function get_exec_query(ctx, aqlQueryString) {
    const response = ctx.session.get(`/query/aql/?$q=${aqlQueryString}`);

    expect(response.status, "exec query").to.equal(200);
    expect(response).to.have.validJsonBody();

    const resp = JSON.parse(response.body);

    return resp;
}

// r.GET("/query/:qualified_query_name", a.Query.ExecStoredQuery)
export function get_exec_stored_query(ctx, qualified_query_name) {
    const response = ctx.session.get(`/query/${qualified_query_name}`);

    expect(response.status, "exec stored query").to.equal(200);
    expect(response).to.have.validJsonBody();

    const resp = JSON.parse(response.body);

    return resp;
}
