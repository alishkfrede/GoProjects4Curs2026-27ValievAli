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

    run("docker", "info")
    run("docker", "info", "--format", "{{.Driver}}")
    run("docker", "info", "--format", "{{.DockerRootDir}}")
    run("docker", "system", "df")
    fmt.Println("\nПроверка daemon.json выполняется после размещения файла в /etc/docker/daemon.json на Linux-хосте.")
}
