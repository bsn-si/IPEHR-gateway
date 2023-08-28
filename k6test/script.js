import encoding from 'k6/encoding';
import http from 'k6/http';
import { sleep, fail } from 'k6';
import { Rate } from 'k6/metrics';
import { uuidv4 } from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';
import { randomString } from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';
import chai, { describe, expect } from 'https://jslib.k6.io/k6chaijs/4.3.4.3/index.js';
import { Httpx, Get, Post } from 'https://jslib.k6.io/httpx/0.0.6/index.js';
import { ServerUrl, EHRSystemID } from "./api/consts.js";
import * as user from './api/user.js' // import all exported functions from user.js;
import * as ehr from './api/ehr.js' // import all exported functions from ehr.js;

let session = new Httpx({
    baseURL: ServerUrl,
    headers: {
        EhrSystemId: EHRSystemID,
    },
});

chai.config.logFailures = true;

const formFailRate = new Rate('failed form fetches');
const submitFailRate = new Rate('failed form submits');

export const options = {
    iterations: 100,
    vus: 20,
    // vus: 50,
    // vus: 20,
    duration: '3600s',
    ext: {
        loadimpact: {
            // Project: Default project
            projectID: 3644252,
            // Test runs with the same name groups test runs together
            name: 'YOUR TEST NAME'
        }
    },
    thresholds: {
        'failed form submits': ['rate<0.1'],
        'failed form fetches': ['rate<0.1'],
        // 'http_req_duration': ['p(95)<400'],
        // fail the test if any checks fail or any requests fail
        'checks': ['rate == 1.00'],
        'http_req_failed': ['rate == 0.00'],
    },
};

export default function testSuite() {
    let ctx = {
        session: session,
    };

    let u = null;

    let ok = describe('Register user', () => {
        u = user.register_user(ctx);
        console.log("USER:" + JSON.stringify(u));
    });

    // wait up to 15 minutes with 10 seconds interval
    if (!ok || !retries(ctx, 90, 10, ctx.RequestId)) {
        fail('failed to register user');
        return;
    }

    sleep(60);

    ok = describe('Login user', () => {
        user.login_user(ctx, u.userID, u.password);
    });

    if (!ok) {
        fail('failed to login user');
        return;
    }

    let ehrDoc = null;

    ok = describe('Create EHR', () => {
        ehrDoc = ehr.create_ehr(ctx, u);
        console.log("EHR:" + JSON.stringify(ehrDoc));
    });

    // wait up to 15 minutes with 10 seconds interval
    if (!ok || !retries(ctx, 90, 10, ctx.RequestId)) {
        fail('failed to create ehr');
        return;
    }

    sleep(60);

    ok = describe('Login user after EHR created', () => {
        ctx.session.clearHeader('Authorization');
        user.login_user(ctx, u.userID, u.password);
    });

    if (!ok) {
        fail('failed to login user');
        return;
    }

    ok = describe('Get EHR', () => {
        const doc = ehr.get_ehr(ctx, ehrDoc.ehr_id.value);
        console.log("EHR:" + JSON.stringify(doc));

        expect(JSON.stringify(ehrDoc), 'try to compare ehrs').to.equal(JSON.stringify(doc));
    });

    if (!ok) {
        fail('failed to get ehr');
        return;
    }

    ok = describe('Get user', () => {
        const userInfo = user.get_user_info(ctx, u);

        console.log("UserInfo:" + JSON.stringify(userInfo));
        expect(ehrDoc.ehr_id.value, 'try to compare ehrs').to.equal(userInfo.ehrID);
    });

    if (!ok) {
        fail('failed to get user info');
        return;
    }

    if (!describe('Logout user', () => {
        user.log_out(ctx);
    })) {
        fail('failed to logout user');
        return;
    }
}

function retries(ctx, retriesCount, interval, requestId) {
    while (retriesCount > 0) {
        const response = ctx.session.get(`/requests/${requestId}`);

        if (response.status !== 200) {
            return false;
        }

        const body = JSON.parse(response.body);

        switch (body.status) {
            case 'Success':
                return true;
            case 'Failed':
                return false;
            case 'Unknown':
                return false;
            case 'Pending', 'Processing':
                let isEtherFinished = true;

                const ethereum = body.ethereum;
                for (let i = 0; i < ethereum.length; i++) {
                    if (ethereum[i].Status !== 'Success') {
                        isEtherFinished = false;
                        break;
                    }
                }

                if (isEtherFinished) {
                    return true;
                }
        }

        retriesCount--;
        sleep(interval);
        continue;
    }

    return false;
}
