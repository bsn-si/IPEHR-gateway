import encoding from 'k6/encoding';
import http from 'k6/http';
import { sleep } from 'k6';
import { uuidv4 } from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';
import { randomString } from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';
import chai, { describe, expect } from 'https://jslib.k6.io/k6chaijs/4.3.4.3/index.js';
import { Httpx, Get, Post } from 'https://jslib.k6.io/httpx/0.0.6/index.js';
import { ServerUrl, EHRSystemID } from "./consts.js";

export function create_ehr(ctx, user) {
    const ehr = {
        _type: "EHR",
        archetype_node_id: "openEHR-EHR-EHR_EXAMPLE.v1",
        name: {
            value: "ehr_name"
        },
        subject: {
            external_ref: {
                id: {
                    _type: "GENERIC_ID",
                    value: user.userID,
                    scheme: "ehrId",
                },
                namespace: "ehr",
                type: "PERSON",
            },
        },
        isModifiable: true,
        isQueryable: true,
    }

    const payload = JSON.stringify(ehr);

    const session = new Httpx({
        baseURL: ServerUrl,
        headers: {
            AuthUserId: user.userID,
            EhrSystemId: EHRSystemID,
            Authorization: `Bearer ${ctx.access_token}`,
            Prefer: "return=representation",
        },
    });

    const response = session.post(`/ehr/`, payload);

    expect(response.status, "create ehr status").to.equal(201);
    expect(response).to.have.validJsonBody();

    const resp = JSON.parse(response.body);

    expect(resp.system_id.value).to.be.a('string');
    expect(resp.ehr_id.value).to.be.a('string');

    ctx.RequestId = response.headers['Requestid'];

    return resp;
}

export function create_ehr_with_id(ctx, user, ehrID) {
    const ehr = {
        _type: "EHR",
        archetype_node_id: "openEHR-EHR-EHR_EXAMPLE.v1",
        name: {
            value: "ehr_name"
        },
        subject: {
            external_ref: {
                id: {
                    _type: "GENERIC_ID",
                    value: user.userID,
                    scheme: "ehrId",
                },
                namespace: "ehr",
                type: "PERSON",
            },
        },
        isModifiable: true,
        isQueryable: true,
        ehr_id: uuidv4(),
    }

    const payload = JSON.stringify(ehr);

    const response = ctx.session.put(`/ehr/` + ehrID, payload);

    expect(response.status, "create ehr status").to.equal(201);
    // expect(response).to.have.validJsonBody();

    return ehr;
}

export function get_ehr(ctx, ehrID) {
    const response = ctx.session.get(`/ehr/` + ehrID);

    expect(response.status, "get ehr status").to.equal(200);
    expect(response).to.have.validJsonBody();

    const resp = JSON.parse(response.body);

    expect(resp.system_id.value, "Check 'system_id' is string").to.be.a('string');
    expect(resp.ehr_id.value, "Check 'ehr_id' is string").to.be.a('string');

    return resp;
}

export function get_ehr_by_subject_id_and_namespace(ctx, subjectID, namespace) {
    const response = ctx.session.get(`/ehr?subjectId=` + subjectID + `&subjectNamespace=` + namespace);

    expect(response.status, "get ehr status").to.equal(200);
    expect(response).to.have.validJsonBody();

    const resp = JSON.parse(response.body);
    console.log("EHR:" + JSON.stringify(resp));

    expect(resp.system_id).to.be.a('string');
    expect(resp.ehr_id).to.be.a('string');

    return resp;
}