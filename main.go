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

    run("docker", "pull", "nginx:alpine")
    run("docker", "pull", "alpine:latest")
    run("docker", "network", "ls")
    run("docker", "volume", "ls")

    run("docker", "rm", "-f", "web1", "tools1", "tmpfs-demo", "writer", "reader")
    run("docker", "network", "rm", "lab-net", "isolated-net")
    run("docker", "volume", "rm", "shared-data")

    run("docker", "run", "-d", "--name", "web1", "-p", "8080:80", "nginx:alpine")
    run("docker", "run", "-d", "--name", "tools1", "alpine:latest", "sleep", "1d")
    run("docker", "ps")

    run("docker", "network", "create", "lab-net")
    run("docker", "network", "connect", "lab-net", "web1")
    run("docker", "network", "connect", "lab-net", "tools1")
    run("docker", "exec", "tools1", "ping", "-c", "4", "web1")
    run("docker", "exec", "tools1", "wget", "-qO-", "http://web1")

    run("docker", "network", "create", "--internal", "isolated-net")
    run("docker", "network", "connect", "isolated-net", "tools1")
    run("docker", "exec", "tools1", "sh", "-c", "apk update")

    run("docker", "run", "-d", "--name", "tmpfs-demo", "--tmpfs", "/app/tmp:rw,size=64m", "alpine:latest", "sleep", "1d")
    run("docker", "inspect", "tmpfs-demo")
    run("docker", "exec", "tmpfs-demo", "sh", "-c", "mount | grep /app/tmp")

    run("docker", "volume", "create", "shared-data")
    run("docker", "run", "-d", "--name", "writer", "-v", "shared-data:/shared", "alpine:latest", "sleep", "1d")
    run("docker", "run", "-d", "--name", "reader", "-v", "shared-data:/shared", "alpine:latest", "sleep", "1d")
    run("docker", "exec", "writer", "sh", "-c", `echo "Docker shared volume test" > /shared/info.txt`)
    run("docker", "exec", "writer", "sh", "-c", "date > /shared/date.txt")
    run("docker", "exec", "reader", "ls", "-la", "/shared")
    run("docker", "exec", "reader", "cat", "/shared/info.txt")
    run("docker", "volume", "inspect", "shared-data")
}
