# Cryptography

![Deliverable 3](https://user-images.githubusercontent.com/98888366/170701371-64966c5b-05ae-465a-8ebd-50a06160d98a.png)

### Encryption
When stored in the repository, EHR documents are pre-encrypted using the ChaCha20-Poly1305 streaming algorithm with message authentication. The protocol is standardized by IETF in RFC 7539, in software implementations it is much more efficient and faster than AES. To encrypt each document a unique key is generated - a random sequence of 256 bits (32 bytes) + a unique 96 bits (12 bytes). Document ID is used as an authentication tag.

Public key cryptography (asymmetric encryption) is used to encrypt the symmetric key. A public and private key pair is generated for each user.
A set of algorithms Curve25519, XSalsa20, Poly1305 is used for encryption.
Curve25519 is an elliptic curve and a set of parameters for it, chosen so as to provide better performance (on average, 20-25%) and get rid of some security problems with traditional ECDH. Described in RFC 7748 in 2016.

There are many implementations of the presented algorithms in libraries for different programming languages.  

**golang**:  
<http://golang.org/x/crypto/chacha20poly1305>  
<http://golang.org/x/crypto/sha3>  
<http://golang.org/x/crypto/nacl>

**js**:  
<https://github.com/emn178/js-sha3>  
<https://tweetnacl.js.org>

### Hashing

To calculate hash sums, SHA3-256 is used.
SHA3(Keccak) is a cryptographic hash function that won the US National Institute of Standards and Technology (NIST) competition in 2007 - 2012. On August 5, 2015, the algorithm was approved and published as the FIPS 202 standard.

### Authentication

A digital signature is added to each document that is saved to authenticate the identity of the user/official who is saving the document.
Legally significant certificates can be used to sign documents. For the EU, these are certificates issued by licensed certification centers according to the unified eIDAS standard.  
<https://ec.europa.eu/digital-building-blocks/wikis/display/DIGITAL/eSignature>  
For technical implementation of the signing mechanism software supporting DSS library can be used.  
<https://ec.europa.eu/digital-building-blocks/DSS/webapp-demo/doc/dss-documentation.html>

## Implementation

The packages that implement the encryption functionality are located in the `pkg/crypto` project directory.

Running tests to demonstrate the work of the packages:

```
go test -v ./pkg/crypto/...
```
