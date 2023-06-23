import encoding from 'k6/encoding';
import http from 'k6/http';
import { sleep } from 'k6';
import { uuidv4 } from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';
import { randomString } from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';
import chai, { describe, expect } from 'https://jslib.k6.io/k6chaijs/4.3.4.3/index.js';
import { Httpx, Get, Post } from 'https://jslib.k6.io/httpx/0.0.6/index.js';
import { ServerUrl, EHRSystemID } from "./consts.js";


export function create_user_group(ctx, users) {
    let user_group = {
        groupID: uuidv4(),
        name: randomString(10, `aeioubcdfghijpqrstuv`),
        description: "some description",
        members: users,
    }

    const payload = JSON.stringify(user_group);

    let response = ctx.session.post(`/user/group/`, payload);

    expect(response.status, "Create User Group").to.equal(201);
    expect(response).to.have.validJsonBody();

    let group = JSON.parse(response.body);

    return group;
}

export function get_user_group(ctx, groupID) {
}