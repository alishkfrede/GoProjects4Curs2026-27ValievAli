import React from "react";
import { createRoot } from "react-dom/client";

function App() {
  return (
    <main style={{fontFamily: "Arial, sans-serif", padding: "40px"}}>
      <h1>React-приложение в Docker</h1>
      <p>Многоэтапная сборка: Node.js → nginx.</p>
      <p>Практическая работа №2 выполнена.</p>
    </main>
  );
}

createRoot(document.getElementById("root")).render(<App />);
