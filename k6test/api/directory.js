import encoding from 'k6/encoding';
import http from 'k6/http';
import { uuidv4 } from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';
import { randomString } from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';
import chai, { describe, expect } from 'https://jslib.k6.io/k6chaijs/4.3.4.3/index.js';
import { ServerUrl, EHRSystemID } from "./consts.js";
import exp from 'constants';


// r.POST("/:ehrid/directory", a.Directory.Create)
export function create_directory(ctx, ehrID, userID) {
    const directory = {
        _type: "VERSIONED_FOLDER",
        name: {
            _type: "DV_TEXT",
            value: "TEST_DIRECTORY" + uuidv4(),
        },
        archetype_node_id: "openEHR-EHR-COMPOSITION.health_summary.v1",
        feeder_audit: {
            _type: "FEEDER_AUDIT",
            originating_system_audit: {
                _type: "PARTY_IDENTIFIED",
                name: "TEST",
            },
        },
        folders: [],
        items: [],
        details: {
            _type: "ITEM_SINGLE",
            name: {
                _type: "DV_TEXT",
                value: "TEST_DIRECTORY" + uuidv4(),
            },
            archetype_node_id: "openEHR-EHR-COMPOSITION.health_summary.v1",
            item: {
                _type: "VERSIONED_OBJECT_REF",
            }
        },
    }

    const payload = JSON.stringify(directory);

    const versiondID = uuidv4();

    const response = ctx.session.post(`/ehr/${ehrID}/directory/${versiondID}?patient_id=${userID}`, payload);

    expect(response.status, "create directory").to.equal(201); 
}

// r.PUT("/:ehrid/directory", a.Directory.Update)
export function update_directory(ctx, ehrID, directoryID, directory) {
    //TODO - add update directory functionality
} 

// r.DELETE("/:ehrid/directory", a.Directory.Delete)
export function delete_directory(ctx, ehrID, directoryID) {
    //TODO - add delete directory functionality
}

// r.GET("/:ehrid/directory", a.Directory.GetByTime)
export function get_directory_by_time(ctx, ehrID, path) {
    const response = ctx.session.get(`/ehr/${ehrID}/directory?path=${path}`);

    expect(response.status, "get directory by time").to.equal(200);
    expect(response).to.have.validJsonBody();

    const resp = JSON.parse(response.body);

    return resp;

}
// r.GET("/:ehrid/directory/:version_uid", a.Directory.GetByVersion)
export function get_directory_by_version(ctx, ehrID, directoryID) {
    const response = ctx.session.get(`/ehr/${ehrID}/directory/${directoryID}`);

    expect(response.status, "get directory by version").to.equal(200);
    expect(response).to.have.validJsonBody();

    const resp = JSON.parse(response.body);

    return resp;
}