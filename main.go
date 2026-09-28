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

    run("docker", "compose", "up", "--build", "-d")
    run("docker", "compose", "ps")
    fmt.Println("\nОткройте: http://localhost:5000")
    fmt.Println("Для просмотра логов: docker compose logs")
}
