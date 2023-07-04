import encoding from 'k6/encoding';
import http from 'k6/http';
import { uuidv4 } from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';
import { randomString } from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';
import chai, { describe, expect } from 'https://jslib.k6.io/k6chaijs/4.3.4.3/index.js';
import { ServerUrl, EHRSystemID } from "./consts.js";
import exp from 'constants';


// r.POST("/:ehrid/composition", a.Composition.Create)
export function create_composition(ctx, ehrID) {
    const payload = JSON.stringify({
        _type: "COMPOSITION",
        name: "TEST_COMPOSITION" + uuidv4(),
        content: []
    });

    const response = ctx.session.post(`/ehr/${ehrID}/composition`, payload);
    expect(response.status, "create composition").to.equal(201);
    expect(response).to.have.validJsonBody();

    const resp = JSON.parse(response.body);

    return resp;
}

// r.GET("/:ehrid/composition", a.Composition.GetList)
export function get_composition_list(ctx, ehrID) {
    const response = ctx.session.get(`/ehr/${ehrID}/composition`);

    expect(response.status, "create composition").to.equal(200);
    expect(response).to.have.validJsonBody();

    const resp = JSON.parse(response.body);

    return resp;
}

// r.GET("/:ehrid/composition/:version_uid", a.Composition.GetByID)
export function get_composition_by_id(ctx, ehrID, compositionID) {
    const response = ctx.session.get(`/ehr/${ehrID}/composition/${compositionID}`);

    expect(response.status, "create composition").to.equal(202);
    expect(response).to.have.validJsonBody();

    const resp = JSON.parse(response.body);

    return resp;
}

// r.DELETE("/:ehrid/composition/:preceding_version_uid", a.Composition.Delete)
export function delete_composition(ctx, ehrID, compositionID) {
    const response = ctx.session.delete(`/ehr/${ehrID}/composition/${compositionID}`);

    expect(response.status, "delete composition").to.equal(204);
}

// r.PUT("/:ehrid/composition/:versioned_object_uid", a.Composition.Update)
export function update_composition(ctx, ehrID, compositionID, composition) {
    const payload = JSON.stringify(composition);

    const response = ctx.session.put(`/ehr/${ehrID}/composition/${compositionID}`, payload);

    expect(response.status, "update composition").to.equal(200);
    expect(response).to.have.validJsonBody();
}
