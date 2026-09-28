package main

import (
    "fmt"
    "os"
    "os/exec"
)

func run(name string, args ...string) {
    fmt.Printf("\n>>> docker %s\n", join(args))
    cmd := exec.Command(name, args...)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    if err := cmd.Run(); err != nil {
        fmt.Printf("Команда завершилась с ошибкой: %v\n", err)
    }
}

func join(a []string) string {
    s := ""
    for i, x := range a {
        if i > 0 { s += " " }
        s += x
    }
    return s
}

func main() {

    run("docker", "rm", "-f", "registry")
    run("docker", "run", "-d", "--name", "registry", "-p", "5000:5000", "--restart=always", "registry:2")
    run("docker", "pull", "nginx:alpine")
    run("docker", "tag", "nginx:alpine", "localhost:5000/my-nginx:1.0")
    run("docker", "push", "localhost:5000/my-nginx:1.0")
    run("docker", "run", "--rm", "curlimages/curl:8.10.1", "http://host.docker.internal:5000/v2/_catalog")
    run("docker", "run", "--rm", "curlimages/curl:8.10.1", "http://host.docker.internal:5000/v2/my-nginx/tags/list")
    fmt.Println("\nДля Docker Hub используйте отдельную команду docker login, затем укажите собственное имя аккаунта.")
}
