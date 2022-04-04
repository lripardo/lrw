package middlewares

const DefaultCustomHeaders = `[
	"Cache-Control: no-cache, no-store, max-age=0, must-revalidate",
	"X-Content-Type-Options: nosniff",
	"X-Frame-Options: DENY",
	"Referrer-Policy: no-referrer",
	"X-XSS-Protection: 1; mode=block"
]`

/*
Example of headers

access-control-allow-origin: {http://myapp.com}
Access-Control-Expose-Headers: X-Frame-Options, X-Download-Options, Access-Control-Expose-Headers, X-Permitted-Cross-Domain-Policies, X-Xss-Protection, X-Signature, Content-Type, Access-Control-Allow-Origin, X-Content-Type-Options, X-Download-Options, Referrer-Policy, Strict-Transport-Security, X-Frame-Options, Strict-Transport-Security, X-Permitted-Cross-Domain-Policies, Pragma, X-Content-Type-Options, Content-Security-Policy, Content-Type, X-Http2-Stream-Id, X-Xss-Protection, Content-Security-Policy, Cache-Control
Cache-Control: no-store
Connection: keep-alive
Content-Length: {length}
content-security-policy: object-src 'none'; script-src 'unsafe-inline' 'unsafe-eval' 'strict-dynamic' https: http:;
Content-Type: application/json;charset=utf-8
Pragma: no-cache
Referrer-Policy: same-origin
Strict-Transport-Security: max-age=31536000; includeSubdomains
X-Content-Type-Options: nosniff
X-Download-Options: noopen
x-frame-options: DENY
x-http2-stream-id: {streamID}
X-Permitted-Cross-Domain-Policies: none
x-signature: {signatureHashHere}
x-xss-protection: 1; mode=block

*/
