import encoding from 'k6/encoding';
import http from 'k6/http';
import { uuidv4 } from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';
import { randomString } from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';
import chai, { describe, expect } from 'https://jslib.k6.io/k6chaijs/4.3.4.3/index.js';
import { ServerUrl, EHRSystemID } from "./consts.js";


export function update_ehr_status(ctx, ehrID, status) {
    const payload = JSON.stringify({
        _type: "EHR_STATUS",
        is_queryable: true,
        is_modifiable: true,
    });

    const response = ctx.session.put(`/ehr/${ehrID}/status`, payload);

    expect(response.status, "update ehr status").to.equal(200);
    expect(response).to.have.validJsonBody();
}

export function get_ehr_status_by_time(ctx, ehrID, time) {
    const response = ctx.session.get(`/ehr/${ehrID}/ehr_status?version_at_time=${time}`);

    expect(response.status, "get ehr status by time").to.equal(200);
    expect(response).to.have.validJsonBody();
}

export function get_ehr_status_by_version(ctx, ehrID, version) {
    const response = ctx.session.get(`/ehr/${ehrID}/ehr_status/${version}`);

    expect(response.status, "get ehr status by version").to.equal(200);
    expect(response).to.have.validJsonBody();
}