## Directory of users and groups

### Roles

At this stage there are two roles used in the system: `Patient` and `Doctor`, whose capabilities will differ.

### Users

In IPEHR the `user` is the following structure:

```
struct User {
	userID   [32]byte
	systemID [32]byte
	role     Role
	groups   [][32]byte
	pwdHash  []byte
}
```

### User groups:

```
struct UserGroup {
	groupID     [32]byte
	owner       [32]byte
	description string
	members     map[string]Access
}
```

Only a member with `Owner` or `Admin` access rights can add users to the group.


Users and user groups are stored in [IPEHR smart contract](https://github.com/bsn-si/IPEHR-blockchain-indexes)

```
mapping (address => User)      users;
mapping (bytes32 => UserGroup) userGroups;
```

## Registration

Pre-registration is required to work with the IPEHR system. When registering, the following parameters must be specified:

- userID - user ID
- systemID - HMS identifier
- role - user role ("Patient" or " Doctor")
- password

User account information is written to a smart contract.  
Before saving, the password is hashed using a special function [scrypt](https://en.wikipedia.org/wiki/Scrypt)  
The API specification of the registration method can be found [here]()

## Authentication

The authorization of requests to the IPEHR gateway API is done via the JWT token.

To get a JWT token, the user performs an authentication procedure using the API method described [here]().

The validity period of the tokens is set in the IPEHR gateway configuration file.

## Request authorization

All requests must contain the following header:

```
Authorization: Bearer <JWT>
```

If the IPEHR gateway successfully authorizes and validates the `userID`, it finds the user key in the `keystore` and executes the request.

For authorization and access rights verification when communicating with an IPEHR smart contract, an ECDSA signature is added to the request with the key of the user who initiated the request.