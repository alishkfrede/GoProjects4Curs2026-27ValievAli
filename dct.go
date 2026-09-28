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

    if len(os.Args) != 2 {
        fmt.Println("Использование: go run dct.go YOUR_DOCKERHUB_USERNAME")
        os.Exit(1)
    }
    user := os.Args[1]
    os.Setenv("DOCKER_CONTENT_TRUST", "1")
    run("docker", "login")
    run("docker", "push", user+"/my-nginx:1.0")
}
