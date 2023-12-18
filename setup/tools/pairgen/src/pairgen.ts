import express from "express";
import * as bodyParser from "body-parser";

import jose from "node-jose";
import jwt from "jsonwebtoken";
import * as forge from "node-forge";

import * as http from "http";
import * as https from "https";
import * as fs from "fs";
import * as path from "path";

const appHttps = express();
const appHttp = express();

const HTTP_PORT = parseInt(process.env.HTTP_PORT ?? "3000");
const HTTPS_PORT = parseInt(process.env.HTTPS_PORT ?? "3001");

const CA_CERT_FILES = process.env.CA_CERT_FILES;
const KEYS_FOLDERS = process.env.KEYS_FOLDERS;
const PAIRGEN_DNS_NAMES = process.env.PAIRGEN_DNS_NAMES ?? "pairgen";
const INTERMEDIATE_FILES = process.env.INTERMEDIATE_FILES ?? "inter";

// prettier-ignore
const attrsCA = [
    { name:      "commonName",       value: "MEx CA"      },
    { name:      "countryName",      value: "DE"          },
    { shortName: "ST",               value: "Brandenburg" },
    { name:      "localityName",     value: "Potsdam"     },
    { name:      "organizationName", value: "MEx CA"      },
    { shortName: "OU",               value: "MEx CA"      }
];

// prettier-ignore
const attrsCert = [
    { name:      "countryName",      value: "DE"          },
    { shortName: "ST",               value: "Brandenburg" },
    { name:      "localityName",     value: "Potsdam"     },
    { name:      "organizationName", value: "Pairgen"     },
    { shortName: "OU",               value: "Pairgen"     }
];

