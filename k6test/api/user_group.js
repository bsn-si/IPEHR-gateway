import encoding from 'k6/encoding';
import http from 'k6/http';
import { sleep } from 'k6';
import { uuidv4 } from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';
import { randomString } from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';
import chai, { describe, expect } from 'https://jslib.k6.io/k6chaijs/4.3.4.3/index.js';
import { Httpx, Get, Post } from 'https://jslib.k6.io/httpx/0.0.6/index.js';
import { ServerUrl, EHRSystemID } from "./consts.js";
import exp from 'constants';


export function create_user_group(ctx, users) {
    const user_group = {
        groupID: uuidv4(),
        name: randomString(10, `aeioubcdfghijpqrstuv`),
        description: "some description",
        members: users,
    }

    const payload = JSON.stringify(user_group);

    const response = ctx.session.post(`/user/group/`, payload);

    expect(response.status, "Create User Group").to.equal(201);
    expect(response).to.have.validJsonBody();

    const group = JSON.parse(response.body);

    return group;
}

export function get_user_group_by_id(ctx, group_id) {
    const response = ctx.session.get('/user/group/' + group_id);

    expect(response.status, "Get User Group By ID").to.equal(200);
    expect(response.body).to.have.validJsonBody();

    const group = JSON.parse(response.body);

    expect(group.groupID, "Check Group ID").to.equal(group_id);

    return group;
}

export function add_user_to_group(ctx, group_id, user_id, access_level) {
    const url = "/user/group/" + group_id + "/user_add/" + user_id + "/" + access_level
    const response = ctx.put(url)

    expect(response.status, "Get User Group By ID").to.equal(200);
    expect(response.body).to.have.validJsonBody();
}

export function remove_user_from_group(ctx, group_id, user_id ){ 
    const url = "/user/group/" + group_id + "/remove_add/" + user_id;
    
    const response = ctx.post(url);

    expect(response.status, "Get User Group By ID").to.equal(200);
    expect(response.body).to.have.validJsonBody();
}

export function get_list_of_user_groups(ctx) {
    const response = ctx.session.get('/user/group/');

    expect(response.status, "Get User Groups Lust").to.equal(200);
    expect(response.body).to.have.validJsonBody();

    const groups = JSON.parse(response.body);

    return groups;
}