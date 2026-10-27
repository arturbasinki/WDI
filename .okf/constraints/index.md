# constraints

* [Host port 8080 collision](/constraints/host-port-collision.md) - Port 8080 is taken on the host by an unrelated service, so local access, probes and documentation use 127.0.0.1:8099 while the container keeps 8080.
* [Standard library only](/constraints/stdlib-only.md) - The product uses only the language's standard library and the browser's built-ins: no third-party dependencies, no frontend build step.
* [Workspace-local build cache](/constraints/workspace-cache.md) - The working environment keeps build cache, toolchain workspace and container-CLI home inside the repository (read-only home), forcing scoped checks and local-only verification scripts.