(async function () {
    process.stdout.write(`STDOUT: pairgen ${process.env.npm_package_version}\n`);
    process.stderr.write(`STDERR: pairgen ${process.env.npm_package_version}\n`);

    const keystore = jose.JWK.createKeyStore();

    const KEY_IDS = ((process.env.KEY_IDS ?? "") + " default").split(" ");
    if (KEY_IDS.length === 0) {
        throw new Error("no KEY_IDS defined");
    }

    const certCA = generateCACertificate(attrsCA);
    if (typeof CA_CERT_FILES !== "undefined") {
        for (let f of CA_CERT_FILES.split(",")) {
            fs.writeFileSync(f + ".pem", forge.pki.certificateToPem(certCA));
            console.log("wrote CA cert to", f + ".pem");
            fs.writeFileSync(f + ".key.pem", forge.pki.privateKeyToPem(certCA.privateKey));
            console.log("wrote CA cert key to", f + ".key.pem");
        }
    }

    const cert = generateIntermediateCertificate(PAIRGEN_DNS_NAMES.split(" "), attrsCert, certCA);
    for (let f of INTERMEDIATE_FILES.split(",")) {
        fs.writeFileSync(f + ".pem", forge.pki.certificateToPem(cert));
        console.log("wrote intermediate cert to", f + ".pem");
        fs.writeFileSync(f + ".key.pem", forge.pki.privateKeyToPem(cert.privateKey));
        console.log("wrote intermediate cert  key to", f + ".key.pem");
    }

    console.log("key IDs:", KEY_IDS);
    for (const keyId of KEY_IDS) {
        if (keyId.length === 0) {
            continue;
        }
        console.log("-", keyId);
        const key = await keystore.generate("RSA", 2048, {
            kid: keyId,
            alg: process.env.KEY_ALG ?? "RS256",
            use: process.env.KEY_USE ?? "sig",
        });

        if (typeof KEYS_FOLDERS !== "undefined") {
            for (let f of KEYS_FOLDERS.split(",")) {
                let keyFile = path.resolve(f, key.kid + ".key.pem");
                fs.writeFileSync(keyFile, key.toPEM(true));
                console.log("    - written:", keyFile);

                keyFile = path.resolve(f, key.kid + ".pubkey.pem");
                fs.writeFileSync(keyFile, key.toPEM(false));
                console.log("    - written:", keyFile);
            }
        }

        console.log("    - generated, kid:", key.kid);
    }

    appHttp.use((req, res, next) => {
        console.log(req.ip, req.url);
        next();
    });

    appHttps.use((req, res, next) => {
        console.log(req.ip, req.url);
        next();
    });

    appHttps.get("/public_keys", async (req, res) => {
        const keys = keystore.all();
        res.json({
            keys: keys.map((key) => (key as unknown as jose.JWK.Key).toJSON()),
        });
    });

    appHttp.get("/", async (req, res) => {
        res.send("ready");
    });

    appHttp.get("/private_keys/:kid", async (req, res) => {
        const keys = keystore.all({ kid: req.params.kid });
        if (keys.length === 0) {
            res.sendStatus(404);
        } else if (keys.length > 1) {
            res.sendStatus(409);
        } else {
            res.set("content-type", "text/plain");
            res.send((keys[0] as unknown as jose.JWK.Key).toPEM(true));
        }
    });

    appHttp.get("/public_keys/:kid", async (req, res) => {
        const format = req.query.format ?? "pem";
        const keys = keystore.all({ kid: req.params.kid });
        if (keys.length === 0) {
            res.sendStatus(404);
        } else if (keys.length > 1) {
            res.sendStatus(409);
        } else {
            res.set("content-type", "text/plain");
            const pem = (keys[0] as unknown as jose.JWK.Key).toPEM(false);
            switch (format) {
                case "pem":
                    res.send(pem);
                    break;
                case "sshd":
                    res.send(`ssh-rsa ${pem2sshd(pem)} ${req.params.kid}\n`);
                    break;
                default:
                    res.sendStatus(400);
            }
        }
    });

    appHttp.post("/jwts/:kid", bodyParser.json(), async (req, res) => {
        console.log("minting JWT with key:", req.params.kid);
        console.log(req.body);

        const keys = keystore.all({ kid: req.params.kid });
        if (keys.length === 0) {
            res.sendStatus(404);
        } else if (keys.length > 1) {
            res.sendStatus(409);
        } else {
            try {
                const token = jwt.sign(req.body, (keys[0] as unknown as jose.JWK.Key).toPEM(true), {
                    keyid: req.params.kid,
                    algorithm: "RS256",
                    ...(req.query.expires ? { expiresIn: req.query.expires as string } : {}),
                });

                res.set("content-type", "text/plain");
                res.set("X-Auth-Token", token);
                res.send(token);
            } catch (error) {
                res.status(400);
                res.send("payload misformed");
            }
        }
    });

    appHttp.post("/token", bodyParser.urlencoded({ extended: true }), async (req, res) => {
        // using client_id as key ID
        console.log(req.body);
        const kid = req.body.client_id;
        console.log("minting JWT with key:", kid);

        const keys = keystore.all({ kid });
        if (keys.length === 0) {
            res.status(404);
            res.send("unknown key ID: " + kid);
        } else if (keys.length > 1) {
            res.status(409);
            res.send(`too many keys for key ID <${kid}>: ${keys.length}`);
        } else {
            try {
                const token = jwt.sign(req.body, (keys[0] as unknown as jose.JWK.Key).toPEM(true), {
                    keyid: kid,
                    algorithm: "RS256",
                    ...(req.query.expires ? { expiresIn: req.query.expires as string } : {}),
                });

                res.json({
                    access_token: token,
                });
            } catch (error) {
                res.status(400);
                res.send("payload misformed");
            }
        }
    });

    appHttp.get("/ca", async (req, res) => {
        console.log("/ca");
        res.set("content-type", "text/plain");
        res.send(forge.pki.certificateToPem(certCA));
    });

    const httpsServer = https.createServer(
        {
            ca: forge.pki.certificateToPem(certCA),
            cert: [forge.pki.certificateToPem(cert)],
            key: forge.pki.privateKeyToPem(cert.privateKey),
        },
        appHttps
    );

    const httpServer = http.createServer(appHttp);

    httpsServer.listen(HTTPS_PORT, () => {
        console.log(`pairgen listening, HTTPS port: ${HTTPS_PORT}`);
    });

    httpServer.listen(HTTP_PORT, () => {
        console.log(`pairgen listening, HTTP port:  ${HTTP_PORT}`);
    });
})();

