import encoding from 'k6/encoding';
import http from 'k6/http';
import { sleep } from 'k6';
import { uuidv4 } from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';
import { randomString } from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';
import chai, { describe, expect } from 'https://jslib.k6.io/k6chaijs/4.3.4.3/index.js';
import { Httpx, Get, Post } from 'https://jslib.k6.io/httpx/0.0.6/index.js';
import { ServerUrl, EHRSystemID } from "./consts.js";
import { register_user, login_user } from './user.js';

let session = new Httpx({
    baseURL: ServerUrl,
    headers: {
        EhrSystemId: EHRSystemID,
    },
});

chai.config.logFailures = true;

export const options = {
    vus: 1,
    iterations: 5,
    // vus: 10,
    // duration: '10s',
    ext: {
        loadimpact: {
            // Project: Default project
            projectID: 3644252,
            // Test runs with the same name groups test runs together
            name: 'YOUR TEST NAME'
        }
    },
    thresholds: {
        // fail the test if any checks fail or any requests fail
        checks: ['rate == 1.00'],
        http_req_failed: ['rate == 0.00'],
    },
};

export default function testSuite() {
    let ctx = {
        session: session,
    };

    let user = null;

    describe('Register user', () => {
        user = register_user(ctx);
    });

    describe('Login user', () => {
        login_user(ctx, user.userID, user.password);
    });
}
