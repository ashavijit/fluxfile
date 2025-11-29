# Docker plugin for FluxFile
# This file provides docker-related tasks and utilities

task docker-shell:
    docker: true
    run:
        docker run -it --rm ${PROJECT}:${VERSION} /bin/sh

task docker-logs:
    run:
        docker logs -f ${PROJECT}

task docker-clean:
    run:
        docker system prune -af
