import encoding from 'k6/encoding';
import http from 'k6/http';
import { sleep } from 'k6';
import { uuidv4 } from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';
import { randomString } from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';
import chai, { describe, expect } from 'https://jslib.k6.io/k6chaijs/4.3.4.3/index.js';
import { Httpx, Get, Post } from 'https://jslib.k6.io/httpx/0.0.6/index.js';
import { ServerUrl, EHRSystemID } from "./consts.js";


const USER_ROLE = 0;
const DOCTOR_ROLE = 1;

// register new user
export function register_user(ctx) {
    return register_user_with_role(ctx, USER_ROLE);
}

// register new Doctor
export function register_doctor(ctx) {
    return register_user_with_role(ctx, DOCTOR_ROLE);
}

function register_user_with_role(ctx, role) {
    let user = {
        userID: uuidv4(),
        username: randomString(10, `aeioubcdfghijpqrstuv`),
        password: randomString(10, `aeioubcdfghijpqrstuv`),
    }

    ctx.session.addHeader('AuthYserId', user.userID);

    const payload = JSON.stringify({
        role: role,
        userID: user.userID,
        name: user.username,
        password: user.password,
        address: "some docktor adderess",
        description: "any description for doctor",
        pictureURL: "some url for doctor picture",
    });

    let response = ctx.session.post(`/user/register/`, payload);

    expect(response.status, "registration status").to.equal(201);
    expect(response).to.have.validJsonBody();
    
    console.log(response.body);

    return user;
}

// login user
export function login_user(ctx, userID, password) {
    const payload = JSON.stringify({
        userID: userID,
        password: password,
    });

    let response = ctx.session.post(`/user/login/`, payload);

    expect(response.status, "login status").to.equal(200);
    expect(response).to.have.validJsonBody()

    const access_token = response.json('access_token');
    const refresh_token = response.json('refresh_token');

    ctx.session = new Httpx({
        baseURL: ServerUrl,
        headers: {
            AuthUserId: userID,
            EhrSystemId: EHRSystemID,
            Authorization: `Bearer ${access_token}`,
        },
    });

    ctx.access_token = access_token;
    ctx.refresh_token = refresh_token;
}

export function refresh_token(ctx) {
    ctx.session = new Httpx({
        baseURL: ServerUrl,
        headers: {
            AuthUserId: userID,
            EhrSystemId: EHRSystemID,
            Authorization: `Bearer ${ctx.refresh_token}`,
        },
    });

    let response = ctx.session.get(`/user/refresh/`);

    expect(response.status, "refresh JWT token").to.equal(200);
    expect(response).to.have.validJsonBody()

    const access_token = response.json('access_token');
    const refresh_token = response.json('refresh_token');

    ctx.session = new Httpx({
        baseURL: ServerUrl,
        headers: {
            AuthUserId: userID,
            EhrSystemId: EHRSystemID,
            Authorization: `Bearer ${access_token}`,
        },
    });

    ctx.access_token = access_token;
    ctx.refresh_token = refresh_token;
}

export function get_user_info(ctx, user) {
    let response = ctx.session.get('/user/' + user.userID);

    expect(response.status, "Get User Info").to.equal(200);
    expect(response).to.have.validJsonBody()

    expect(response.json('userID'), "Correct User ID").to.equal(user.userID);
}

export function get_doctor_info(ctx, user, code) {
    const response = ctx.session.get('/user/code/' + code);

    expect(response.status, "Get Doctor By Code").to.equal(200);
}

export function log_out(ctx) {
    let response = ctx.session.get('/user/logout');

    expect(response.status, "User Logout").to.equal(200);
    expect(response).to.have.validJsonBody()
}