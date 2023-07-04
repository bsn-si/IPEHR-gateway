import encoding from 'k6/encoding';
import http from 'k6/http';
import { uuidv4 } from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';
import { randomString } from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';
import chai, { describe, expect } from 'https://jslib.k6.io/k6chaijs/4.3.4.3/index.js';
import { ServerUrl, EHRSystemID } from "./consts.js";
import exp from 'constants';

//  /ehr/{ehr_id}/contribution/{contribution_uid} [get]
export function get_contribution_by_id(ctx, ehr_id, contribution_id) {
    const response = ctx.session.get(`/ehr/${ehr_id}/contribution/${contribution_id}`);

    expect(response.status, "get contribution by id").to.equal(200);
    expect(response).to.have.validJsonBody();

    const resp = JSON.parse(response.body);

    return resp;
}

// /ehr/{ehr_id}/contribution [post]
export function create_contribution(ctx, ehr_id) {
    const contribution = {
        _type: "CONTRIBUTION",
        name: {
            _type: "DV_TEXT",
            value: "TEST_CONTRIBUTION" + uuidv4(),
        },
        archetype_node_id: "openEHR-EHR-COMPOSITION.health_summary.v1",
        uid: {
            _type: "UID_BASED_ID",
            value: uuidv4(),
        },
        versions: [],
        audit: {
            _type: "AUDIT_DETAILS",
            system_id: EHRSystemID,
        },
    }

    const payload = JSON.stringify(contribution);

    const response = ctx.session.post(`/ehr/${ehr_id}/contribution`, payload);

    expect(response.status, "create contribution").to.equal(201);

    return contribution;
}
