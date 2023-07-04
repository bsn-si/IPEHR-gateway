import encoding from 'k6/encoding';
import { uuidv4 } from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';
import { randomString } from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';
import chai, { describe, expect } from 'https://jslib.k6.io/k6chaijs/4.3.4.3/index.js';
import exp from 'constants';


// query.GET("/definition/query/:qualified_query_name", a.Query.ListStored)
export function get_query_list(ctx, qualified_query_name) {
    const response = ctx.session.get(`/definition/query/${qualified_query_name}`);

    expect(response.status, "get query list").to.equal(200);
    expect(response).to.have.validJsonBody();

    const resp = JSON.parse(response.body);

    return resp;
}

// query.GET("/definition/query/:qualified_query_name/:version", a.Query.GetStoredByVersion)
export function store_query(ctx, qualified_query_name, version) {
    const response = ctx.session.get(`/definition/query/${qualified_query_name}/${version}`);

    expect(response.status, "store query").to.equal(200);
    expect(response).to.have.validJsonBody();

    const resp = JSON.parse(response.body);

    return resp;
}

// query.PUT("/definition/query/:qualified_query_name", a.Query.Store)
export function store_query(ctx) {
    const qualified_query_name = "test_query" + uuidv4();

    const aqlQueryString = "SELECT c FROM EHR[ehr_id/value='123'] CONTAINS COMPOSITION c[openEHR-EHR-COMPOSITION.health_summary.v1]";
    
    const response = ctx.session.put(`/definition/query/${qualified_query_name}`, aqlQueryString);

    expect(response.status, "store query").to.equal(201);
    return qualified_query_name;
}

// query.PUT("/definition/query/:qualified_query_name/:version", a.Query.StoreVersion)
export function store_query_version(ctx) {
    const qualified_query_name = "test_query" + uuidv4();
    const version = "1.0.0";

    const aqlQueryString = "SELECT c FROM EHR[ehr_id/value='123'] CONTAINS COMPOSITION c[openEHR-EHR-COMPOSITION.health_summary.v1]";

    const response = ctx.session.put(`/definition/query/${qualified_query_name}/${version}`, aqlQueryString);

    expect(response.status, "store query version").to.equal(201);

    // check response header Location
    const location = response.headers["Location"];
    expect(location).to.contain(`/definition/query/${qualified_query_name}/${version}`);
    
    return {name: qualified_query_name, version: version};
}