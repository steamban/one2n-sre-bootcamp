#!/bin/bash

set -e

install_dependencies() {
    echo "Installing dependencies..."
    sudo apt-get update
    sudo apt-get install -y curl build-essential
}

install_docker() {
    echo "Installing Docker..."
    sudo apt-get install -y docker.io
}

install_docker_compose() {
    echo "Installing Docker Compose..."
    sudo curl -L "https://github.com/docker/compose/releases/download/v5.1.3/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
    sudo chmod +x /usr/local/bin/docker-compose
}

configure_docker() {
    echo "Configuring Docker service..."
    sudo systemctl enable docker
    sudo systemctl start docker
    sudo usermod -aG docker vagrant
}

verify_installation() {
    echo "Verifying installations..."
    docker --version
    docker-compose --version
}

main() {
    install_dependencies
    install_docker
    install_docker_compose
    configure_docker
    verify_installation
    echo "Provisioning complete!"
}

main "$@"