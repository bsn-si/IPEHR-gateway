export class IperhApi {
    public apiUrl: string
    public username: string
    public password: string
    public betterToken: string
    public systemId: string

    public constructor(init?:Partial<IperhApi>) {
        Object.assign(this, init);
    }

    public uuidRegexp = /[0-9a-fA-F]{8}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{12}/gm;

    public async login() {
        const res = await fetch(this.apiUrl + '/user/login', {
          method: 'POST',
          headers: {
            AuthUserId: this.username,
            EhrSystemId: this.systemId
          },
          body: JSON.stringify({
            password: this.password,
            userID: this.username
          })
        })

        if (res.status !== 200) {
            console.log(new Error('IPEHR Login error:', res.json()));
            return ''
        }
        const json = await res.json()
        return json.access_token || ''
    }

    public async createEhr(url: string, ehrId: string) {
        if (!url) {
            console.log(new Error('createEhr: url is empty'));
            return;
        }

        if (!ehrId) {
            console.log(new Error('createEhr: ehrId is empty'));
            return;
        }

        const searchParams = new URLSearchParams(url.split("?")[1]);
        const subjectId = searchParams.get('subjectId');
        const subjectNamespace = searchParams.get('subjectNamespace');

        if (subjectId && subjectNamespace) {
          try {
            const token: any = await this.login();
            if (token) {
              console.log('LOGIN RESULT', token);
              const res = await fetch(`${this.apiUrl}/ehr/${ehrId}`, {
                method: 'PUT',
                headers: {
                  Authorization: `Bearer ${token}`,
                  'AuthUserId': this.username,
                  'EhrSystemId': this.systemId,
                  Prefer: 'return=representation',
                },
                body: JSON.stringify({
                  "_type": "EHR_STATUS",
                  "archetype_node_id": "openEHR-EHR-EHR_STATUS.generic.v1",
                  "isModifiable": true,
                  "isQueryable": true,
                  "name": {
                    "value": "EHR Status"
                  },
                  "subject": {
                    "external_ref": {
                      "id": {
                        "_type": "GENERIC_ID",
                        "scheme": "id_scheme",
                        "value": subjectId
                      },
                      "namespace": subjectNamespace,
                      "type": "PERSON"
                    }
                  }
                })
              })

              if (res.status !== 201) {
                  console.log(new Error('ERROR CREATING EHR:', res.json()));
                  return ''
              }
            }
          } catch (e) {
            console.log('ERROR CREATING EHR', e)
          }
        }
    }

    public async createComposition(ehrId: string, compositionUid: string) {
        console.log('CREATING COMPOSITION');

        if (!ehrId) {
            console.log(new Error('createComposition: ehrId is empty'));
            return;
        }

        if (!compositionUid) {
            console.log(new Error('createComposition: compositionUid is empty'));
            return;
        }

        try {
            const getCompositionRes = await fetch(`https://sandbox.better.care/ehr/rest/openehr/v1/ehr/${ehrId}/composition/${compositionUid}`, {
              headers: {
                'Authorization': `Bearer ${this.betterToken}`,
              }
            })

            if (getCompositionRes.status !== 200) {
                console.log(new Error('FETCH COMPOSITION ERROR:', getCompositionRes.json()));
                return;
            }

            const compositionBody = await getCompositionRes.text();
            console.log('FETCH COMPOSITION FROM BETTER API', compositionBody)

            const token: any = await this.login()
            if (token) {
                console.log('LOGIN RESULT', token)
                const res = await fetch(`${this.apiUrl}/ehr/${ehrId}/composition`, {
                    method: 'POST',
                    headers: {
                      Authorization: `Bearer ${token}`,
                      'AuthUserId': this.username,
                      'EhrSystemId': this.systemId,
                      Prefer: 'return=representation',
                    },
                    body: compositionBody
                });

                if (res.status !== 201) {
                    console.log(new Error('ERROR CREATING EHR:', res.json()));
                    return;
                }
            }
        } catch (e) {
            console.log(e)
        }
    }

}
