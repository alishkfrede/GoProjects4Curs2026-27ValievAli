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

    run("docker", "build", "-t", "my-react-app:1.0", ".")
    run("docker", "images")
    run("docker", "rm", "-f", "react-app")
    run("docker", "run", "-d", "--name", "react-app", "-p", "8080:80", "my-react-app:1.0")
    run("docker", "ps")
    fmt.Println("\nОткройте в браузере: http://localhost:8080")
}
