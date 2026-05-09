# Sample Domain-Driven Design Application

### Applied:
- Hexagonal Architecture
- Domain-Driven Design
- CQRS (command side)

### Testing

##### Requirements

- [Docker](https://docs.docker.com/get-docker/)
- [just](https://github.com/casey/just)
- [gotestsum](https://github.com/gotestyourself/gotestsum)

```bash
# E2E tests (start the API first with `just run-test`)
just run-test
just e2e
```
