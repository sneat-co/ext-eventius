# Eventius contract

`@sneat/extension-eventius-contract` is the public, implementation-free API
for Eventius. It exposes Eventius DTOs, service interfaces, and Angular DI
tokens. Hosts provide concrete implementations from a private Eventius runtime;
consumers must not import implementation classes from this package.