function generateCACertificate(attrs: forge.pki.CertificateField[]): forge.pki.Certificate {
    const pki = forge.pki;

    const keys = pki.rsa.generateKeyPair(2048);
    const cert = pki.createCertificate();

    cert.publicKey = keys.publicKey;
    cert.privateKey = keys.privateKey;

    cert.serialNumber = "01";
    cert.validity.notBefore = new Date();
    cert.validity.notAfter = new Date();
    cert.validity.notAfter.setFullYear(cert.validity.notBefore.getFullYear() + 1);
    cert.setSubject(attrs);

    cert.setIssuer(attrs);

    cert.setExtensions([
        {
            name: "basicConstraints",
            cA: true,
        },
        {
            name: "keyUsage",
            keyCertSign: true,
            digitalSignature: true,
            nonRepudiation: true,
            keyEncipherment: true,
            dataEncipherment: true,
        },
        {
            name: "extKeyUsage",
            serverAuth: true,
            clientAuth: true,
            codeSigning: true,
            emailProtection: true,
            timeStamping: true,
        },
        {
            name: "nsCertType",
            client: true,
            server: true,
            email: true,
            objsign: true,
            sslCA: true,
            emailCA: true,
            objCA: true,
        },
        // Keep for future reference.
        // {
        //     name: "subjectAltName",
        //     altNames: [
        //         {
        //             type: 6, // URI
        //             value: "http://example.org/webid#me"
        //         },
        //         {
        //             type: 7, // IP
        //             ip: "127.0.0.1"
        //         }
        //     ]
        // },
        // {
        //     name: "subjectKeyIdentifier"
        // }
    ]);

    cert.sign(cert.privateKey, forge.md.sha256.create());

    return cert;
}

function generateIntermediateCertificate(
    dnsNames: string[],
    attrs: forge.pki.CertificateField[],
    signingCert: forge.pki.Certificate
): forge.pki.Certificate {
    const pki = forge.pki;

    const keys = pki.rsa.generateKeyPair(2048);
    const cert = pki.createCertificate();

    cert.publicKey = keys.publicKey;
    cert.privateKey = keys.privateKey;

    cert.serialNumber = "01";
    cert.validity.notBefore = new Date();
    cert.validity.notAfter = new Date();
    cert.validity.notAfter.setFullYear(cert.validity.notBefore.getFullYear() + 1);

    cert.setSubject([
        ...attrs,
        {
            name: "commonName",
            value: dnsNames[0],
        },
    ]);

    cert.setIssuer(signingCert.subject.attributes);

    cert.setExtensions([
        {
            name: "basicConstraints",
            cA: false,
        },
        // Keep for future reference.
        // {
        //     name: "keyUsage",
        //     keyCertSign: true,
        //     digitalSignature: true,
        //     nonRepudiation: true,
        //     keyEncipherment: true,
        //     dataEncipherment: true
        // },
        // {
        //     name: "extKeyUsage",
        //     serverAuth: true,
        //     clientAuth: true,
        //     codeSigning: false,
        //     emailProtection: false,
        //     timeStamping: false
        // },
        // {
        //     name: "nsCertType",
        //     client: true,
        //     server: true,
        //     email: true,
        //     objsign: true,
        //     sslCA: true,
        //     emailCA: true,
        //     objCA: true
        // },
        {
            name: "subjectAltName",
            altNames: [
                ...dnsNames.map((dnsName) => ({ type: 2, value: dnsName })),

                // Keep for future reference.
                // {
                //     type: 6, // URI
                //     value: "https://localhost:5081"
                // },
                // {
                //     type: 7, // IP
                //     ip: "127.0.0.1"
                // }
            ],
        },
    ]);

    cert.sign(signingCert.privateKey, forge.md.sha256.create());

    return cert;
}

function pem2sshd(pem: string): string {
    const lines = pem.split("\r\n");
    if (lines.length < 2) {
        return "<error>";
    }

    return lines.slice(1, -2).join("");
}

async function pause(millis: number): Promise<void> {
    return new Promise(function (resolve, _reject) {
        setTimeout(resolve, millis);
    });
}
