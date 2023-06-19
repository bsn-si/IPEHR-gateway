import encoding from 'k6/encoding';
import http from 'k6/http';
import { sleep } from 'k6';
import { uuidv4 } from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';
import { randomString } from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';
import chai, { describe, expect } from 'https://jslib.k6.io/k6chaijs/4.3.4.3/index.js';
import { Httpx, Get, Post } from 'https://jslib.k6.io/httpx/0.0.6/index.js';
import {ServerUrl, EHRSystemID} from "./consts.js";



// register new user
export function register_user(ctx) {
    let user = {
        userID: uuidv4(),
        username: randomString(10, `aeioubcdfghijpqrstuv`),
        password: randomString(10, `aeioubcdfghijpqrstuv`),
    }
    //a const username = randomString(10, `aeioubcdfghijpqrstuv`);
    // const password = randomString(10, `aeioubcdfghijpqrstuv`);
    // const userID = uuidv4();

    ctx.session.addHeader('AuthYserId', user.userID);

    const payload = JSON.stringify({
        address: "enim eu fugiat sunt",
        description: "dolore in eu Lorem dolor",
        name: user.username,
        pictureURL: "des",
        password: user.password,
        role: 0,
        userID: user.userID,
    });

    let response = ctx.session.post(`/user/register/`, payload);

    expect(response.status, "registration status").to.equal(201);
    expect(response).to.have.validJsonBody();

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

    ctx.session.addHeader('Authorization', `Bearer ${access_token}`);
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